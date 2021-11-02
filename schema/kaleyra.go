package schema

import "go-app/model"

type SendOTPOpts struct {
	PhoneNo model.PhoneNumber `json:"to"`
	OTP     string            `json:"otp"`
}
