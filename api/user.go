package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func (a *API) me(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	requestCTX.SetAppResponse(requestCTX.UserClaim, http.StatusOK)
}

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

func (a *API) verifyEmail(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.VerifyEmailOpts
	var returnToken bool
	resp := make(map[string]interface{})

	var isWeb bool
	isWeb, _ = strconv.ParseBool(r.URL.Query().Get("isWeb"))
	returnToken, _ = strconv.ParseBool(r.URL.Query().Get("returnToken"))

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
			if returnToken {
				resp["token"] = token
			}
		}
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

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
	url := a.App.KeeperUser.Login()
	requestCTX.SetRedirectResponse(url, http.StatusTemporaryRedirect)
}

func (a *API) keeperLoginCallback(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	claim, err := a.App.KeeperUser.Callback(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	token, err := a.TokenAuth.SignKeeperToken(claim)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	redirectURL := fmt.Sprintf("%s?token=%s", a.Config.KeeperLoginRedirectURL, token)
	requestCTX.SetRedirectResponse(redirectURL, http.StatusPermanentRedirect)
}

func (a *API) updateUserAuthInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateUserAuthOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	err := a.App.User.UpdateUserAuthInfo(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusAccepted)
}

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

func (a *API) verifyUserAuthUpdate(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.VerifyUserAuthUpdate
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
	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	claim, err := a.App.User.VerifyUserAuthUpdate(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
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
		requestCTX.SetAppResponse(true, http.StatusOK)
		return
	}
	requestCTX.SetAppResponse(token, http.StatusAccepted)
}
