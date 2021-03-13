package app

import (
	"context"
	"go-app/mock"
	"go-app/model"
	"go-app/schema"
	"go-app/server/auth"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"syreclabs.com/go/faker"
)

func TestUserImpl_validateCreateUser(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateUserOpts
	}

	type TC struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     error
		prepare func(*TC)
	}
	tests := []TC{
		{
			name: "[Ok] when no user with provided email exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.opts = schema.GetRandomCreateUserOpts()
			},
		},
		{
			name: "[Ok] when user with provided email exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.opts = schema.GetRandomCreateUserOpts()
				doc := bson.M{"email": tt.args.opts.Email}
				tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), doc)
				tt.err = errors.Errorf("user with email:%s already exists", tt.args.opts.Email)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			err := ui.validateCreateUser(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, err)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestUserImpl_generateUniqueUsername(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		email string
	}

	type TC struct {
		name    string
		fields  fields
		args    args
		want    string
		prepare func(*TC)
	}

	tests := []TC{
		{
			name: "[Ok] when no prefix username exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.email = faker.Internet().FreeEmail()
				tt.want = strings.Split(tt.args.email, "@")[0]
			},
		},
		{
			name: "[Ok] when 1 prefix username exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.email = faker.Internet().FreeEmail()
				doc := bson.M{"username": strings.Split(tt.args.email, "@")[0]}
				tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), doc)
				tt.want = strings.Split(tt.args.email, "@")[0] + "1"
			},
		},
		{
			name: "[Ok] when 2 prefix username exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.email = faker.Internet().FreeEmail()
				doc := []interface{}{
					bson.M{"username": strings.Split(tt.args.email, "@")[0]},
					bson.M{"username": strings.Split(tt.args.email, "@")[0] + "1"},
				}
				tt.fields.DB.Collection(model.UserColl).InsertMany(context.TODO(), doc)
				tt.want = strings.Split(tt.args.email, "@")[0] + "2"
			},
		},
		{
			name: "[Ok] when 10 prefix username exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.email = faker.Internet().FreeEmail()
				doc := []interface{}{
					bson.M{"username": strings.Split(tt.args.email, "@")[0]},
					bson.M{"username": strings.Split(tt.args.email, "@")[0] + "1"},
					bson.M{"username": strings.Split(tt.args.email, "@")[0] + "2"},
					bson.M{"username": strings.Split(tt.args.email, "@")[0] + "3"},
					bson.M{"username": strings.Split(tt.args.email, "@")[0] + "4"},
					bson.M{"username": strings.Split(tt.args.email, "@")[0] + "5"},
					bson.M{"username": strings.Split(tt.args.email, "@")[0] + "6"},
					bson.M{"username": strings.Split(tt.args.email, "@")[0] + "7"},
					bson.M{"username": strings.Split(tt.args.email, "@")[0] + "8"},
					bson.M{"username": strings.Split(tt.args.email, "@")[0] + "9"},
				}
				tt.fields.DB.Collection(model.UserColl).InsertMany(context.TODO(), doc)
				tt.want = strings.Split(tt.args.email, "@")[0] + "10"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			got := ui.generateUniqueUsername(tt.args.email)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserImpl_CreateUser(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateUserOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     *schema.CreateUserResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, *schema.CreateUserResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.opts = schema.GetRandomCreateUserOpts()
				tt.args.opts.Type = model.CustomerType
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreateUserResp) {
				assert.False(t, resp.ID.IsZero())
				var doc model.User
				tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				assert.WithinDuration(t, time.Now().UTC(), doc.CreatedAt, 200*time.Millisecond)
				assert.Equal(t, tt.args.opts.Email, doc.Email)
				assert.Equal(t, model.CustomerType, doc.Type)
				assert.Equal(t, strings.Split(tt.args.opts.Email, "@")[0], doc.Username)
				assert.Len(t, doc.EmailVerificationCode, 6)
				assert.True(t, CheckPasswordHash(tt.args.opts.Password, doc.Password))
			},
		},
		{
			name: "[Ok] Only Email",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.opts = schema.GetRandomCreateUserOpts()
				tt.args.opts.MobileNo = nil
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreateUserResp) {
				assert.False(t, resp.ID.IsZero())
				var doc model.User
				tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				assert.WithinDuration(t, time.Now().UTC(), doc.CreatedAt, 200*time.Millisecond)
				assert.Equal(t, tt.args.opts.Email, doc.Email)
				assert.Equal(t, strings.Split(tt.args.opts.Email, "@")[0], doc.Username)
				assert.Nil(t, doc.PhoneNo)
				assert.Len(t, doc.EmailVerificationCode, 6)
				assert.True(t, CheckPasswordHash(tt.args.opts.Password, doc.Password))
			},
		},
		{
			name: "[Ok] Only MobileNo",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				tt.args.opts = schema.GetRandomCreateUserOpts()
				tt.args.opts.Email = ""
			},
			validate: func(t *testing.T, tt *TC, resp *schema.CreateUserResp) {
				assert.False(t, resp.ID.IsZero())
				var doc model.User
				tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				assert.WithinDuration(t, time.Now().UTC(), doc.CreatedAt, 200*time.Millisecond)
				assert.Equal(t, tt.args.opts.MobileNo.Number, doc.PhoneNo.Number)
				assert.Equal(t, tt.args.opts.MobileNo.Prefix, doc.PhoneNo.Prefix)
				assert.Len(t, doc.EmailVerificationCode, 6)
				assert.True(t, CheckPasswordHash(tt.args.opts.Password, doc.Password))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			got, err := ui.CreateUser(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestUserImpl_VerifyEmail(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.VerifyEmailOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     bool
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, bool)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)
				tt.args.opts = &schema.VerifyEmailOpts{
					Email:            res.Email,
					VerificationCode: user.EmailVerificationCode,
				}
				tt.want = true
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				assert.Equal(t, tt.want, resp)
				var user model.User
				err := tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"email": tt.args.opts.Email}).Decode(&user)
				assert.Nil(t, err)
				assert.WithinDuration(t, time.Now().UTC(), user.EmailVerifiedAt, 200*time.Millisecond)
			},
		},
		{
			name: "[Error] Invalid Verification Code",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)
				tt.args.opts = &schema.VerifyEmailOpts{
					Email:            res.Email,
					VerificationCode: faker.Numerify("######"),
				}
				tt.err = errors.New("invalid verification code")
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				assert.Equal(t, tt.want, resp)
				var user model.User
				err := tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"email": tt.args.opts.Email}).Decode(&user)
				assert.Nil(t, err)
				assert.WithinDuration(t, time.Now().UTC(), user.EmailVerifiedAt, 200*time.Millisecond)
			},
		},
		{
			name: "[Error] Email already verified",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)

				tt.args.opts = &schema.VerifyEmailOpts{
					Email:            res.Email,
					VerificationCode: user.EmailVerificationCode,
				}
				tt.fields.App.User.VerifyEmail(tt.args.opts)
				tt.err = errors.Errorf("email:%s already verified", user.Email)
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				assert.Equal(t, tt.want, resp)
				var user model.User
				err := tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"email": tt.args.opts.Email}).Decode(&user)
				assert.Nil(t, err)
				assert.WithinDuration(t, time.Now().UTC(), user.EmailVerifiedAt, 200*time.Millisecond)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			got, err := ui.VerifyEmail(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.VerifyEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestUserImpl_ResendConfirmationEmail(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		user *model.User
		opts *schema.ResendVerificationEmailOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     bool
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, bool)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)
				tt.args.user = &user
				tt.args.opts = &schema.ResendVerificationEmailOpts{Email: user.Email}
				tt.want = true
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": tt.args.user.ID}).Decode(&user)
				assert.NotEqual(t, tt.args.user.EmailVerificationCode, user.EmailVerificationCode)
				assert.Len(t, user.EmailVerificationCode, 6)
			},
		},
		{
			name: "[Error] Invalid email",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)
				tt.args.user = &user
				tt.args.opts = &schema.ResendVerificationEmailOpts{Email: faker.Internet().FreeEmail()}
				tt.err = errors.Errorf("user with email:%s not found: mongo: no documents in result", tt.args.opts.Email)
			},
		},
		{
			name: "[Error] email verified",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)
				opts := schema.VerifyEmailOpts{
					Email:            user.Email,
					VerificationCode: user.EmailVerificationCode,
				}
				tt.fields.App.User.VerifyEmail(&opts)
				tt.args.user = &user
				tt.args.opts = &schema.ResendVerificationEmailOpts{Email: user.Email}
				tt.err = errors.Errorf("email:%s already verified", user.Email)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			got, err := ui.ResendConfirmationEmail(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.ResendConfirmationEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestUserImpl_GetUserByEMail(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		email string
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     *schema.GetUserResp
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, *schema.GetUserResp)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				resp, _ := tt.fields.App.User.CreateUser(createOpts)
				var doc model.User
				tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				tt.args.email = resp.Email
				tt.want = &schema.GetUserResp{
					ID:         doc.ID,
					Type:       doc.Type,
					Role:       doc.Role,
					Email:      doc.Email,
					PhoneNo:    doc.PhoneNo,
					Username:   doc.Username,
					CreatedVia: doc.CreatedVia,
					CreatedAt:  doc.CreatedAt,
					UpdatedAt:  doc.UpdatedAt,
				}
			},
			validate: func(t *testing.T, tt *TC, resp *schema.GetUserResp) {
				assert.Equal(t, tt.want, resp)
			},
		},
		{
			name: "[Error] user does not exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				resp, _ := tt.fields.App.User.CreateUser(createOpts)
				var doc model.User
				tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": resp.ID}).Decode(&doc)
				tt.args.email = faker.Internet().Email()
				tt.err = errors.Errorf("user with email:%s: mongo: no documents in result", tt.args.email)
			},
			validate: func(t *testing.T, tt *TC, resp *schema.GetUserResp) {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			got, err := ui.GetUserByEMail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.GetUserByEMail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestCustomerImpl_EmailLoginCustomerUser(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.EmailLoginCustomerOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     auth.Claim
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, auth.Claim)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				createUserOpts := schema.CreateUserOpts{
					Type:     model.CustomerType,
					Email:    faker.Internet().SafeEmail(),
					Password: faker.Internet().Password(6, 10),
				}
				tt.args.opts = &schema.EmailLoginCustomerOpts{
					Email:    createUserOpts.Email,
					Password: createUserOpts.Password,
				}

				user, _ := tt.fields.App.User.CreateUser(&createUserOpts)

				customer := model.GetRandomCustomer()
				customer.UserID = user.ID
				res, _ := tt.fields.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customer)
				customer.ID = res.InsertedID.(primitive.ObjectID)

				claim := &auth.UserClaim{
					ID:            user.ID.Hex(),
					CustomerID:    customer.ID.Hex(),
					Role:          model.UserRole,
					Type:          user.Type,
					FullName:      customer.FullName,
					DOB:           customer.DOB,
					Email:         user.Email,
					ProfileImage:  customer.ProfileImage,
					Gender:        *customer.Gender,
					EmailVerified: false,
					PhoneVerified: false,
				}

				tt.want = claim
			},
			validate: func(t *testing.T, tt *TC, resp auth.Claim) {
				w := tt.want.(*auth.UserClaim)
				g := resp.(*auth.UserClaim)
				assert.WithinDuration(t, w.DOB, g.DOB, 10*time.Millisecond)
				w.DOB = time.Time{}
				g.DOB = time.Time{}
				assert.Equal(t, w, g)
			},
		},
		{
			name: "[Error] logging with via social account",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				socialLoginOpts := schema.LoginWithSocial{
					Type:     model.CreatedViaFacebook,
					Email:    faker.Internet().SafeEmail(),
					FullName: faker.Name().Name(),
					ProfileImage: &schema.Img{
						SRC: faker.Avatar().Url("png", 50, 50),
					},
				}
				tt.args.opts = &schema.EmailLoginCustomerOpts{
					Email:    socialLoginOpts.Email,
					Password: faker.RandomString(6),
				}

				user, _ := tt.fields.App.User.LoginWithSocial(&socialLoginOpts)
				uc := user.(*auth.UserClaim)
				id, _ := primitive.ObjectIDFromHex(uc.ID)
				customer := model.Customer{
					UserID:   id,
					FullName: socialLoginOpts.FullName,
					ProfileImage: &model.IMG{
						SRC:    socialLoginOpts.ProfileImage.SRC,
						Width:  50,
						Height: 50,
					},
				}
				customer.UserID = id
				res, _ := tt.fields.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customer)
				customer.ID = res.InsertedID.(primitive.ObjectID)

				tt.err = errors.New("invalid password")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			got, err := ui.EmailLoginCustomerUser(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.EmailLoginCustomerUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestUserImpl_ForgotPassword(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.ForgotPasswordOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     bool
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, bool)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)
				tt.args.opts = &schema.ForgotPasswordOpts{
					Email: res.Email,
				}
				tt.want = true
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"email": tt.args.opts.Email}).Decode(&user)
				assert.Len(t, user.PasswordResetCode, 6)
			},
		},
		{
			name: "[Error] Invalid email",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)
				tt.args.opts = &schema.ForgotPasswordOpts{
					Email: faker.Internet().FreeEmail(),
				}
				tt.err = errors.Errorf("user with email:%s not found", tt.args.opts.Email)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			got, err := ui.ForgotPassword(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.ForgotPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestUserImpl_ResetPassword(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.ResetPasswordOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     bool
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, bool)
	}

	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				_, _ = tt.fields.App.User.ForgotPassword(&schema.ForgotPasswordOpts{Email: res.Email})

				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)
				tt.args.opts = &schema.ResetPasswordOpts{
					Email:    res.Email,
					OTP:      user.PasswordResetCode,
					Password: faker.Internet().Password(6, 10),
				}
				tt.args.opts.ConfirmPassword = tt.args.opts.Password
				tt.want = true
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"email": tt.args.opts.Email}).Decode(&user)
				assert.Len(t, user.PasswordResetCode, 0)
				assert.True(t, CheckPasswordHash(tt.args.opts.Password, user.Password))
			},
		},
		{
			name: "[Error] invalid otp",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				_, _ = tt.fields.App.User.ForgotPassword(&schema.ForgotPasswordOpts{Email: res.Email})

				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)
				tt.args.opts = &schema.ResetPasswordOpts{
					Email:    res.Email,
					OTP:      faker.Letterify("######"),
					Password: faker.Internet().Password(6, 10),
				}
				tt.args.opts.ConfirmPassword = tt.args.opts.Password
				tt.err = errors.New("invalid otp")
			},
			wantErr: true,
		},
		{
			name: "[Error] invalid email",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				createOpts := schema.GetRandomCreateUserOpts()
				res, _ := tt.fields.App.User.CreateUser(createOpts)
				_, _ = tt.fields.App.User.ForgotPassword(&schema.ForgotPasswordOpts{Email: res.Email})

				var user model.User
				_ = tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": res.ID}).Decode(&user)
				tt.args.opts = &schema.ResetPasswordOpts{
					Email:    faker.Internet().SafeEmail(),
					OTP:      faker.Letterify("######"),
					Password: faker.Internet().Password(6, 10),
				}
				tt.args.opts.ConfirmPassword = tt.args.opts.Password
				tt.err = errors.Errorf("user with email:%s not found: mongo: no documents in result", tt.args.opts.Email)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			got, err := ui.ResetPassword(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestUserImpl_MobileLoginCustomerUser(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.MobileLoginCustomerUserOpts
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     auth.Claim
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, auth.Claim)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				userDoc := model.User{
					Role: model.UserRole,
					PhoneNo: &model.PhoneNumber{
						Prefix: faker.PhoneNumber().SubscriberNumber(3),
						Number: faker.PhoneNumber().SubscriberNumber(10),
					},
					LoginOTP: &model.LoginOTP{
						Type:      model.PhoneLoginOTPType,
						OTP:       faker.Numerify("######"),
						CreatedAt: time.Now().UTC(),
					},
				}
				res, _ := tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), userDoc)
				userDoc.ID = res.InsertedID.(primitive.ObjectID)

				customerDoc := model.GetRandomCustomer()
				customerDoc.UserID = res.InsertedID.(primitive.ObjectID)
				res1, _ := tt.fields.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customerDoc)
				customerDoc.ID = res1.InsertedID.(primitive.ObjectID)

				tt.args.opts = &schema.MobileLoginCustomerUserOpts{
					PhoneNo: &schema.PhoneNoOpts{
						Prefix: userDoc.PhoneNo.Prefix,
						Number: userDoc.PhoneNo.Number,
					},
					OTP: userDoc.LoginOTP.OTP,
				}
				claim := &auth.UserClaim{
					ID:           userDoc.ID.Hex(),
					CustomerID:   customerDoc.ID.Hex(),
					Type:         userDoc.Type,
					Role:         userDoc.Role,
					Email:        userDoc.Email,
					PhoneNo:      userDoc.PhoneNo,
					CreatedVia:   userDoc.CreatedVia,
					FullName:     customerDoc.FullName,
					Gender:       *customerDoc.Gender,
					DOB:          customerDoc.DOB,
					ProfileImage: customerDoc.ProfileImage,
				}
				if !userDoc.EmailVerifiedAt.IsZero() {
					claim.EmailVerified = true
				}
				if !userDoc.PhoneVerifiedAt.IsZero() {
					claim.PhoneVerified = true
				}
				tt.want = claim
			},
			validate: func(t *testing.T, tt *TC, got auth.Claim) {
				w := tt.want.(*auth.UserClaim)
				g := got.(*auth.UserClaim)
				assert.WithinDuration(t, w.DOB, g.DOB, 10*time.Millisecond)
				w.DOB = time.Time{}
				g.DOB = time.Time{}
				assert.Equal(t, w, g)
			},
		},
		{
			name: "[Error] expired otp",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				userDoc := model.User{
					Role: model.UserRole,
					PhoneNo: &model.PhoneNumber{
						Prefix: faker.PhoneNumber().SubscriberNumber(3),
						Number: faker.PhoneNumber().SubscriberNumber(10),
					},
					LoginOTP: &model.LoginOTP{
						Type:      model.PhoneLoginOTPType,
						OTP:       faker.Numerify("######"),
						CreatedAt: time.Now().UTC().Add(-(20 * time.Minute)),
					},
				}
				res, _ := tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), userDoc)
				userDoc.ID = res.InsertedID.(primitive.ObjectID)

				customerDoc := model.GetRandomCustomer()
				customerDoc.UserID = res.InsertedID.(primitive.ObjectID)
				res1, _ := tt.fields.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customerDoc)
				customerDoc.ID = res1.InsertedID.(primitive.ObjectID)

				tt.args.opts = &schema.MobileLoginCustomerUserOpts{
					PhoneNo: &schema.PhoneNoOpts{
						Prefix: userDoc.PhoneNo.Prefix,
						Number: userDoc.PhoneNo.Number,
					},
					OTP: userDoc.LoginOTP.OTP,
				}
				tt.err = errors.New("otp expired")
			},
			wantErr: true,
		},
		{
			name: "[Error] invalid otp",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				userDoc := model.User{
					Role: model.UserRole,
					PhoneNo: &model.PhoneNumber{
						Prefix: faker.PhoneNumber().SubscriberNumber(3),
						Number: faker.PhoneNumber().SubscriberNumber(10),
					},
					LoginOTP: &model.LoginOTP{
						Type:      model.PhoneLoginOTPType,
						OTP:       faker.Numerify("######"),
						CreatedAt: time.Now().UTC(),
					},
				}
				res, _ := tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), userDoc)
				userDoc.ID = res.InsertedID.(primitive.ObjectID)

				customerDoc := model.GetRandomCustomer()
				customerDoc.UserID = res.InsertedID.(primitive.ObjectID)
				res1, _ := tt.fields.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customerDoc)
				customerDoc.ID = res1.InsertedID.(primitive.ObjectID)

				tt.args.opts = &schema.MobileLoginCustomerUserOpts{
					PhoneNo: &schema.PhoneNoOpts{
						Prefix: userDoc.PhoneNo.Prefix,
						Number: userDoc.PhoneNo.Number,
					},
					OTP: faker.Numerify("######"),
				}
				tt.err = errors.New("invalid otp")
			},
			wantErr: true,
		},
		{
			name: "[Error] phone number not found",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				userDoc := model.User{
					Role: model.UserRole,
					PhoneNo: &model.PhoneNumber{
						Prefix: faker.PhoneNumber().SubscriberNumber(3),
						Number: faker.PhoneNumber().SubscriberNumber(10),
					},
					LoginOTP: &model.LoginOTP{
						Type:      model.PhoneLoginOTPType,
						OTP:       faker.Numerify("######"),
						CreatedAt: time.Now().UTC(),
					},
				}
				res, _ := tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), userDoc)
				userDoc.ID = res.InsertedID.(primitive.ObjectID)

				customerDoc := model.GetRandomCustomer()
				customerDoc.UserID = res.InsertedID.(primitive.ObjectID)
				res1, _ := tt.fields.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customerDoc)
				customerDoc.ID = res1.InsertedID.(primitive.ObjectID)

				tt.args.opts = &schema.MobileLoginCustomerUserOpts{
					PhoneNo: &schema.PhoneNoOpts{
						Prefix: userDoc.PhoneNo.Prefix,
						Number: faker.PhoneNumber().SubscriberNumber(10),
					},
					OTP: faker.Numerify("######"),
				}
				tt.err = errors.Errorf("user with phone:%s%s not found: mongo: no documents in result", tt.args.opts.PhoneNo.Prefix, tt.args.opts.PhoneNo.Number)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			got, err := ui.MobileLoginCustomerUser(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.MobileLoginCustomerUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestUserImpl_GenerateMobileLoginOTP(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.GenerateMobileLoginOTPOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       bool
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockSNS)
		validate   func(*testing.T, *TC, bool)
	}

	tests := []TC{
		{
			name: "[Ok] user does not exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.GenerateMobileLoginOTPOpts{
					PhoneNo: &schema.PhoneNoOpts{
						Prefix: faker.Numerify("+##"),
						Number: faker.PhoneNumber().SubscriberNumber(10),
					},
				}
				tt.args.opts = &opts
			},
			buildStubs: func(tt *TC, m *mock.MockSNS) {
				resp := &sns.PublishOutput{}
				m.EXPECT().Publish(gomock.Any()).Times(1).Return(resp, nil)
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				var userDoc model.User
				tt.fields.DB.Collection(model.UserColl).FindOne(
					context.TODO(),
					bson.M{"phone_no": bson.M{"prefix": tt.args.opts.PhoneNo.Prefix, "number": tt.args.opts.PhoneNo.Number}},
				).Decode(&userDoc)

				assert.NotNil(t, userDoc.LoginOTP)
				assert.Len(t, userDoc.LoginOTP.OTP, 6)
				assert.Equal(t, model.PhoneLoginOTPType, userDoc.LoginOTP.Type)
				assert.WithinDuration(t, time.Now().UTC(), userDoc.LoginOTP.CreatedAt, 50*time.Millisecond)
				assert.Equal(t, model.UserRole, userDoc.Role)
				assert.Equal(t, model.CustomerType, userDoc.Type)

				var customerDoc model.Customer
				tt.fields.DB.Collection(model.CustomerColl).FindOne(
					context.TODO(),
					bson.M{"user_id": userDoc.ID},
				).Decode(&customerDoc)
				assert.False(t, customerDoc.ID.IsZero())
				assert.False(t, customerDoc.UserID.IsZero())
			},
		},
		{
			name: "[Ok] user exists",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				userDoc := model.User{
					PhoneNo: &model.PhoneNumber{
						Prefix: faker.Letterify("+##"),
						Number: faker.PhoneNumber().SubscriberNumber(10),
					},
				}
				tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), userDoc)

				opts := schema.GenerateMobileLoginOTPOpts{
					PhoneNo: &schema.PhoneNoOpts{
						Prefix: userDoc.PhoneNo.Prefix,
						Number: userDoc.PhoneNo.Number,
					},
				}
				tt.args.opts = &opts
			},
			buildStubs: func(tt *TC, m *mock.MockSNS) {
				resp := &sns.PublishOutput{}
				m.EXPECT().Publish(gomock.Any()).Times(1).Return(resp, nil)
			},
			validate: func(t *testing.T, tt *TC, resp bool) {
				var userDoc model.User
				tt.fields.DB.Collection(model.UserColl).FindOne(
					context.TODO(),
					bson.M{"phone_no": bson.M{"prefix": tt.args.opts.PhoneNo.Prefix, "number": tt.args.opts.PhoneNo.Number}},
				).Decode(&userDoc)

				assert.NotNil(t, userDoc.LoginOTP)
				assert.Len(t, userDoc.LoginOTP.OTP, 6)
				assert.Equal(t, model.PhoneLoginOTPType, userDoc.LoginOTP.Type)
				assert.WithinDuration(t, time.Now().UTC(), userDoc.LoginOTP.CreatedAt, 50*time.Millisecond)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSNS := mock.NewMockSNS(ctrl)
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.fields.App.SNS = mockSNS
			tt.prepare(&tt)
			tt.buildStubs(&tt, mockSNS)
			got, err := ui.GenerateMobileLoginOTP(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.GenerateMobileLoginOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestUserImpl_LoginWithSocial(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.LoginWithSocial
	}

	type TC struct {
		name     string
		fields   fields
		args     args
		want     auth.Claim
		wantErr  bool
		err      error
		prepare  func(*TC)
		validate func(*testing.T, *TC, auth.Claim)
	}
	tests := []TC{
		{
			name: "[Ok] first time login (user doesn't exist)",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName + "_1"),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.LoginWithSocial{
					Type:     model.CreatedViaGoogle,
					Email:    faker.Internet().FreeEmail(),
					FullName: faker.Name().Name(),
					ProfileImage: &schema.Img{
						SRC: faker.Avatar().Url("png", 50, 50),
					},
				}
				tt.args.opts = &opts
			},
			validate: func(t *testing.T, tt *TC, got auth.Claim) {
				g := got.(*auth.UserClaim)
				id, _ := primitive.ObjectIDFromHex(g.ID)
				cid, _ := primitive.ObjectIDFromHex(g.CustomerID)
				var userDoc model.User
				err := tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&userDoc)
				assert.Nil(t, err)
				assert.WithinDuration(t, time.Now().UTC(), userDoc.CreatedAt, 2*time.Second)

				var customerDoc model.Customer
				err = tt.fields.DB.Collection(model.CustomerColl).FindOne(context.TODO(), bson.M{"_id": cid}).Decode(&customerDoc)
				assert.Nil(t, err)
				assert.WithinDuration(t, time.Now().UTC(), userDoc.CreatedAt, 2*time.Second)

				assert.False(t, id.IsZero())
				assert.False(t, cid.IsZero())
				assert.Equal(t, tt.args.opts.Type, g.CreatedVia)
				assert.Equal(t, model.UserRole, g.Role)
				assert.Equal(t, model.CustomerType, g.Type)

				assert.Equal(t, customerDoc.UserID, userDoc.ID)
				assert.Equal(t, tt.args.opts.FullName, customerDoc.FullName)
				assert.Equal(t, tt.args.opts.ProfileImage.SRC, customerDoc.ProfileImage.SRC)
				assert.Equal(t, 50, customerDoc.ProfileImage.Height)
				assert.Equal(t, 50, customerDoc.ProfileImage.Width)

				userCount, _ := tt.fields.DB.Collection(model.UserColl).CountDocuments(context.TODO(), bson.M{})
				customerCount, _ := tt.fields.DB.Collection(model.CustomerColl).CountDocuments(context.TODO(), bson.M{})

				assert.Equal(t, int64(1), userCount)
				assert.Equal(t, int64(1), customerCount)
			},
		},
		{
			name: "[Ok] user exists but 2nd one is different user",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName + "_2"),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.LoginWithSocial{
					Type:     model.CreatedViaGoogle,
					Email:    faker.Internet().FreeEmail(),
					FullName: faker.Name().Name(),
					ProfileImage: &schema.Img{
						SRC: faker.Avatar().Url("png", 50, 50),
					},
				}
				tt.args.opts = &opts

				user := model.User{
					Role:       model.UserRole,
					Type:       model.CustomerType,
					Email:      faker.Internet().FreeEmail(),
					CreatedVia: model.CreatedViaGoogle,
					CreatedAt:  time.Now().UTC(),
				}
				res, _ := tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), user)
				user.ID = res.InsertedID.(primitive.ObjectID)

				customer := model.Customer{
					UserID:    user.ID,
					CreatedAt: time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customer)
			},
			validate: func(t *testing.T, tt *TC, got auth.Claim) {
				g := got.(*auth.UserClaim)
				id, _ := primitive.ObjectIDFromHex(g.ID)
				cid, _ := primitive.ObjectIDFromHex(g.CustomerID)
				var userDoc model.User
				err := tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&userDoc)
				assert.Nil(t, err)
				assert.WithinDuration(t, time.Now().UTC(), userDoc.CreatedAt, 2*time.Second)

				var customerDoc model.Customer
				err = tt.fields.DB.Collection(model.CustomerColl).FindOne(context.TODO(), bson.M{"_id": cid}).Decode(&customerDoc)
				assert.Nil(t, err)
				assert.WithinDuration(t, time.Now().UTC(), userDoc.CreatedAt, 2*time.Second)

				assert.False(t, id.IsZero())
				assert.False(t, cid.IsZero())
				assert.Equal(t, tt.args.opts.Type, g.CreatedVia)
				assert.Equal(t, model.UserRole, g.Role)
				assert.Equal(t, model.CustomerType, g.Type)

				assert.Equal(t, customerDoc.UserID, userDoc.ID)
				assert.Equal(t, tt.args.opts.FullName, customerDoc.FullName)
				assert.Equal(t, tt.args.opts.ProfileImage.SRC, customerDoc.ProfileImage.SRC)
				assert.Equal(t, 50, customerDoc.ProfileImage.Height)
				assert.Equal(t, 50, customerDoc.ProfileImage.Width)

				userCount, _ := tt.fields.DB.Collection(model.UserColl).CountDocuments(context.TODO(), bson.M{})
				customerCount, _ := tt.fields.DB.Collection(model.CustomerColl).CountDocuments(context.TODO(), bson.M{})

				assert.Equal(t, int64(2), userCount)
				assert.Equal(t, int64(2), customerCount)
			},
		},
		{
			name: "[Ok] user logging in 2nd time",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName + "_3"),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.LoginWithSocial{
					Type:     model.CreatedViaGoogle,
					Email:    faker.Internet().FreeEmail(),
					FullName: faker.Name().Name(),
					ProfileImage: &schema.Img{
						SRC: faker.Avatar().Url("png", 50, 50),
					},
				}
				tt.args.opts = &opts

				user := model.User{
					Role:       model.UserRole,
					Type:       model.CustomerType,
					Email:      opts.Email,
					CreatedVia: model.CreatedViaGoogle,
					CreatedAt:  time.Now().UTC(),
				}
				res, _ := tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), user)
				user.ID = res.InsertedID.(primitive.ObjectID)

				customer := model.Customer{
					UserID:    user.ID,
					CreatedAt: time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customer)
			},
			validate: func(t *testing.T, tt *TC, got auth.Claim) {
				g := got.(*auth.UserClaim)
				id, _ := primitive.ObjectIDFromHex(g.ID)
				cid, _ := primitive.ObjectIDFromHex(g.CustomerID)
				var userDoc model.User
				err := tt.fields.DB.Collection(model.UserColl).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&userDoc)
				assert.Nil(t, err)
				assert.WithinDuration(t, time.Now().UTC(), userDoc.CreatedAt, 2*time.Second)

				var customerDoc model.Customer
				err = tt.fields.DB.Collection(model.CustomerColl).FindOne(context.TODO(), bson.M{"_id": cid}).Decode(&customerDoc)
				assert.Nil(t, err)
				assert.WithinDuration(t, time.Now().UTC(), userDoc.CreatedAt, 2*time.Second)

				assert.False(t, id.IsZero())
				assert.False(t, cid.IsZero())
				assert.Equal(t, tt.args.opts.Type, g.CreatedVia)
				assert.Equal(t, model.UserRole, g.Role)
				assert.Equal(t, model.CustomerType, g.Type)

				assert.Equal(t, customerDoc.UserID, userDoc.ID)
				assert.Equal(t, tt.args.opts.FullName, customerDoc.FullName)
				assert.Equal(t, tt.args.opts.ProfileImage.SRC, customerDoc.ProfileImage.SRC)
				assert.Equal(t, 50, customerDoc.ProfileImage.Height)
				assert.Equal(t, 50, customerDoc.ProfileImage.Width)

				userCount, _ := tt.fields.DB.Collection(model.UserColl).CountDocuments(context.TODO(), bson.M{})
				customerCount, _ := tt.fields.DB.Collection(model.CustomerColl).CountDocuments(context.TODO(), bson.M{})

				assert.Equal(t, int64(1), userCount)
				assert.Equal(t, int64(1), customerCount)
			},
		},
		{
			name: "[Error] user logging in 2nd time with different social login but same email",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.LoginWithSocial{
					Type:     model.CreatedViaFacebook,
					Email:    faker.Internet().FreeEmail(),
					FullName: faker.Name().Name(),
					ProfileImage: &schema.Img{
						SRC: faker.Avatar().Url("png", 50, 50),
					},
				}
				tt.args.opts = &opts

				user := model.User{
					Role:       model.UserRole,
					Type:       model.CustomerType,
					Email:      opts.Email,
					CreatedVia: model.CreatedViaGoogle,
					CreatedAt:  time.Now().UTC(),
				}
				res, _ := tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), user)
				user.ID = res.InsertedID.(primitive.ObjectID)

				customer := model.Customer{
					UserID:    user.ID,
					CreatedAt: time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customer)

				tt.err = errors.New("cannot use facebook login: this account was created via google")
			},
			wantErr: true,
		},
		{
			name: "[Error] user logging in 2nd time with different social login but same email",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {
				opts := schema.LoginWithSocial{
					Type:     model.CreatedViaGoogle,
					Email:    faker.Internet().FreeEmail(),
					FullName: faker.Name().Name(),
					ProfileImage: &schema.Img{
						SRC: faker.Avatar().Url("png", 50, 50),
					},
				}
				tt.args.opts = &opts

				user := model.User{
					Role:       model.UserRole,
					Type:       model.CustomerType,
					Email:      opts.Email,
					CreatedVia: model.CreatedViaFacebook,
					CreatedAt:  time.Now().UTC(),
				}
				res, _ := tt.fields.DB.Collection(model.UserColl).InsertOne(context.TODO(), user)
				user.ID = res.InsertedID.(primitive.ObjectID)

				customer := model.Customer{
					UserID:    user.ID,
					CreatedAt: time.Now().UTC(),
				}
				tt.fields.DB.Collection(model.CustomerColl).InsertOne(context.TODO(), customer)

				tt.err = errors.New("cannot use google login: this account was created via facebook")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := &UserImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.User = ui
			tt.prepare(&tt)
			got, err := ui.LoginWithSocial(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.LoginWithSocial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("UserImpl.MobileLoginCustomerUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
			if tt.wantErr {
				assert.Nil(t, got)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}
