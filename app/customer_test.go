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
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
