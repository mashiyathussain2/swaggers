package app

import (
	"context"
	"go-app/mock"
	"go-app/model"
	"go-app/schema"
	"go-app/server/auth"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"syreclabs.com/go/faker"
)

func TestCustomerImpl_SignUp(t *testing.T) {
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
		name       string
		fields     fields
		args       args
		want       auth.Claim
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockUser)
		validate   func(*testing.T, *TC, auth.Claim)
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
				opts := schema.GetRandomCreateUserOpts()
				tt.args.opts = opts
			},
			buildStubs: func(tt *TC, m *mock.MockUser) {
				resp := schema.CreateUserResp{
					ID:    primitive.NewObjectID(),
					Type:  model.CustomerType,
					Email: tt.args.opts.Email,
				}
				m.EXPECT().CreateUser(tt.args.opts).Times(1).Return(&resp, nil)
				tt.want = &auth.UserClaim{
					ID:      resp.ID.Hex(),
					Type:    resp.Type,
					Email:   resp.Email,
					PhoneNo: resp.PhoneNo,
					Role:    model.UserRole,
				}
			},
			validate: func(t *testing.T, tt *TC, resp auth.Claim) {
				w := tt.want.(*auth.UserClaim)
				g := resp.(*auth.UserClaim)
				assert.WithinDuration(t, w.DOB, g.DOB, 10*time.Millisecond)
				w.DOB = time.Time{}
				g.DOB = time.Time{}

				var customer *model.Customer
				cID, _ := primitive.ObjectIDFromHex(g.CustomerID)
				f := bson.M{"_id": cID}
				err := tt.fields.DB.Collection(model.CustomerColl).FindOne(context.TODO(), f).Decode(&customer)
				assert.NotNil(t, customer)
				assert.Nil(t, err)
				w.CustomerID = customer.ID.Hex()
				assert.Equal(t, w, g)
				assert.WithinDuration(t, time.Now().UTC(), customer.CreatedAt, 100*time.Millisecond)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userMock := mock.NewMockUser(ctrl)

			ci := &CustomerImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Customer = ci
			tt.fields.App.User = userMock
			tt.prepare(&tt)
			tt.buildStubs(&tt, userMock)
			got, err := ci.SignUp(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomerImpl.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validate(t, &tt, got)
			}
		})
	}
}

func TestCustomerImpl_UpdateCustomer(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		createOpts         *schema.CreateUserOpts
		createCustomerResp *auth.UserClaim
		opts               *schema.UpdateCustomerOpts
	}

	type TC struct {
		name       string
		fields     fields
		args       args
		want       *schema.GetCustomerInfoResp
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(*TC, *mock.MockUser)
		validate   func(*testing.T, *TC, *schema.GetCustomerInfoResp)
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

				dob, _ := time.Parse(time.RFC3339, "1996-10-21T00:00:00+00:00")
				userid, _ := primitive.ObjectIDFromHex(tt.args.createCustomerResp.ID)
				id, _ := primitive.ObjectIDFromHex(tt.args.createCustomerResp.CustomerID)
				tt.args.opts = &schema.UpdateCustomerOpts{
					ID:       id,
					FullName: faker.Name().Name(),
					DOB:      dob,
					Gender:   faker.RandomChoice([]string{model.Male, model.Female, model.Others}),
					ProfileImage: &schema.Img{
						SRC: faker.Avatar().Url("png", 100, 200),
					},
				}
				tt.want = &schema.GetCustomerInfoResp{
					ID:       id,
					UserID:   userid,
					FullName: tt.args.opts.FullName,
					DOB:      tt.args.opts.DOB,
					Gender:   &tt.args.opts.Gender,
					ProfileImage: &model.IMG{
						SRC:    tt.args.opts.ProfileImage.SRC,
						Width:  100,
						Height: 200,
					},
				}
			},
			buildStubs: func(tt *TC, m *mock.MockUser) {
				opts := schema.GetRandomCreateUserOpts()
				tt.args.createOpts = opts
				resp := schema.CreateUserResp{
					ID:    primitive.NewObjectID(),
					Type:  model.CustomerType,
					Email: tt.args.createOpts.Email,
				}
				m.EXPECT().CreateUser(tt.args.createOpts).Times(1).Return(&resp, nil)
				resp1, _ := tt.fields.App.Customer.SignUp(opts)
				tt.args.createCustomerResp = resp1.(*auth.UserClaim)
			},
			validate: func(t *testing.T, tt *TC, got *schema.GetCustomerInfoResp) {
				assert.WithinDuration(t, tt.want.DOB, got.DOB, 10*time.Millisecond)
				tt.want.DOB = time.Time{}
				got.DOB = time.Time{}
				assert.Equal(t, tt.want, got)
			},
		},
		{
			name: "[Error] no fields to update",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.opts = &schema.UpdateCustomerOpts{
					ID: primitive.NewObjectID(),
				}
				tt.err = errors.New("no field update found for customer")
			},
			buildStubs: func(tt *TC, m *mock.MockUser) {

			},
		},
		{
			name: "[Error] invalid gender value",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.opts = &schema.UpdateCustomerOpts{
					ID:     primitive.NewObjectID(),
					Gender: "X",
				}
				tt.err = errors.Errorf("%s is invalid gender value", tt.args.opts.Gender)
			},
			buildStubs: func(tt *TC, m *mock.MockUser) {

			},
		},
		{
			name: "[Error] invalid image url",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args:    args{},
			wantErr: true,
			prepare: func(tt *TC) {
				tt.args.opts = &schema.UpdateCustomerOpts{
					ID: primitive.NewObjectID(),
					ProfileImage: &schema.Img{
						SRC: "abc",
					},
				}
				tt.err = errors.New("Get \"abc\": unsupported protocol scheme \"\"")
			},
			buildStubs: func(tt *TC, m *mock.MockUser) {

			},
		},
		{
			name: "[Ok] Only profile image url",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.UserConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			prepare: func(tt *TC) {

				userid, _ := primitive.ObjectIDFromHex(tt.args.createCustomerResp.ID)
				id, _ := primitive.ObjectIDFromHex(tt.args.createCustomerResp.CustomerID)
				tt.args.opts = &schema.UpdateCustomerOpts{
					ID: id,
					ProfileImage: &schema.Img{
						SRC: faker.Avatar().Url("png", 100, 200),
					},
				}
				tt.want = &schema.GetCustomerInfoResp{
					ID:       id,
					UserID:   userid,
					FullName: tt.args.opts.FullName,
					DOB:      tt.args.opts.DOB,
					Gender:   nil,
					ProfileImage: &model.IMG{
						SRC:    tt.args.opts.ProfileImage.SRC,
						Width:  100,
						Height: 200,
					},
				}
			},
			buildStubs: func(tt *TC, m *mock.MockUser) {
				opts := schema.GetRandomCreateUserOpts()
				tt.args.createOpts = opts
				resp := schema.CreateUserResp{
					ID:    primitive.NewObjectID(),
					Type:  model.CustomerType,
					Email: tt.args.createOpts.Email,
				}
				m.EXPECT().CreateUser(tt.args.createOpts).Times(1).Return(&resp, nil)
				resp1, _ := tt.fields.App.Customer.SignUp(opts)
				tt.args.createCustomerResp = resp1.(*auth.UserClaim)
			},
			validate: func(t *testing.T, tt *TC, got *schema.GetCustomerInfoResp) {
				assert.WithinDuration(t, tt.want.DOB, got.DOB, 10*time.Millisecond)
				tt.want.DOB = time.Time{}
				got.DOB = time.Time{}
				assert.Equal(t, tt.want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userMock := mock.NewMockUser(ctrl)

			ci := &CustomerImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Customer = ci
			tt.fields.App.User = userMock
			tt.buildStubs(&tt, userMock)
			tt.prepare(&tt)
			got, err := ci.UpdateCustomer(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomerImpl.UpdateCustomer() error = %v, wantErr %v", err, tt.wantErr)
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
