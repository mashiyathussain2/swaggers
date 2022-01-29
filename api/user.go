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

func (a *API) me(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	requestCTX.SetAppResponse(requestCTX.UserClaim, http.StatusOK)
}

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
	requestCTX.SetRedirectResponse(redirectURL, http.StatusPermanentRedirect)
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
