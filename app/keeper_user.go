package app

import (
	"context"
	"encoding/json"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"go-app/server/auth"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KeeperUser interface {
	Login() string
	Callback(state, code string) (auth.Claim, error)
	AddNewSessionID(userID primitive.ObjectID, sessionID string) error
}

type KeeperUserOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

type KeeperUserImpl struct {
	App         *App
	DB          *mongo.Database
	Logger      *zerolog.Logger
	GoogleOAuth GoogleOAuth
}

func InitKeeperUser(opts *KeeperUserOpts) KeeperUser {
	ku := KeeperUserImpl{
		App:         opts.App,
		Logger:      opts.Logger,
		DB:          opts.DB,
		GoogleOAuth: NewGoogleOAuth(&GoogleOAuthOpts{Config: &opts.App.Config.GoogleOAuth}),
	}
	return &ku
}

func (ku *KeeperUserImpl) Login() string {
	return ku.GoogleOAuth.AuthCodeURL()
}

func (ku *KeeperUserImpl) verifyCallback(state string, code string) (*schema.KeeperUserLoginOpts, error) {
	if state != ku.App.Config.GoogleOAuth.State {
		return nil, errors.New("invalid oauth state")
	}
	token, err := ku.GoogleOAuth.Exchange(code)
	if err != nil {
		return nil, errors.Wrap(err, "code exchange failed")
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting user info")
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	var res schema.KeeperUserLoginOpts
	if err := json.Unmarshal(contents, &res); err != nil {
		return nil, fmt.Errorf("failed docode response body: %s", err.Error())
	}

	return &res, nil
}

func (ku *KeeperUserImpl) Callback(state, code string) (auth.Claim, error) {
	res, err := ku.verifyCallback(state, code)
	if err != nil {
		return nil, err
	}
	if res.Domain != "hypd.in" {
		return nil, errors.Wrap(err, "not allowed")
	}

	user, err := ku.CreateOrUpdateUser(res)
	if err != nil {
		return nil, err
	}

	cliam := auth.UserClaim{
		ID:            user.UserInfo.ID.Hex(),
		KeeperUserID:  user.ID.Hex(),
		Type:          user.UserInfo.Type,
		Role:          user.UserInfo.Role,
		FullName:      user.FullName,
		Email:         user.UserInfo.Email,
		ProfileImage:  user.ProfileImage,
		CreatedVia:    user.UserInfo.CreatedVia,
		EmailVerified: user.UserInfo.EmailVerified,
	}

	return &cliam, nil
}

func (ku *KeeperUserImpl) CreateOrUpdateUser(opts *schema.KeeperUserLoginOpts) (*schema.KeeperUserInfoResp, error) {
	filter := bson.M{
		"email": opts.Email,
	}
	userSetOnInsert := model.User{
		Type:            model.KeeperType,
		Role:            model.UserRole,
		CreatedVia:      "google",
		CreatedAt:       time.Now().UTC(),
		EmailVerifiedAt: time.Now().UTC(),
	}

	userSet := model.User{
		UpdatedAt: time.Now().UTC(),
	}

	update := bson.M{
		"$set":         userSet,
		"$setOnInsert": userSetOnInsert,
	}

	queryOpts := options.Update().SetUpsert(true)
	_, err := ku.DB.Collection(model.UserColl).UpdateOne(context.TODO(), filter, update, queryOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create/update user")
	}

	opts2 := schema.CreateOrUpdateKeeperUser{
		Email:        opts.Email,
		FullName:     opts.Name,
		ProfileImage: &schema.Img{SRC: opts.Picture},
	}
	keeperUser, err := ku.CreateOrUpdateKeeperUser(&opts2)
	if err != nil {
		return nil, err
	}
	return keeperUser, nil
}

func (ku *KeeperUserImpl) CreateOrUpdateKeeperUser(opts *schema.CreateOrUpdateKeeperUser) (*schema.KeeperUserInfoResp, error) {
	user, err := ku.App.User.GetUserByEMail(opts.Email)
	if err != nil {
		return nil, err
	}

	keeperUserSetOnInsert := model.KeeperUser{
		FullName:  opts.FullName,
		CreatedAt: time.Now().UTC(),
	}
	keeperSet := model.KeeperUser{
		FullName: opts.FullName,
	}
	profileImage := model.IMG{SRC: opts.ProfileImage.SRC}
	if err := profileImage.LoadFromURL(); err == nil {
		keeperSet.ProfileImage = &profileImage
		keeperUserSetOnInsert.ProfileImage = &profileImage
	}

	filter := bson.M{
		"user_id": user.ID,
	}

	update := bson.M{
		"$set": keeperSet,
		"$setOnInsert": bson.M{
			"created_at": time.Now().UTC(),
		},
	}

	var res model.KeeperUser
	queryOpts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	if err := ku.DB.Collection(model.KeeperUserColl).FindOneAndUpdate(context.TODO(), filter, update, queryOpts).Decode(&res); err != nil {
		return nil, errors.Wrap(err, "failed to create/update keeper user")
	}

	resp := schema.KeeperUserInfoResp{
		ID:           res.ID,
		UserID:       res.UserID,
		UserInfo:     user,
		FullName:     res.FullName,
		ProfileImage: res.ProfileImage,
		CreatedAt:    res.CreatedAt,
	}

	return &resp, nil
}

func (ku *KeeperUserImpl) AddNewSessionID(userID primitive.ObjectID, sessionID string) error {
	filter := bson.M{
		"user_id": userID,
	}
	update := bson.M{
		"$addToSet": bson.M{
			"session_ids": sessionID,
		},
	}
	_, err := ku.DB.Collection(model.KeeperUserColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return errors.Wrap(err, "failed to add session id")
	}
	return nil
}
