package api

import (
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/pkg/errors"
	"github.com/vasupal1996/goerror"
)

func (a *API) expressCheckout(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ExpressCheckoutOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.UserID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	fullName := requestCTX.UserClaim.(*auth.UserClaim).FullName

	platform := r.URL.Query().Get("platform")

	// To remove after app update
	if platform == "" {
		platform = "android"
	}
	if platform != "web" && platform != "android" && platform != "ios" {
		requestCTX.SetErr(goerror.New("platform incorrect", &goerror.BadRequest), http.StatusBadRequest)
	}

	resp, err := a.App.ExpressCheckout.ExpressCheckoutComplete(&s, fullName, platform)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) expressCheckoutWeb(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ExpressCheckoutWebOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.UserID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	fullName := requestCTX.UserClaim.(*auth.UserClaim).FullName

	resp, err := a.App.ExpressCheckout.ExpressCheckoutWeb(&s, fullName)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}
