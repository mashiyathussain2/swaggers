package api

import (
	"fmt"
	"go-app/model"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:route  GET /me ME me
// me
//
// This endpoint will returns the updated user info stored in the token.
//
// Endpoint: /me
//
// Method: GET
//
// parameters:
// + name: cookie
//   in: header
//   description:Login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
//   in: header
//   description:Token required for successful response.
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description:Invalid User
//  200: UserClaim description: OK
func (a *API) me(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	requestCTX.SetAppResponse(requestCTX.UserClaim, http.StatusOK)
}

// swagger:route  POST /me ME updateMe
// updateMe
//
// This endpoint will returns the updated user info stored in the token.
//
// Endpoint: /me
//
// Method: POST
//
// parameters:
// + name: cookie
//   in: header
//   description:Login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
//   in: header
//   description:Token required for successful response.
//   required: true
//
// parameters:
// + name: isWeb
//   in: query
//   description: If value is set to True, token is omitted from response
//   schema:
//   type: boolean
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description:Invalid User
//  200: UserClaim description: OK
func (a *API) updateMe(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var userID primitive.ObjectID
	isWeb, _ := strconv.ParseBool(r.URL.Query().Get("isWeb"))

	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	}
	res, err := a.App.Customer.UpdateToken(userID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	token, err := a.TokenAuth.SignToken(res)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if isWeb {
		if err := a.SessionAuth.Create(token, w); err != nil {
			requestCTX.SetErr(fmt.Errorf("failed to login user: %s", err), http.StatusInternalServerError)
			return
		}
		requestCTX.SetAppResponse(map[string]interface{}{"data": res}, http.StatusOK)
		return
	}
	requestCTX.SetAppResponse(map[string]interface{}{"token": token, "data": res}, http.StatusOK)
}

func (a *API) keeperUpdateMe(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	sID, err := a.SessionAuth.GetSessionID(r)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	userSession, err := a.SessionAuth.GetToken(sID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	claim, err := a.TokenAuth.VerifyToken(userSession.Token)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(map[string]interface{}{"token": userSession.Token, "data": claim}, http.StatusOK)
}

// swagger:route  POST /user/forgot-password Password forgotPassword
// forgotPassword
//
// This endpoint will help the user to recover the password.
//
// Endpoint: /user/forgot-password
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: ForgotPasswordOpts
//     "$ref": "#/definitions/ForgotPasswordOpts"
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description:Invalid User
//  200: description: true
func (a *API) forgotPassword(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ForgotPasswordOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.User.ForgotPassword(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /user/reset-password Password resetPassword
// resetPassword
//
// This endpoint will help the user reset the password.
//
// Endpoint: /user/reset-password
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: ResetPasswordOpts
//     "$ref": "#/definitions/ResetPasswordOpts"
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description:Invalid User
//  200: description: true
func (a *API) resetPassword(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ResetPasswordOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.User.ResetPassword(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /user/auth/email/verify Verification verifyEmailAuth
// verifyEmailAuth
//
// This endpoint will verify the user email.
//
// Endpoint: /user/auth/email/verify
//
// Method: POST
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: VerifyEmailOpts
//     "$ref": "#/definitions/VerifyEmailOpts"
//   required: true
//
// parameters:
// + name: isWeb
//   in: query
//   description: If value is set to True, token is omitted from response
//   schema:
//   type: boolean
//   required: true
//
//
// parameters:
// + name: cookie
//   in: header
//   description: Customer login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
//   in: header
//   description:Token required for successful response.
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description:Invalid User
//  500: AppErr description:Failed to login user
//  200: description: true
func (a *API) verifyEmail(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.VerifyEmailOpts
	resp := make(map[string]interface{})
	var isWeb bool
	isWeb, _ = strconv.ParseBool(r.URL.Query().Get("isWeb"))

	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.User.VerifyEmail(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	if res {
		resp["email_verified"] = true
		claim := requestCTX.UserClaim.(*auth.UserClaim)
		if claim != nil {
			claim.EmailVerified = true
			claim.Email = s.Email
			token, err := a.TokenAuth.SignToken(claim)
			if err != nil {
				requestCTX.SetErr(err, http.StatusBadRequest)
				return
			}
			if isWeb {
				if err := a.SessionAuth.Create(token, w); err != nil {
					requestCTX.SetErr(fmt.Errorf("failed to login user: %s", err), http.StatusInternalServerError)
					return
				}
				requestCTX.SetAppResponse(resp, http.StatusOK)
				return
			}
			resp["token"] = token
		}
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  POST /user/auth/email/check checkEmail checkEmail
// Check User Email for Auth
//
// Endpoint: /user/auth/email/check
//
// This endpoint will check the user email exists or not.
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   description: Check user email
//   schema:
//   type: CheckEmailOpts
//     "$ref": "#/definitions/CheckEmailOpts"
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: CommonError description: Error
//  200: description: payload : true
func (a *API) checkEmail(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CheckEmailOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.User.CheckEmail(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  POST /user/auth/phone/check checkPhoneNo checkPhoneNo
// Check User Phone No for Auth
//
// Endpoint: /user/auth/phone/check
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CheckPhoneNoOpts
//     "$ref": "#/definitions/CheckPhoneNoOpts"
//   required: true
//
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: CommonError description: Error
//  200: description: payload : true
func (a *API) checkPhoneNo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CheckPhoneNoOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.User.CheckPhoneNo(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  POST /user/auth/phone/verify Verification verifyPhoneNoAuth
// verifyPhoneNoAuth
//
// This endpoint will verify user phone number.
//
// Endpoint: /user/auth/phone/verify
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: VerifyPhoneNoOpts
//     "$ref": "#/definitions/VerifyPhoneNoOpts"
//   required: true
//
// parameters:
// + name: isWeb
//   in: query
//   enum: true
//   type: bool
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: CommonError description: Error
//  200: description: payload : true
func (a *API) verifyPhoneNo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.VerifyPhoneNoOpts
	resp := make(map[string]interface{})

	var isWeb bool
	isWeb, _ = strconv.ParseBool(r.URL.Query().Get("isWeb"))

	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.User.VerifyPhoneNo(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	if res {
		resp["phone_verified"] = true
		claim := requestCTX.UserClaim.(*auth.UserClaim)
		if claim != nil {
			claim.PhoneVerified = true
			claim.PhoneNo = &model.PhoneNumber{
				Prefix: s.PhoneNo.Prefix,
				Number: s.PhoneNo.Number,
			}
			token, err := a.TokenAuth.SignToken(claim)
			if err != nil {
				requestCTX.SetErr(err, http.StatusBadRequest)
				return
			}
			if isWeb {
				if err := a.SessionAuth.Create(token, w); err != nil {
					requestCTX.SetErr(fmt.Errorf("failed to update user: %s", err), http.StatusInternalServerError)
					return
				}
				requestCTX.SetAppResponse(resp, http.StatusOK)
				return
			}
			resp["token"] = token
		}
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  POST /user/verify-email/resend Verification resendEmailVerificationCode
// resendEmailVerificationCode
//
// This endpoint will resend email verification code
//
// Endpoint: /user/verify-email/resend
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: ResendVerificationEmailOpts
//     "$ref": "#/definitions/ResendVerificationEmailOpts"
//   required: true
//
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: CommonError description: Error
//  200: description:true
func (a *API) resendEmailVerificationCode(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ResendVerificationEmailOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.User.ResendConfirmationEmail(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /customer/otp/generate login LoginViaMobileOTP
// LoginViaMobileOTP
//
// This endpoint will generate the otp when login via mobile
//
//
// Endpoint: /customer/otp/generate
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GenerateMobileLoginOTPOpts
//     "$ref": "#/definitions/GenerateMobileLoginOTPOpts"
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: CommonError description: Error
//  200: description: payload : true
func (a *API) loginViaMobileOTP(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GenerateMobileLoginOTPOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.User.GenerateMobileLoginOTP(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /customer/social/login login loginViaSocial
// loginViaSocial
//
// This endpoint will login the user via social
//
// Endpoint: /customer/social/login
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   description: Login user via social
//   schema:
//   type: LoginWithSocial
//     "$ref": "#/definitions/LoginWithSocial"
//   required: true
//
// parameters:
// + name: isWeb
//   in: query
//   description: If value is set to True, token is omitted from response
//   schema:
//   type: boolean
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: Error
//  200: SuccessfulLogin description: Success
func (a *API) loginViaSocial(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.LoginWithSocial
	var isWeb bool
	isWeb, _ = strconv.ParseBool(r.URL.Query().Get("isWeb"))
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.User.LoginWithSocial(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	token, err := a.TokenAuth.SignToken(res)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if isWeb {
		if err := a.SessionAuth.Create(token, w); err != nil {
			requestCTX.SetErr(fmt.Errorf("failed to login user: %s", err), http.StatusInternalServerError)
			return
		}
		requestCTX.SetAppResponse(true, http.StatusOK)
		return
	}
	requestCTX.SetAppResponse(token, http.StatusOK)
}

// swagger:route  POST /customer/apple/login login loginViaApple
// loginViaApple
//
// This endpoint will login the user via apple.
//
// Endpoint: /customer/apple/login
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   description: Login user via apple
//   schema:
//   type: LoginWithApple
//     "$ref": "#/definitions/LoginWithApple"
//   required: true
//
// parameters:
// + name: isWeb
//   in: query
//   description: If value is set to True, token is omitted from response
//   schema:
//   type: boolean
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: CommonError description: Error
//  200: SuccessfulLogin description: Success
func (a *API) loginViaApple(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.LoginWithApple
	var isWeb bool
	isWeb, _ = strconv.ParseBool(r.URL.Query().Get("isWeb"))
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.User.LoginWithApple(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	token, err := a.TokenAuth.SignToken(res)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if isWeb {
		if err := a.SessionAuth.Create(token, w); err != nil {
			requestCTX.SetErr(fmt.Errorf("failed to login user: %s", err), http.StatusInternalServerError)
			return
		}
		requestCTX.SetAppResponse(true, http.StatusOK)
		return
	}
	requestCTX.SetAppResponse(token, http.StatusOK)
}

// swagger:route  POST /customer/otp/confirm login confrimloginViaMobileOtp
// confrimloginViaMobileOtp
//
// This endpoint will confirm login via mobile otp.
//
// Endpoint: /customer/otp/confirm
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: MobileLoginCustomerUserOpts
//     "$ref": "#/definitions/MobileLoginCustomerUserOpts"
//   required: true
//
// parameters:
// + name: isWeb
//   in: query
//   description: If value is set to True, token is omitted from response
//   schema:
//   type: boolean
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: CommonError description: Error
//  200: SuccessfulLogin description: Success
func (a *API) confirmLoginViaMobileOTP(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.MobileLoginCustomerUserOpts
	var isWeb bool
	isWeb, _ = strconv.ParseBool(r.URL.Query().Get("isWeb"))
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.User.MobileLoginCustomerUser(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	token, err := a.TokenAuth.SignToken(res)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if isWeb {
		if err := a.SessionAuth.Create(token, w); err != nil {
			requestCTX.SetErr(fmt.Errorf("failed to login user: %s", err), http.StatusInternalServerError)
			return
		}
		requestCTX.SetAppResponse(true, http.StatusOK)
		return
	}
	requestCTX.SetAppResponse(token, http.StatusOK)
}

func (a *API) getUserInfoByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetUserInfoByIDOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.User.GetUserInfoByID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) keeperLogin(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	// cookie, err := r.Cookie("session")
	// if err != nil {
	// 	requestCTX.SetErr(err, http.StatusBadRequest)
	// 	return
	// }
	url := a.App.KeeperUser.Login()
	requestCTX.SetRedirectResponse(url, http.StatusTemporaryRedirect)
}

func (a *API) keeperLoginCallback(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	claim, err := a.App.KeeperUser.Callback(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	fmt.Println(1)
	token, err := a.TokenAuth.SignKeeperToken(claim)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	sid, err := a.SessionAuth.CreateAndReturn(token, w)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(claim.(*auth.UserClaim).ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	err = a.App.KeeperUser.AddNewSessionID(id, sid)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	redirectURL := fmt.Sprintf("%s?token=%s", a.Config.KeeperLoginRedirectURL, token)
	fmt.Println(redirectURL)
	requestCTX.SetRedirectResponse(redirectURL, http.StatusPermanentRedirect)
}

// swagger:route  GET /user/auth/logoutt logout logoutUser
// logoutUser
//
// This endpoint will logout the user.
//
// Endpoint: /user/auth/logoutt
//
// Method: GET
//
//
// parameters:
// + name: cookie
//   in: header
//   description: Login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
//   in: header
//   description:Token required for successful response.
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description:Invalid User
//  200: description: true
func (a *API) logoutUser(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	if requestCTX.UserClaim != nil {
		a.SessionAuth.Delete(r)
	}
	cookie := &http.Cookie{
		Name:     "",
		Value:    "",
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now(),
	}
	http.SetCookie(w, cookie)
	requestCTX.SetAppResponse(true, http.StatusAccepted)
}

// swagger:route  GET /user/influencerid getUserIDByInfluencerID getUserIDByInfluencerID
// getUserIDByInfluencerID
//
// This endpoint will return user id by influencer id.
//
// Endpoint: /user/influencerid
//
// Method: GET
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetUserInfoByIDOpts
//     "$ref": "#/definitions/GetUserInfoByIDOpts"
//   required: true
//
// parameters:
// + name: id
//   in: query
//   schema:
//   type: ObjectID
//   required: true
//
// parameters:
// + name: cookie
//   in: header
//   description: Login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
//   in: header
//   description:Token required for successful response.
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description:Invalid User
//  200: ObjectID description: OK
func (a *API) getUserIDByInfluencerID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetUserInfoByIDOpts
	id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.ID = id
	res, err := a.App.User.GetUserIDByInfluencerID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) setUserGroups(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.SetUserGroupsOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	claim, sIDs, err := a.App.KeeperUser.SetUserGroups(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	token, err := a.TokenAuth.SignToken(*claim)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	for _, sessionID := range sIDs {
		if err := a.SessionAuth.UpdateSession(sessionID, token); err != nil {
			requestCTX.SetErr(errors.Wrap(err, "failed to update session id"), http.StatusBadRequest)
			return
		}
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getKeeperUsers(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetKeeperUsersOpts
	s.Query = r.URL.Query().Get("query")
	s.Page = uint(GetPageValue(r))
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.KeeperUser.GetKeeperUsers(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
