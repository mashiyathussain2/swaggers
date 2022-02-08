//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_user.go -package=mock go-app/app User

package app

import (
	"context"
	"fmt"
	"sync"

	"go-app/model"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/config"

	"strconv"
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

// User contains methods to implement authentication and authorization of 3 different types of users.
// 1. Customer 2. Brand 3. Influencer
type User interface {
	CreateUser(*schema.CreateUserOpts) (*schema.CreateUserResp, error)
	GetUserByEMail(string) (*schema.GetUserResp, error)
	GetUserInfoByID(*schema.GetUserInfoByIDOpts) (bson.M, error)
	EmailLoginCustomerUser(*schema.EmailLoginCustomerOpts) (auth.Claim, error)
	VerifyEmail(*schema.VerifyEmailOpts) (bool, error)
	CheckEmail(*schema.CheckEmailOpts) (bool, error)
	VerifyPhoneNo(*schema.VerifyPhoneNoOpts) (bool, error)
	CheckPhoneNo(*schema.CheckPhoneNoOpts) (bool, error)
	ResendConfirmationEmail(*schema.ResendVerificationEmailOpts) (bool, error)
	ForgotPassword(*schema.ForgotPasswordOpts) (bool, error)
	ResetPassword(*schema.ResetPasswordOpts) (bool, error)
	MobileLoginCustomerUser(*schema.MobileLoginCustomerUserOpts) (auth.Claim, error)
	GenerateMobileLoginOTP(*schema.GenerateMobileLoginOTPOpts) (bool, error)
	LoginWithSocial(*schema.LoginWithSocial) (auth.Claim, error)
	LoginWithApple(*schema.LoginWithApple) (auth.Claim, error)
	GetUserByID(primitive.ObjectID) (*model.User, error)
	GetUserClaim(*model.User, *model.Customer) auth.Claim

	UpdateUserEmail(*schema.UpdateUserEmailOpts) error
	UpdateUserPhoneNo(*schema.UpdateUserPhoneNoOpts) error
	GetUserIDByInfluencerID(opts *schema.GetUserInfoByIDOpts) (primitive.ObjectID, error)
}

// UserImpl implements user interface methods
type UserImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// UserImplOpts contains args required to create
type UserImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitUser returns new instance of user implementation
func InitUser(opts *UserImplOpts) User {
	ui := UserImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &ui
}

// CreateUser adds a new customer type user into user collection
func (ui *UserImpl) CreateUser(opts *schema.CreateUserOpts) (*schema.CreateUserResp, error) {
	// Validating user info before inserting it into DB
	if err := ui.validateCreateUser(opts); err != nil {
		return nil, err
	}

	user := model.User{
		CreatedAt: time.Now().UTC(),
	}
	if opts.Email != "" {
		user.Email = opts.Email
		user.Username = ui.generateUniqueUsername(opts.Email)
	}
	if opts.MobileNo != nil {
		user.PhoneNo = &model.PhoneNumber{
			Prefix: opts.MobileNo.Prefix,
			Number: opts.MobileNo.Number,
		}
	}

	switch opts.Type {
	case model.CustomerType:
		user.Role = model.UserRole
		user.Type = model.CustomerType
	}

	// Setting up encrypted password
	p, err := HashPassword(opts.Password, ui.App.Config.TokenAuthConfig.HashPasswordCost)
	if err != nil {
		return nil, err
	}
	user.Password = p

	// Setting up 6 digit confirmation code
	user.EmailVerificationCode, _ = GenerateOTP(ui.App.Config.TokenAuthConfig.OTPLength)

	if err := ui.sendConfirmationEmail(&user); err != nil {
		ui.Logger.Err(err).Interface("user", user).Msgf("failed to send confirmation email:%s", opts.Email)
	}

	res, err := ui.DB.Collection(model.UserColl).InsertOne(context.TODO(), user)
	if err != nil {
		ui.Logger.Err(err).Interface("user", user).Msg("failed to create user")
		return nil, errors.Wrap(err, "failed to create user")
	}

	return &schema.CreateUserResp{
		ID:      res.InsertedID.(primitive.ObjectID),
		Type:    user.Type,
		Email:   user.Email,
		PhoneNo: user.PhoneNo,
	}, nil
}

func (ui *UserImpl) validateCreateUser(opts *schema.CreateUserOpts) error {
	// Checking if user email already exists
	count, err := ui.DB.Collection(model.UserColl).CountDocuments(context.TODO(), bson.M{"email": opts.Email})
	if err != nil {
		ui.Logger.Err(err).Msgf("failed count if any document with email:%s already exists", opts.Email)
		return errors.Wrapf(err, "failed to validate user with email:%s", opts.Email)
	}
	if count != 0 {
		return errors.Errorf("user with email:%s already exists", opts.Email)
	}
	return nil
}

func (ui *UserImpl) generateUniqueUsername(email string) string {
	cmp := strings.Split(email, "@")
	username := cmp[0]
	filter := bson.M{
		"username": primitive.Regex{
			Pattern: username,
			Options: "i",
		},
	}
	count, _ := ui.DB.Collection(model.UserColl).CountDocuments(context.TODO(), filter)
	if count != 0 {
		return fmt.Sprintf("%s%s", username, strconv.Itoa(int(count)))
	}
	return username
}

func (ui *UserImpl) sendConfirmationEmail(u *model.User) error {
	htmlBody := fmt.Sprintf(`
		<p>Welcome! Thanks for signing up. Here's your verification otp:</p>
		<h3>%s</h3>
		<br>
		<p>Cheers!</p>
		<p>Team hypd!</p>`, u.EmailVerificationCode,
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
				Data:    aws.String("HYPD: Email Verification OTP"),
			},
		},
		Source: aws.String("hello@hypd.in"),
	}
	_, err := ui.App.SES.SendEmail(input)
	if err != nil {
		ui.Logger.Err(err).Msgf("failed to send verification otp to email:%s", u.Email)
		return err
	}
	return nil
}

func (ui *UserImpl) sendConfirmationOTP(u *model.User) error {
	// Sending OTP to phone number via SNS
	params := schema.SendOTPOpts{
		PhoneNo: *u.PhoneNo,
		OTP:     u.PhoneVerificationCode,
	}
	err := ui.App.Kaleyra.SendOTP(&params)
	if err != nil {
		return errors.Wrapf(err, "failed to send otp")
	}
	return nil
}

func (ui *UserImpl) sendForgotPasswordOTPEmail(u *model.User) error {
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
	_, err := ui.App.SES.SendEmail(input)
	if err != nil {
		ui.Logger.Err(err).Msgf("failed to send password reset otp to email:%s", u.Email)
		return err
	}
	return nil
}

func (ui *UserImpl) GetUserClaim(user *model.User, customer *model.Customer) auth.Claim {
	claim := auth.UserClaim{
		ID:           user.ID.Hex(),
		CustomerID:   customer.ID.Hex(),
		CartID:       customer.CartID.Hex(),
		Type:         user.Type,
		Role:         user.Role,
		Email:        user.Email,
		PhoneNo:      user.PhoneNo,
		CreatedVia:   user.CreatedVia,
		FullName:     customer.FullName,
		DOB:          customer.DOB,
		ProfileImage: customer.ProfileImage,
	}
	if customer.Gender != nil {
		claim.Gender = *customer.Gender
	}
	if !user.EmailVerifiedAt.IsZero() {
		claim.EmailVerified = true
	}
	if !user.PhoneVerifiedAt.IsZero() {
		claim.PhoneVerified = true
	}
	if !user.InfluencerID.IsZero() {
		if influencer, err := ui.App.Influencer.GetInfluencerByID(user.InfluencerID); err == nil {
			claim.InfluencerInfo = &auth.InfluencerInfo{
				ID:           user.InfluencerID.Hex(),
				Name:         influencer.Name,
				ProfileImage: influencer.ProfileImage,
			}
		}
	}

	return &claim
}

// VerifyEmail verify email with received verification code
func (ui *UserImpl) VerifyEmail(opts *schema.VerifyEmailOpts) (bool, error) {
	var user model.User
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"email": opts.Email}).Decode(&user); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return false, errors.Wrapf(err, "user with email:%s not found", opts.Email)
		}
	}
	if !user.EmailVerifiedAt.IsZero() {
		return false, errors.Errorf("email:%s already verified", user.Email)
	}
	if user.EmailVerificationCode != opts.VerificationCode {
		return false, errors.New("invalid verification code")
	}

	f := bson.M{"_id": user.ID}
	u := bson.M{
		"$set": bson.M{
			"email_verified_at": time.Now().UTC(),
		},
		"$unset": bson.M{
			"email_verification_code": 1,
		},
	}
	_, err := ui.DB.Collection(model.UserColl).UpdateOne(context.TODO(), f, u)
	if err != nil {
		return false, errors.Wrapf(err, "failed to update email:%s as verified", user.Email)
	}

	return true, nil
}

// CheckEmail checks email if valid of not
func (ui *UserImpl) CheckEmail(opts *schema.CheckEmailOpts) (bool, error) {
	ctx := context.TODO()
	var user model.User
	filter := bson.M{"email": opts.Email}
	if err := ui.DB.Collection(model.UserColl).FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return true, nil
		}
		return false, errors.Errorf("failed to check for user with email: %s", opts.Email)
	}
	if user != (model.User{}) {
		return false, errors.Errorf("user with email: %s already exists", opts.Email)
	}
	return true, nil
}

// CheckPhoneNo checks phone no if valid of not
func (ui *UserImpl) CheckPhoneNo(opts *schema.CheckPhoneNoOpts) (bool, error) {
	ctx := context.TODO()
	var user model.User
	filter := bson.M{"phone_no.prefix": opts.PhoneNo.Prefix, "phone_no.number": opts.PhoneNo.Number}
	if err := ui.DB.Collection(model.UserColl).FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return true, nil
		}
		return false, errors.Errorf("failed to check for user with phone no: %s", opts.PhoneNo.Number)
	}
	if user != (model.User{}) {
		return false, errors.Errorf("user with phone no: %s already exists", opts.PhoneNo.Number)
	}
	return true, nil
}

// VerifyEmail verify phone no with received verification code
func (ui *UserImpl) VerifyPhoneNo(opts *schema.VerifyPhoneNoOpts) (bool, error) {
	var user model.User
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"phone_no.prefix": opts.PhoneNo.Prefix, "phone_no.number": opts.PhoneNo.Number}).Decode(&user); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return false, errors.Wrapf(err, "user with phone no:%s not found", opts.PhoneNo.Number)
		}
	}
	if !user.PhoneVerifiedAt.IsZero() {
		return false, errors.Errorf("phone :%s already verified", user.PhoneNo.Number)
	}
	if user.PhoneVerificationCode != opts.VerificationCode {
		return false, errors.New("invalid verification code")
	}

	f := bson.M{"_id": user.ID}
	u := bson.M{
		"$set": bson.M{
			"phone_verified_at": time.Now().UTC(),
		},
		"$unset": bson.M{
			"phone_verification_code": 1,
		},
	}
	_, err := ui.DB.Collection(model.UserColl).UpdateOne(context.TODO(), f, u)
	if err != nil {
		return false, errors.Wrapf(err, "failed to update phone no:%s as verified", user.PhoneNo.Number)
	}

	return true, nil
}

// ResendConfirmationEmail generates a new confirmation token and resends its to user email
func (ui *UserImpl) ResendConfirmationEmail(opts *schema.ResendVerificationEmailOpts) (bool, error) {
	var user model.User
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"email": opts.Email}).Decode(&user); err != nil {
		return false, errors.Wrapf(err, "user with email:%s not found", opts.Email)
	}
	if !user.EmailVerifiedAt.IsZero() {
		return false, errors.Errorf("email:%s already verified", user.Email)
	}
	verificationCode, _ := GenerateOTP(ui.App.Config.TokenAuthConfig.OTPLength)

	u := bson.M{
		"$set": bson.M{
			"email_verification_code": verificationCode,
		},
	}
	if _, err := ui.DB.Collection(model.UserColl).UpdateOne(context.TODO(), bson.M{"_id": user.ID}, u); err != nil {
		return false, errors.Wrap(err, "failed to generate verification code")
	}

	if err := ui.sendConfirmationEmail(&user); err != nil {
		ui.Logger.Err(err).Interface("user", user).Msgf("failed to send confirmation email:%s", opts.Email)
	}
	return true, nil
}

// EmailLoginCustomerUser allows customer to login via email and password
func (ui *UserImpl) EmailLoginCustomerUser(opts *schema.EmailLoginCustomerOpts) (auth.Claim, error) {
	var user model.User
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"email": opts.Email}).Decode(&user); err != nil {
		return nil, errors.Wrapf(err, "user with email:%s not found", opts.Email)
	}
	if !CheckPasswordHash(opts.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	var customer model.Customer
	if err := ui.DB.Collection(model.CustomerColl).FindOne(context.TODO(), bson.M{"user_id": user.ID}).Decode(&customer); err != nil {
		return nil, errors.Wrapf(err, "customer with email:%s not found", opts.Email)
	}

	claim := ui.GetUserClaim(&user, &customer)
	return claim, nil
}

// GetUserByEMail returns user info with filtered by email
func (ui *UserImpl) GetUserByEMail(email string) (*schema.GetUserResp, error) {
	var user model.User
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Wrapf(err, "user with email:%s", email)
		}
		return nil, errors.Wrapf(err, "failed to find user with email:%s", email)
	}
	resp := schema.GetUserResp{
		ID:         user.ID,
		Type:       user.Type,
		Role:       user.Role,
		Email:      user.Email,
		PhoneNo:    user.PhoneNo,
		Username:   user.Username,
		CreatedVia: user.CreatedVia,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	if !user.EmailVerifiedAt.IsZero() {
		resp.EmailVerified = true
	}
	if !user.PhoneVerifiedAt.IsZero() {
		resp.PhoneVerified = true
	}
	return &resp, nil
}

// ForgotPassword sends an otp to email to allow user to reset password
func (ui *UserImpl) ForgotPassword(opts *schema.ForgotPasswordOpts) (bool, error) {
	otp, _ := GenerateOTP(ui.App.Config.TokenAuthConfig.OTPLength)
	filter := bson.M{"email": opts.Email}
	update := bson.M{
		"$set": bson.M{
			"password_reset_code": otp,
		},
	}

	res, err := ui.DB.Collection(model.UserColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ui.Logger.Err(err).Interface("user", opts).Msgf("failed to generate password_reset_code for user email:%s", opts.Email)
		return false, errors.Wrapf(err, "failed to generate password_reset_code for user email:%s", opts.Email)
	}
	if res.MatchedCount == 0 {
		return false, errors.Errorf("user with email:%s not found", opts.Email)
	}

	// Sending Email
	if err := ui.sendForgotPasswordOTPEmail(&model.User{Email: opts.Email, PasswordResetCode: otp}); err != nil {
		ui.Logger.Err(err).Msgf("failed to send forgot password otp over email:%s", opts.Email)
	}
	return true, nil
}

// ResetPassword change existing user password by matching the otp from user and in password_reset_field
func (ui *UserImpl) ResetPassword(opts *schema.ResetPasswordOpts) (bool, error) {
	var user model.User
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"email": opts.Email}).Decode(&user); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return false, errors.Wrapf(err, "user with email:%s not found", opts.Email)
		}
	}
	if user.PasswordResetCode != opts.OTP {
		return false, errors.New("invalid otp")
	}
	password, _ := HashPassword(opts.Password, ui.App.Config.TokenAuthConfig.HashPasswordCost)
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"password_reset_code": "",
		"email_verified_at":   time.Now().UTC(),
		"password":            password,
	}}

	res, err := ui.DB.Collection(model.UserColl).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ui.Logger.Err(err).Interface("user", opts).Msgf("failed to reset password for user email:%s", opts.Email)
		return false, errors.Wrapf(err, "failed to reset password for user user email:%s", opts.Email)
	}
	if res.MatchedCount == 0 {
		return false, errors.Errorf("user with email:%s not found", opts.Email)
	}
	return true, nil
}

// MobileLoginCustomerUser allows user to login via phone
func (ui *UserImpl) MobileLoginCustomerUser(opts *schema.MobileLoginCustomerUserOpts) (auth.Claim, error) {
	var user model.User
	filter := bson.M{"phone_no.prefix": opts.PhoneNo.Prefix, "phone_no.number": opts.PhoneNo.Number}
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), filter).Decode(&user); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Wrapf(err, "user with phone:%s%s not found", opts.PhoneNo.Prefix, opts.PhoneNo.Number)
		}
	}

	if user.LoginOTP == nil || opts.OTP != user.LoginOTP.OTP || user.LoginOTP.Type != model.PhoneLoginOTPType {
		return nil, errors.New("invalid otp")
	}

	if int(time.Now().UTC().Sub(user.LoginOTP.CreatedAt.UTC()).Minutes()) > 15 {
		return nil, errors.New("otp expired")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		filter := bson.M{
			"_id": user.ID,
		}
		update := bson.D{
			{
				Key: "$unset",
				Value: bson.M{
					"login_otp": 1,
				},
			},
		}
		if user.PhoneVerifiedAt.IsZero() {
			update = append(update, bson.E{Key: "$set", Value: bson.M{"phone_verified_at": time.Now().UTC()}})
		}
		_, err := ui.DB.Collection(model.UserColl).UpdateOne(context.TODO(), filter, update)
		ui.Logger.Err(err).Msg("failed to unset otp")
	}()

	var customer model.Customer
	if err := ui.DB.Collection(model.CustomerColl).FindOne(context.TODO(), bson.M{"user_id": user.ID}).Decode(&customer); err != nil {
		return nil, errors.Wrapf(err, "customer with phone_no:%s%s not found", opts.PhoneNo.Prefix, opts.PhoneNo.Number)
	}
	claim := ui.GetUserClaim(&user, &customer)

	wg.Wait()
	return claim, nil
}

// GenerateMobileLoginOTP checks provided phone number in DB and sends an otp.
// If phone number does not exists then it create a new user and sends the otp.
// GenerateMobileLoginOTP checks provided phone number in DB and sends an otp.
// If phone number does not exists then it create a new user and sends the otp.
func (ui *UserImpl) GenerateMobileLoginOTP(opts *schema.GenerateMobileLoginOTPOpts) (bool, error) {
	ctx := context.TODO()
	filter := bson.M{"phone_no.prefix": opts.PhoneNo.Prefix, "phone_no.number": opts.PhoneNo.Number}
	count, err := ui.DB.Collection(model.UserColl).CountDocuments(context.TODO(), filter)
	if err != nil {
		ui.Logger.Err(err).Msgf("failed to check for user with phone_no:%s%s", opts.PhoneNo.Prefix, opts.PhoneNo.Number)
		return false, errors.Wrapf(err, "failed to check for user with phone_no:%s%s", opts.PhoneNo.Prefix, opts.PhoneNo.Number)
	}
	otp, _ := GenerateOTP(6)
	if config.GetConfig().ServerConfig.Env == "staging" {
		if opts.PhoneNo.Number == "1000000001" || opts.PhoneNo.Number == "1000000002" ||
			opts.PhoneNo.Number == "2000000001" || opts.PhoneNo.Number == "2000000002" ||
			opts.PhoneNo.Number == "3000000001" || opts.PhoneNo.Number == "3000000002" ||
			opts.PhoneNo.Number == "9999999999" {
			otp = "000000"
		}
	}
	switch count {
	// When no user exists thus creating a new one
	case 0:
		// creating new user
		user := model.User{
			Type: model.CustomerType,
			Role: model.UserColl,
			PhoneNo: &model.PhoneNumber{
				Prefix: opts.PhoneNo.Prefix,
				Number: opts.PhoneNo.Number,
			},
			LoginOTP: &model.LoginOTP{
				Type:      model.PhoneLoginOTPType,
				OTP:       otp,
				CreatedAt: time.Now().UTC(),
			},
			CreatedVia: model.CreateViaMobile,
			CreatedAt:  time.Now().UTC(),
		}
		res, err := ui.DB.Collection(model.UserColl).InsertOne(ctx, user)
		if err != nil {
			ui.Logger.Err(err).Msgf("failed to generate user using phone_no:%s%s", opts.PhoneNo.Prefix, opts.PhoneNo.Number)
			return false, errors.Wrapf(err, "failed to generate user using phone_no:%s%s", opts.PhoneNo.Prefix, opts.PhoneNo.Number)
		}
		user.ID = res.InsertedID.(primitive.ObjectID)

		// creating customer and linking user_id
		customer := model.Customer{
			UserID:    user.ID,
			CreatedAt: time.Now().UTC(),
		}
		res, err = ui.DB.Collection(model.CustomerColl).InsertOne(ctx, customer)
		if err != nil {
			ui.Logger.Err(err).Msgf("failed to generate customer using phone_no:%s%s", opts.PhoneNo.Prefix, opts.PhoneNo.Number)
			return false, errors.Wrapf(err, "failed to generate customer using phone_no:%s%s", opts.PhoneNo.Prefix, opts.PhoneNo.Number)
		}
	// When user exists
	default:
		// Setting up otp
		update := bson.M{
			"$set": bson.M{
				"login_otp": model.LoginOTP{
					Type:      model.PhoneLoginOTPType,
					OTP:       otp,
					CreatedAt: time.Now().UTC(),
				},
			},
		}
		_, err := ui.DB.Collection(model.UserColl).UpdateOne(ctx, filter, update)
		if err != nil {
			return false, errors.Wrap(err, "failed to generate otp for login")
		}
	}

	if config.GetConfig().ServerConfig.Env == "staging" {
		if opts.PhoneNo.Number == "1000000001" || opts.PhoneNo.Number == "1000000002" ||
			opts.PhoneNo.Number == "2000000001" || opts.PhoneNo.Number == "2000000002" ||
			opts.PhoneNo.Number == "3000000001" || opts.PhoneNo.Number == "3000000002" ||
			opts.PhoneNo.Number == "9999999999" {
			return true, nil
		}
	}
	// Sending OTP to phone number via SNS
	// params := &sns.PublishInput{
	// 	Message:     aws.String(fmt.Sprintf("OTP for login: %s", otp)),
	// 	PhoneNumber: aws.String(fmt.Sprintf("%s%s", opts.PhoneNo.Prefix, opts.PhoneNo.Number)),
	// }
	err = ui.App.Kaleyra.SendOTP(&schema.SendOTPOpts{
		PhoneNo: model.PhoneNumber{
			Prefix: opts.PhoneNo.Prefix,
			Number: opts.PhoneNo.Number,
		},
		OTP: otp,
	})
	// _, err = ui.App.SNS.Publish(params)
	if err != nil {
		ui.Logger.Err(err).Interface("opts", opts).Msg("failed to send otp")
		return false, errors.Wrap(err, "failed to send otp")
	}
	return true, nil
}

// LoginWithSocial allows customer to login with social accounts such as google, facebook
func (ui *UserImpl) LoginWithSocial(opts *schema.LoginWithSocial) (auth.Claim, error) {
	ctx := context.TODO()
	var user model.User
	var customer model.Customer
	var newUser bool
	filter := bson.M{"email": opts.Email}
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), filter).Decode(&user); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			newUser = true
		} else {
			ui.Logger.Err(err).Msgf("failed to check for social user with email:%s", opts.Email)
			return nil, errors.Wrapf(err, "failed to check for social user with email:%s", opts.Email)
		}
	}

	switch newUser {
	case true:
		// creating new user
		user = model.User{
			Type:            model.CustomerType,
			Role:            model.UserRole,
			Email:           opts.Email,
			EmailVerifiedAt: time.Now().UTC(),
			CreatedAt:       time.Now().UTC(),
		}
		switch opts.Type {
		case "google":
			user.CreatedVia = model.CreatedViaGoogle
		case "facebook":
			user.CreatedVia = model.CreatedViaFacebook
		}
		res, err := ui.DB.Collection(model.UserColl).InsertOne(ctx, user)
		if err != nil {
			ui.Logger.Err(err).Interface("opts", opts).Msgf("failed to create social user using email:%s", opts.Email)
			return nil, errors.Wrapf(err, "failed to create social user using email:%s", opts.Email)
		}
		user.ID = res.InsertedID.(primitive.ObjectID)

		// creating customer and linking user_id
		customer = model.Customer{
			UserID:    user.ID,
			FullName:  opts.FullName,
			CreatedAt: time.Now().UTC(),
		}
		if opts.ProfileImage != nil {
			customer.ProfileImage = &model.IMG{
				SRC: opts.ProfileImage.SRC,
			}
			if err := customer.ProfileImage.LoadFromURL(); err != nil {
				customer.ProfileImage = nil
			}
		}
		res, err = ui.DB.Collection(model.CustomerColl).InsertOne(ctx, customer)
		if err != nil {
			ui.Logger.Err(err).Interface("opts", opts).Msgf("failed to create social customer using email:%s", opts.Email)
			return nil, errors.Wrapf(err, "failed to create social customer using email:%s", opts.Email)
		}
		customer.ID = res.InsertedID.(primitive.ObjectID)
	default:
		if user.CreatedVia != model.CreatedViaFacebook && user.CreatedVia != model.CreatedViaGoogle && user.CreatedVia != model.CreatedViaApple {
			return nil, errors.New("cannot use social login for this user, please use email/otp login")
		}
		if user.CreatedVia == model.CreatedViaGoogle {
			if opts.Type == model.CreatedViaFacebook || opts.Type == model.CreatedViaApple {
				return nil, errors.New("cannot use facebook login: this account was created via other social login")
			}
		}
		if user.CreatedVia == model.CreatedViaFacebook {
			if opts.Type == model.CreatedViaGoogle || opts.Type == model.CreatedViaApple {
				return nil, errors.New("cannot use google login: this account was created via other social login")
			}
		}
		if user.CreatedVia == model.CreatedViaApple {
			if opts.Type == model.CreatedViaGoogle || opts.Type == model.CreatedViaFacebook {
				return nil, errors.New("cannot use google login: this account was created via other social login")
			}
		}
		filterQuery := bson.M{"user_id": user.ID}
		var update bson.D
		if opts.FullName != "" {
			update = append(update, bson.E{Key: "full_name", Value: opts.FullName})
		}
		if opts.ProfileImage != nil {
			img := model.IMG{
				SRC: opts.ProfileImage.SRC,
			}
			if err := img.LoadFromURL(); err == nil {
				update = append(update, bson.E{Key: "profile_image", Value: img})
			}
		}
		updateQuery := bson.M{
			"$set": update,
		}
		optsQuery := options.FindOneAndUpdate().SetReturnDocument(options.After)
		if err := ui.DB.Collection(model.CustomerColl).FindOneAndUpdate(ctx, filterQuery, updateQuery, optsQuery).Decode(&customer); err != nil {
			return nil, errors.Wrapf(err, "failed to update social customer with email:%s", user.Email)
		}
	}

	claim := ui.GetUserClaim(&user, &customer)
	return claim, nil
}

func (ui *UserImpl) LoginWithApple(opts *schema.LoginWithApple) (auth.Claim, error) {
	ctx := context.TODO()
	var user model.User
	var customer model.Customer
	var newUser bool
	filter := bson.M{"social_id": opts.AppleID, "created_via": "apple"}
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), filter).Decode(&user); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			newUser = true
		} else {
			ui.Logger.Err(err).Msgf("failed to check for apple social user with id:%s", opts.AppleID)
			return nil, errors.Wrapf(err, "failed to check for apple social user with id:%s", opts.AppleID)
		}
	}

	switch newUser {
	case true:
		// creating new user
		user = model.User{
			Type:            model.CustomerType,
			Role:            model.UserRole,
			Email:           opts.Email,
			EmailVerifiedAt: time.Now().UTC(),
			CreatedAt:       time.Now().UTC(),
			CreatedVia:      model.CreatedViaApple,
			SocialID:        opts.AppleID,
		}
		res, err := ui.DB.Collection(model.UserColl).InsertOne(ctx, user)
		if err != nil {
			ui.Logger.Err(err).Interface("opts", opts).Msgf("failed to create social user using email:%s", opts.Email)
			return nil, errors.Wrapf(err, "failed to create social user using email:%s", opts.Email)
		}
		user.ID = res.InsertedID.(primitive.ObjectID)

		// creating customer and linking user_id
		customer = model.Customer{
			UserID:    user.ID,
			FullName:  opts.FullName,
			CreatedAt: time.Now().UTC(),
		}
		res, err = ui.DB.Collection(model.CustomerColl).InsertOne(ctx, customer)
		if err != nil {
			ui.Logger.Err(err).Interface("opts", opts).Msgf("failed to create social customer using email:%s", opts.Email)
			return nil, errors.Wrapf(err, "failed to create social customer using email:%s", opts.Email)
		}
		customer.ID = res.InsertedID.(primitive.ObjectID)
	default:
		filterQuery := bson.M{"user_id": user.ID}
		if err := ui.DB.Collection(model.CustomerColl).FindOne(ctx, filterQuery).Decode(&customer); err != nil {
			return nil, errors.Wrapf(err, "failed to apple social customer with id:%s", user.ID)
		}
	}

	claim := ui.GetUserClaim(&user, &customer)
	return claim, nil
}

func (ui *UserImpl) GetUserInfoByID(opts *schema.GetUserInfoByIDOpts) (bson.M, error) {
	matchStage := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"_id": opts.ID,
			},
		},
	}

	lookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         model.CustomerColl,
				"localField":   "_id",
				"foreignField": "user_id",
				"as":           "customer_info",
			},
		},
	}

	setStage := bson.D{
		{
			Key: "$set",
			Value: bson.M{
				"customer_info": bson.M{
					"$arrayElemAt": bson.A{
						"$customer_info",
						0,
					},
				},
			},
		},
	}

	projectStage := bson.D{
		{
			Key: "$project",
			Value: bson.M{
				"email":    1,
				"phone_no": 1,
				"username": 1,
				"type":     1,

				"id":            "$_id",
				"customer_id":   "$customer_info._id",
				"full_name":     "$customer_info.full_name",
				"profile_image": "$customer_info.profile_image",
			},
		},
	}

	var res []bson.M
	ctx := context.TODO()
	cur, err := ui.DB.Collection(model.UserColl).Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, setStage, projectStage})
	if err != nil {
		return nil, errors.Wrapf(err, "query failed to find user by id: %", opts.ID.Hex())
	}

	if err := cur.All(ctx, &res); err != nil {
		return nil, errors.Wrapf(err, "failed to find user by id: %s", opts.ID.Hex())
	}

	if len(res) == 0 {
		return nil, errors.Errorf("user with id: %s not found", opts.ID.Hex())
	}
	return res[0], nil
}

func (ui *UserImpl) GetUserByID(id primitive.ObjectID) (*model.User, error) {
	filter := bson.M{
		"_id": id,
	}
	var user model.User
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return nil, errors.Wrap(err, "failed to find user with id")
	}
	return &user, nil
}

func (ui *UserImpl) UpdateUserEmail(opts *schema.UpdateUserEmailOpts) error {
	ctx := context.TODO()
	// Checking if another user with email already exists
	var filter bson.M
	var update bson.M
	var user model.User
	// var wg sync.WaitGroup
	filter = bson.M{"email": opts.Email}
	if err := ui.DB.Collection(model.UserColl).FindOne(ctx, filter).Decode(&user); err != nil {
		if err != mongo.ErrNoDocuments && err != mongo.ErrNilDocument {
			return errors.Wrapf(err, "failed to check for user with email: %s", opts.Email)
		}
	}

	// user already exists with the email
	if (user != model.User{}) {
		// If its a different user return error
		if user.ID != opts.ID {
			return errors.Errorf("email: %s is associated with different user", opts.Email)
		}
	}

	// no account has provided email thus added email to provided user
	claimOTP, _ := GenerateOTP(6)
	update = bson.M{
		"$set": bson.M{
			"email":                   opts.Email,
			"email_verification_code": claimOTP,
		},
	}
	filter = bson.M{"_id": opts.ID}
	queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := ui.DB.Collection(model.UserColl).FindOneAndUpdate(ctx, filter, update, queryOpts).Decode(&user); err != nil {
		return errors.Wrap(err, "failed to update user info")
	}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	ui.sendConfirmationEmail(&user)
	// }()

	// wg.Wait()
	return nil
}

func (ui *UserImpl) UpdateUserPhoneNo(opts *schema.UpdateUserPhoneNoOpts) error {
	ctx := context.TODO()
	// Checking if another user with phone no already exists
	var filter bson.M
	var update bson.M
	// var wg sync.WaitGroup
	var user model.User
	filter = bson.M{"phone_no.prefix": opts.PhoneNo.Prefix, "phone_no.number": opts.PhoneNo.Number}
	if err := ui.DB.Collection(model.UserColl).FindOne(ctx, filter).Decode(&user); err != nil {
		if err != mongo.ErrNoDocuments && err != mongo.ErrNilDocument {
			return errors.Wrapf(err, "failed to check for user with phone: %s", opts.PhoneNo.Number)
		}
	}

	// user already exists with the phone
	if (user != model.User{}) {
		// If its a different user return error
		if user.ID != opts.ID {
			return errors.Errorf("phone no: %s is associated with different user", opts.PhoneNo.Number)
		}
	}
	// no account has provided phone no thus added phone no to provided user
	claimOTP, _ := GenerateOTP(6)
	filter = bson.M{"_id": opts.ID}
	update = bson.M{
		"$set": bson.M{
			"phone_no": model.PhoneNumber{
				Prefix: opts.PhoneNo.Prefix,
				Number: opts.PhoneNo.Number,
			},
			"phone_verification_code": claimOTP,
		},
	}
	queryOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err := ui.DB.Collection(model.UserColl).FindOneAndUpdate(ctx, filter, update, queryOpts).Decode(&user); err != nil {
		return errors.Wrap(err, "failed to update user info")
	}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	ui.sendConfirmationOTP(&user)
	// }()
	// wg.Wait()
	return nil
}

func (ui *UserImpl) GetUserIDByInfluencerID(opts *schema.GetUserInfoByIDOpts) (primitive.ObjectID, error) {
	filter := bson.M{"influencer_id": opts.ID}
	var user model.User
	if err := ui.DB.Collection(model.UserColl).FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return primitive.NilObjectID, errors.Wrap(err, "failed to find user by influencer id")
	}
	return user.ID, nil
}
