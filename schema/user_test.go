package schema

import (
	"encoding/json"
	"go-app/model"
	"go-app/server/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreateUserOpts
	}{
		{
			name: "[Ok] Create user with email id",
			json: string(`{
				"type": "customer",
				"email": "vasu@hypd.in",
				"password": "abcd1234",
				"confirm_password": "abcd1234"
			}`),
			want: CreateUserOpts{
				Type:            model.CustomerType,
				Email:           "vasu@hypd.in",
				Password:        "abcd1234",
				ConfirmPassword: "abcd1234",
			},
		},
		{
			name: "[Ok] Create user with mobile id",
			json: string(`{
				"type": "customer",
				"phone_no": {
					"prefix": "+91",
					"number": "9988998899"
				},
				"password": "abcd1234",
				"confirm_password": "abcd1234"
			}`),
			want: CreateUserOpts{
				Type: model.CustomerType,
				MobileNo: &PhoneNoOpts{
					Prefix: "+91",
					Number: "9988998899",
				},
				Password:        "abcd1234",
				ConfirmPassword: "abcd1234",
			},
		},
		{
			name: "[error] Without Email and MobileNo",
			json: string(`{
				"type": "customer",
				"password": "abcd1234",
				"confirm_password": "abcd1234"
			}`),
			wantErr: true,
			err:     []string{"Key: 'CreateUserOpts.email' Error:Field validation for 'email' failed on the 'required_without=MobileNo|email' tag"},
		},
		{
			name: "[Ok] Invalid Password",
			json: string(`{
				"type": "customer",
				"phone_no": {
					"prefix": "+91",
					"number": "9988998899"
				},
				"password": "abcd",
				"confirm_password": "abcd"
			}`),
			wantErr: true,
			err:     []string{"password must be at least 6 characters in length"},
		},
		{
			name: "[Ok] Different Password and ConfirmPassword",
			json: string(`{
				"type": "customer",
				"phone_no": {
					"prefix": "+91",
					"number": "9988998899"
				},
				"password": "abcd1234",
				"confirm_password": "abcd12345"
			}`),
			wantErr: true,
			err:     []string{"confirm_password must be equal to Password"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc CreateUserOpts
			err := json.Unmarshal([]byte(tt.json), &sc)
			assert.Nil(t, err)
			errs := tv.Validate(&sc)
			if tt.wantErr {
				assert.Len(t, errs, len(tt.err))
				assert.Equal(t, errs[0].Error(), tt.err[0])
			}
			if !tt.wantErr {
				assert.Len(t, errs, 0)
				assert.Equal(t, tt.want, sc)
			}
		})
	}
}

func TestEmailLoginUserOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    EmailLoginCustomerOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"email": "vasu@hypd.in",
				"password": "abcd1234"
			}`),
			want: EmailLoginCustomerOpts{
				Email:    "vasu@hypd.in",
				Password: "abcd1234",
			},
		},
		{
			name: "[Error] Without email",
			json: string(`{
				"password": "abcd1234"
			}`),
			wantErr: true,
			err:     []string{"email is a required field"},
		},
		{
			name: "[Error] Without Password",
			json: string(`{
				"email": "vasu@hypd.in"
			}`),
			wantErr: true,
			err:     []string{"password is a required field"},
		},
		{
			name: "[Error] Empty Password",
			json: string(`{
				"email": "vasu@hypd.in",
				"password": ""
			}`),
			wantErr: true,
			err:     []string{"password is a required field"},
		},
		{
			name: "[Error] Password less than 6 characters",
			json: string(`{
				"email": "vasu@hypd.in",
				"password": "iam"
			}`),
			wantErr: true,
			err:     []string{"password must be at least 6 characters in length"},
		},
		{
			name: "[Error] Invalid email",
			json: string(`{
				"email": "vasuhypd.in",
				"password": "iamiamiam"
			}`),
			wantErr: true,
			err:     []string{"email must be a valid email address"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc EmailLoginCustomerOpts
			err := json.Unmarshal([]byte(tt.json), &sc)
			assert.Nil(t, err)
			errs := tv.Validate(&sc)
			if tt.wantErr {
				assert.Len(t, errs, len(tt.err))
				assert.Equal(t, errs[0].Error(), tt.err[0])
			}
			if !tt.wantErr {
				assert.Len(t, errs, 0)
				assert.Equal(t, tt.want, sc)
			}
		})
	}
}
