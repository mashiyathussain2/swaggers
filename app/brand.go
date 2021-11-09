package app

import (
	"context"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"go-app/server/auth"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Brand contains methods for brand service functionality
type Brand interface {
	CreateBrand(*schema.CreateBrandOpts) (*schema.CreateBrandResp, error)
	EditBrand(*schema.EditBrandOpts) (*schema.EditBrandResp, error)
	GetBrandByID(primitive.ObjectID) (*schema.GetBrandResp, error)
	CheckBrandByID(primitive.ObjectID) (bool, error)
	GetBrandsByID([]primitive.ObjectID) ([]schema.GetBrandResp, error)
	GetBrands() ([]schema.GetBrandResp, error)

	// SendPasswordEmail(string, string) (bool, error)
	AddFollower(opts *schema.AddBrandFollowerOpts) (bool, error)
	RemoveFollower(opts *schema.AddBrandFollowerOpts) (bool, error)

	CreateBrandAdminUser(*schema.CreateBrandAdminUserOpts) (bool, error)

	BrandUserLogin(*schema.BrandUserLoginOpts) (auth.Claim, error)
	ForgotPassword(*schema.ForgotPasswordOpts) (bool, error)
	ResetPassword(*schema.ResetPasswordOpts) (bool, error)
	CheckBrandUsernameExists(string, *mongo.SessionContext) error
}

// BrandImpl implements brand interface methods
type BrandImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// BrandImplOpts contains args required to create
type BrandImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitBrand returns new instance of brand implementation
func InitBrand(opts *BrandImplOpts) Brand {
	ui := BrandImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &ui
}

// CreateBrand inserts a new brand document into collection
// Note: this method only creates brand profile not brand authentication
func (bi *BrandImpl) CreateBrand(opts *schema.CreateBrandOpts) (*schema.CreateBrandResp, error) {

	ctx := context.TODO()
	session, err := bi.DB.Client().StartSession()
	if err != nil {
		bi.Logger.Err(err).Msg("unable to create db session")
		return nil, errors.Wrap(err, "failed to add follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		bi.Logger.Err(err).Msg("unable to start transaction")
		return nil, errors.Wrap(err, "failed to add follower")
	}
	var b model.Brand
	var res *mongo.InsertOneResult
	var sp []schema.GetSizeProfileForBrandResp

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		b = model.Brand{
			Name:               opts.Name,
			LName:              strings.ToLower(opts.Name),
			RegisteredName:     opts.RegisteredName,
			FulfillmentEmail:   opts.FulfillmentEmail,
			FulfillmentCCEmail: opts.FulfillmentCCEmail,
			Domain:             opts.Domain,
			Website:            opts.Website,
			Logo: &model.IMG{
				SRC: opts.Logo.SRC,
			},
			Bio: opts.Bio,
			CoverImg: &model.IMG{
				SRC: opts.CoverImg.SRC,
			},
			CreatedAt: time.Now().UTC(),
		}
		//TODO: check if username is unique

		err = bi.CheckBrandUsernameExists(opts.Username, &sc)
		if err != nil {
			return err
		}
		b.Username = opts.Username

		if err := b.Logo.LoadFromURL(); err != nil {
			return errors.Wrap(err, "invalid image for brand logo")
		}
		if err := b.CoverImg.LoadFromURL(); err != nil {
			return errors.Wrap(err, "invalid image for brand cover")
		}
		if opts.SocialAccount != nil {
			b.SocialAccount = &model.SocialAccount{}
			if opts.SocialAccount.Facebook != nil {
				b.SocialAccount.Facebook = &model.SocialMedia{
					FollowersCount: uint(opts.SocialAccount.Facebook.FollowersCount),
					URL:            opts.SocialAccount.Facebook.URL,
				}
			}
			if opts.SocialAccount.Instagram != nil {
				b.SocialAccount.Instagram = &model.SocialMedia{
					FollowersCount: uint(opts.SocialAccount.Instagram.FollowersCount),
					URL:            opts.SocialAccount.Instagram.URL,
				}
			}
			if opts.SocialAccount.Youtube != nil {
				b.SocialAccount.Youtube = &model.SocialMedia{
					FollowersCount: uint(opts.SocialAccount.Youtube.FollowersCount),
					URL:            opts.SocialAccount.Youtube.URL,
				}
			}
			if opts.SocialAccount.Twitter != nil {
				b.SocialAccount.Twitter = &model.SocialMedia{
					FollowersCount: uint(opts.SocialAccount.Twitter.FollowersCount),
					URL:            opts.SocialAccount.Twitter.URL,
				}
			}
		}
		if len(opts.SizeProfiles) > 0 {
			b.SizeProfiles = opts.SizeProfiles
		}

		res, err = bi.DB.Collection(model.BrandColl).InsertOne(sc, b)
		if err != nil {
			session.AbortTransaction(sc)
			bi.Logger.Err(err).Interface("opts", opts).Msg("failed to insert brand")
			return errors.Wrap(err, "failed to create brand")
		}
		err = bi.App.SizeProfile.AddBrandToSizeProfile(&schema.AddBrandToSizeProfileOpts{
			IDs:     opts.SizeProfiles,
			BrandID: res.InsertedID.(primitive.ObjectID),
		})
		if err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msg("failed to link size profiles with brand id")
		}

		sp, err = bi.App.SizeProfile.GetSizeProfilesForBrand(res.InsertedID.(primitive.ObjectID))
		if err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msg("failed to get size profiles for brand with id")
		}

		if err := session.CommitTransaction(sc); err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed to commit transaction")
			return errors.Wrap(err, "failed to create brand")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	resp := schema.CreateBrandResp{
		ID:                 res.InsertedID.(primitive.ObjectID),
		Name:               b.Name,
		Username:           b.Username,
		RegisteredName:     b.RegisteredName,
		FulfillmentEmail:   b.FulfillmentEmail,
		FulfillmentCCEmail: b.FulfillmentCCEmail,
		Domain:             b.Domain,
		Website:            b.Website,
		Logo:               b.Logo,
		CoverImg:           b.CoverImg,
		SocialAccount:      b.SocialAccount,
		Bio:                b.Bio,
		CreatedAt:          b.CreatedAt,
		SizeProfiles:       sp,
	}
	return &resp, nil
}

func (bi *BrandImpl) EditBrand(opts *schema.EditBrandOpts) (*schema.EditBrandResp, error) {

	ctx := context.TODO()
	session, err := bi.DB.Client().StartSession()
	if err != nil {
		bi.Logger.Err(err).Msg("unable to create db session")
		return nil, errors.Wrap(err, "failed to add follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		bi.Logger.Err(err).Msg("unable to start transaction")
		return nil, errors.Wrap(err, "failed to add follower")
	}
	var brand model.Brand
	var sp []schema.GetSizeProfileForBrandResp

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		var update bson.D

		if opts.Name != "" {
			update = append(update, bson.E{Key: "name", Value: opts.Name})
			update = append(update, bson.E{Key: "lname", Value: string(opts.Name)})
		}

		if opts.Username != "" {
			//TODO: check if username is valid
			err = bi.CheckBrandUsernameExists(opts.Username, &sc)
			if err != nil {
				return err
			}
			update = append(update, bson.E{Key: "username", Value: opts.Username})
		}

		if opts.Domain != "" {
			update = append(update, bson.E{Key: "domain", Value: opts.Domain})
		}
		if opts.RegisteredName != "" {
			update = append(update, bson.E{Key: "registered_name", Value: opts.RegisteredName})
		}
		if opts.FulfillmentEmail != "" {
			update = append(update, bson.E{Key: "fulfillment_email", Value: opts.FulfillmentEmail})
		}
		if len(opts.FulfillmentCCEmail) != 0 {
			update = append(update, bson.E{Key: "fulfillment_cc_email", Value: opts.FulfillmentCCEmail})
		}
		if opts.Website != "" {
			update = append(update, bson.E{Key: "website", Value: opts.Website})
		}
		if opts.CoverImg != nil {
			img := model.IMG{SRC: opts.CoverImg.SRC}
			if err := img.LoadFromURL(); err != nil {
				return errors.Wrap(err, "invalid image for brand cover")
			}
			update = append(update, bson.E{Key: "cover_img", Value: img})
		}
		if opts.Logo != nil {
			img := model.IMG{SRC: opts.Logo.SRC}
			if err := img.LoadFromURL(); err != nil {
				return errors.Wrap(err, "invalid image for brand logo")
			}
			update = append(update, bson.E{Key: "logo", Value: img})
		}
		if opts.SocialAccount != nil {
			if opts.SocialAccount.Facebook != nil {
				update = append(update, bson.E{Key: "social_account.facebook.url", Value: opts.SocialAccount.Facebook.URL})
				update = append(update, bson.E{Key: "social_account.facebook.followers_count", Value: opts.SocialAccount.Facebook.FollowersCount})
			}
			if opts.SocialAccount.Instagram != nil {
				update = append(update, bson.E{Key: "social_account.instagram.url", Value: opts.SocialAccount.Instagram.URL})
				update = append(update, bson.E{Key: "social_account.instagram.followers_count", Value: opts.SocialAccount.Instagram.FollowersCount})
			}
			if opts.SocialAccount.Youtube != nil {
				update = append(update, bson.E{Key: "social_account.youtube.url", Value: opts.SocialAccount.Youtube.URL})
				update = append(update, bson.E{Key: "social_account.youtube.followers_count", Value: opts.SocialAccount.Youtube.FollowersCount})
			}
			if opts.SocialAccount.Twitter != nil {
				update = append(update, bson.E{Key: "social_account.twitter.url", Value: opts.SocialAccount.Twitter.URL})
				update = append(update, bson.E{Key: "social_account.twitter.followers_count", Value: opts.SocialAccount.Twitter.FollowersCount})
			}
		}
		if opts.Bio != "" {
			update = append(update, bson.E{Key: "bio", Value: opts.Bio})
		}
		if len(opts.SizeProfiles) != 0 {
			update = append(update, bson.E{Key: "size_profiles", Value: opts.SizeProfiles})
			bi.App.SizeProfile.AddBrandToSizeProfile(&schema.AddBrandToSizeProfileOpts{IDs: opts.SizeProfiles, BrandID: opts.ID})
		}
		update = append(update, bson.E{Key: "updated_at", Value: time.Now().UTC()})

		filterQuery := bson.M{"_id": opts.ID}
		updateQuery := bson.M{"$set": update}

		queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		if err := bi.DB.Collection(model.BrandColl).FindOneAndUpdate(sc, filterQuery, updateQuery, queryOpts).Decode(&brand); err != nil {
			if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
				return errors.Wrapf(err, "brand with id:%s not found", opts.ID.Hex())
			}
			return errors.Wrapf(err, "failed to update brand with id:%s", opts.ID.Hex())
		}
		sp, err = bi.App.SizeProfile.GetSizeProfilesForBrand(brand.ID)
		if err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msg("failed to get size profiles for brand with id")
		}

		if err := session.CommitTransaction(sc); err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed to commit transaction")
			return errors.Wrap(err, "failed to create brand")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	resp := schema.EditBrandResp{
		ID:                 brand.ID,
		Name:               brand.Name,
		Username:           brand.Username,
		RegisteredName:     brand.RegisteredName,
		FulfillmentEmail:   brand.FulfillmentEmail,
		FulfillmentCCEmail: brand.FulfillmentCCEmail,
		Domain:             brand.Domain,
		Website:            brand.Website,
		Logo:               brand.Logo,
		CoverImg:           brand.CoverImg,
		SocialAccount:      brand.SocialAccount,
		Bio:                brand.Bio,
		SizeProfiles:       sp,
		CreatedAt:          brand.CreatedAt,
		UpdatedAt:          brand.UpdatedAt,
	}
	return &resp, nil
}

// GetBrandByID returns brand info with matching id
func (bi *BrandImpl) GetBrandByID(id primitive.ObjectID) (*schema.GetBrandResp, error) {

	var resp []schema.GetBrandResp
	ctx := context.TODO()
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": id,
		},
	}}
	lookupStage := bson.D{{

		Key: "$lookup", Value: bson.M{
			"from":         model.SizeProfileColl,
			"localField":   "size_profiles",
			"foreignField": "_id",
			"as":           "size_profiles",
		},
	}}

	cursor, err := bi.DB.Collection(model.BrandColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get brand data")
	}

	if err := cursor.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "error decoding brands")
	}

	return &resp[0], nil
}

// CheckBrandByID check if brand exists with matching id
func (bi *BrandImpl) CheckBrandByID(id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}
	count, err := bi.DB.Collection(model.BrandColl).CountDocuments(context.TODO(), filter)
	if err != nil {
		bi.Logger.Err(err).Msgf("failed to check brand with id:%s", id.Hex())
		return false, errors.Wrapf(err, "failed to check brand with id:%s", id.Hex())
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}

// GetBrandByID returns brand info with matching id
func (bi *BrandImpl) GetBrandsByID(ids []primitive.ObjectID) ([]schema.GetBrandResp, error) {
	var resp []schema.GetBrandResp
	ctx := context.TODO()
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"_id": bson.M{"$in": ids},
		},
	}}
	lookupStage := bson.D{{

		Key: "$lookup", Value: bson.M{
			"from":         model.SizeProfileColl,
			"localField":   "size_profiles",
			"foreignField": "_id",
			"as":           "size_profiles",
		},
	}}

	cursor, err := bi.DB.Collection(model.BrandColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get brand data")
	}

	if err := cursor.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "error decoding brands")
	}

	return resp, nil
}

func (bi *BrandImpl) GetBrands() ([]schema.GetBrandResp, error) {
	var resp []schema.GetBrandResp
	ctx := context.TODO()

	lookupStage := bson.D{{

		Key: "$lookup", Value: bson.M{
			"from":         model.SizeProfileColl,
			"localField":   "size_profiles",
			"foreignField": "_id",
			"as":           "size_profiles",
		},
	}}

	cursor, err := bi.DB.Collection(model.BrandColl).Aggregate(ctx, mongo.Pipeline{
		lookupStage,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get brands data")
	}

	if err := cursor.All(ctx, &resp); err != nil {
		return nil, errors.Wrap(err, "error decoding brands")
	}

	return resp, nil
}

func (bi *BrandImpl) AddFollower(opts *schema.AddBrandFollowerOpts) (bool, error) {
	ctx := context.TODO()
	session, err := bi.DB.Client().StartSession()
	if err != nil {
		bi.Logger.Err(err).Msg("unable to create db session")
		return false, errors.Wrap(err, "failed to add follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		bi.Logger.Err(err).Msg("unable to start transaction")
		return false, errors.Wrap(err, "failed to add follower")
	}

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		isFollowing, err := bi.DB.Collection(model.BrandColl).CountDocuments(sc, bson.M{"_id": opts.BrandID, "followers_id": opts.CustomerID})
		if err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msg("failed to check is user already follow brand")
			session.AbortTransaction(sc)
			return errors.Wrap(err, "failed to follow brand")
		}

		if isFollowing != 0 {
			session.AbortTransaction(sc)
			return errors.New("user already follow the brand")
		}

		filter := bson.M{
			"_id": opts.BrandID,
		}
		update := bson.M{
			"$addToSet": bson.M{
				"followers_id": opts.CustomerID,
			},
			"$inc": bson.M{
				"followers_count": 1,
			},
		}

		res, err := bi.DB.Collection(model.BrandColl).UpdateOne(sc, filter, update)
		if err != nil {
			session.AbortTransaction(sc)
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed add follower")
			return errors.Wrap(err, "failed to add follower")
		}

		if res.MatchedCount == 0 {
			session.AbortTransaction(sc)
			return errors.New("brand not found")
		}

		if err := bi.App.Customer.AddBrandFollowing(sc, opts); err != nil {
			session.AbortTransaction(sc)
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed add brand_id in customer following")
			return errors.Wrap(err, "failed to add follower")
		}

		if err := session.CommitTransaction(sc); err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed to commit transaction")
			return errors.Wrap(err, "failed to add follower")
		}
		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (bi *BrandImpl) RemoveFollower(opts *schema.AddBrandFollowerOpts) (bool, error) {
	ctx := context.TODO()
	session, err := bi.DB.Client().StartSession()
	if err != nil {
		bi.Logger.Err(err).Msg("unable to create db session")
		return false, errors.Wrap(err, "failed to add follower")
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		bi.Logger.Err(err).Msg("unable to start transaction")
		return false, errors.Wrap(err, "failed to add follower")
	}

	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {

		isFollowing, err := bi.DB.Collection(model.BrandColl).CountDocuments(sc, bson.M{"_id": opts.BrandID, "followers_id": opts.CustomerID})
		if err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msg("failed to check is user already follow brand")
			session.AbortTransaction(sc)
			return errors.Wrap(err, "failed to unfollow brand")
		}

		if isFollowing == 0 {
			session.AbortTransaction(sc)
			return errors.New("user does not follow the brand")
		}

		filter := bson.M{
			"_id": opts.BrandID,
		}
		update := bson.M{
			"$pull": bson.M{
				"followers_id": opts.CustomerID,
			},
			"$inc": bson.M{
				"followers_count": -1,
			},
		}

		res, err := bi.DB.Collection(model.BrandColl).UpdateOne(sc, filter, update)
		if err != nil {
			session.AbortTransaction(sc)
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed remove follower")
			return errors.Wrap(err, "failed to remove follower")
		}

		if res.MatchedCount == 0 {
			session.AbortTransaction(sc)
			return errors.New("brand not found")
		}

		if err := bi.App.Customer.RemoveBrandFollowing(sc, opts); err != nil {
			session.AbortTransaction(sc)
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed remove brand_id in customer following")
			return errors.Wrap(err, "failed to remove follower")
		}

		if err := session.CommitTransaction(sc); err != nil {
			bi.Logger.Err(err).Interface("opts", opts).Msgf("failed to commit transaction")
			return errors.Wrap(err, "failed to remove follower")
		}
		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (bi *BrandImpl) CreateBrandAdminUser(opts *schema.CreateBrandAdminUserOpts) (bool, error) {
	ctx := context.TODO()
	count, err := bi.DB.Collection(model.BrandUserColl).CountDocuments(ctx, bson.M{"email": opts.Email})
	if err != nil {
		return false, errors.Wrap(err, "failed to check for existing user")
	}
	if count > 0 {
		return false, errors.Errorf("user with email:%s already exists", opts.Email)
	}
	p, _ := GeneratePassword(8)
	password, _ := HashPassword(p, bi.App.Config.TokenAuthConfig.HashPasswordCost)
	user := model.BrandUser{
		BrandId:   opts.BrandID,
		Email:     opts.Email,
		Password:  password,
		Role:      model.AdminRole,
		CreatedAt: time.Now().UTC(),
	}
	if _, err := bi.DB.Collection(model.BrandUserColl).InsertOne(ctx, user); err != nil {
		return false, errors.Wrap(err, "failed to create brand user")
	}
	bi.sendBrandUserEmail(user.Email, p)
	return true, nil
}

func (bi *BrandImpl) sendBrandUserEmail(email, password string) (bool, error) {
	htmlBody := fmt.Sprintf(`
		<p>Welcome! Here's your email & password for brand dashboard login</p>
		<h3>Email: %s</h3>
		<h3>Password: %s</h3>
		<br>
		<p>Cheers!</p>
		<p>Team hypd!</p>`, email, password,
	)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    aws.String("HYPD: Brand Account Login Details"),
			},
		},
		Source: aws.String("hello@hypd.in"),
	}
	_, err := bi.App.SES.SendEmail(input)
	if err != nil {
		bi.Logger.Err(err).Msgf("failed to brand account details email:%s", email)
	}
	return true, nil
}

func (bi *BrandImpl) SendBrandUserEmail(email string) error {
	filter := bson.M{
		"email": email,
	}
	p, _ := GeneratePassword(8)
	password, _ := HashPassword(p, bi.App.Config.TokenAuthConfig.HashPasswordCost)
	update := bson.M{
		"$set": bson.M{
			"password": password,
		},
	}
	if _, err := bi.DB.Collection(model.BrandUserColl).UpdateOne(context.TODO(), filter, update); err != nil {
		return errors.Wrap(err, "failed to send password")
	}
	bi.sendBrandUserEmail(email, p)
	return nil
}

func (ui *BrandImpl) getBrandUserClaim(user *model.BrandUser) (auth.Claim, error) {
	claim := auth.UserClaim{
		ID:           user.ID.Hex(),
		Type:         model.BrandType,
		Role:         user.Role,
		Email:        user.Email,
		ProfileImage: user.ProfileImage,
	}
	var brand model.BrandClaim
	if err := ui.DB.Collection(model.BrandColl).FindOne(context.TODO(), bson.M{"_id": user.BrandId}).Decode(&brand); err != nil {
		return nil, errors.Wrapf(err, "failed to get brand info associated with this user")
	}
	claim.BrandInfo = &brand
	return &claim, nil
}

func (bi *BrandImpl) BrandUserLogin(opts *schema.BrandUserLoginOpts) (auth.Claim, error) {
	var user model.BrandUser
	if err := bi.DB.Collection(model.BrandUserColl).FindOne(context.TODO(), bson.M{"email": opts.Email}).Decode(&user); err != nil {
		return nil, errors.Wrapf(err, "user with email:%s not found", opts.Email)
	}
	if !CheckPasswordHash(opts.Password, user.Password) {
		return nil, errors.New("invalid password")
	}
	return bi.getBrandUserClaim(&user)
}

func (bi *BrandImpl) sendForgotPasswordOTPEmail(u *model.BrandUser) error {
	htmlBody := fmt.Sprintf(`
		<p>Here's is your to reset your account's otp:</p>
		<h3>%s</h3>
		<br>
		<p>Cheers!</p>
		<p>Team hypd!</p>`, u.PasswordResetCode,
	)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(u.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    aws.String("HYPD: Password Reset OTP"),
			},
		},

		Source: aws.String("hello@hypd.in"),
	}
	_, err := bi.App.SES.SendEmail(input)
	if err != nil {
		bi.Logger.Err(err).Msgf("failed to send password reset otp to email:%s", u.Email)
		return err
	}
	return nil
}

// ForgotPassword sends an otp to email to allow user to reset password
func (bi *BrandImpl) ForgotPassword(opts *schema.ForgotPasswordOpts) (bool, error) {
	otp, _ := GenerateOTP(bi.App.Config.TokenAuthConfig.OTPLength)
	filter := bson.M{"email": opts.Email}
	update := bson.M{
		"$set": bson.M{
			"password_reset_code": otp,
		},
	}

	res, err := bi.DB.Collection(model.BrandUserColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		bi.Logger.Err(err).Interface("user", opts).Msgf("failed to generate password_reset_code for user email:%s", opts.Email)
		return false, errors.Wrapf(err, "failed to generate password_reset_code for user email:%s", opts.Email)
	}
	if res.MatchedCount == 0 {
		return false, errors.Errorf("user with email:%s not found", opts.Email)
	}

	// Sending Email
	if err := bi.sendForgotPasswordOTPEmail(&model.BrandUser{Email: opts.Email, PasswordResetCode: otp}); err != nil {
	}
	return true, nil
}

// ResetPassword change existing user password by matching the otp from user and in password_reset_field
func (bi *BrandImpl) ResetPassword(opts *schema.ResetPasswordOpts) (bool, error) {
	var user model.User
	if err := bi.DB.Collection(model.BrandUserColl).FindOne(context.TODO(), bson.M{"email": opts.Email}).Decode(&user); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return false, errors.Wrapf(err, "user with email:%s not found", opts.Email)
		}
	}
	if user.PasswordResetCode != opts.OTP {
		return false, errors.New("invalid otp")
	}
	password, _ := HashPassword(opts.Password, bi.App.Config.TokenAuthConfig.HashPasswordCost)
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"password_reset_code": "",
		"password":            password,
	}}

	res, err := bi.DB.Collection(model.BrandUserColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		bi.Logger.Err(err).Interface("user", opts).Msgf("failed to reset password for user email:%s", opts.Email)
		return false, errors.Wrapf(err, "failed to reset password for user user email:%s", opts.Email)
	}
	if res.MatchedCount == 0 {
		return false, errors.Errorf("user with email:%s not found", opts.Email)
	}
	return true, nil
}

func (bi *BrandImpl) CheckBrandUsernameExists(username string, sc *mongo.SessionContext) error {
	ctx := context.TODO()
	isAlpha := regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	if !isAlpha(username) {
		errors.Errorf("%s is not valid", username)
	}
	filter := bson.M{
		"username": username,
	}
	var brand *model.Brand
	var err error
	if sc != nil {
		err = bi.DB.Collection(model.BrandColl).FindOne(*sc, filter).Decode(&brand)
	} else {
		err = bi.DB.Collection(model.BrandColl).FindOne(ctx, filter).Decode(&brand)
	}
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNilValue || err == mongo.ErrNoDocuments {
			return nil
		}
		return errors.Wrapf(err, "error checking if username exists or not")
	}
	return errors.Errorf("username: %s alreasy exist", username)
}
