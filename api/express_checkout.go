package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/pkg/errors"
	"github.com/vasupal1996/goerror"
)

// swagger:route  POST /app/express-checkout ExpressCheckout expressCheckout
// expressCheckout
//
// This endpoint is for express checkout.
//
// Endpoint: /app/express-checkout
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: ExpressCheckoutOpts
//     "$ref": "#/definitions/ExpressCheckoutOpts"
//   required: true
//
// parameters:
// + name: platform
//   in: query
//   description: Platform type for example android, web or ios.
//   schema:
//   type: string
//   required: true
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
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description: Invalid User
//  200: OrderInfo description: OK
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

// swagger:route  POST /web/express-checkout ExpressCheckout expressCheckoutWeb
// expressCheckoutWeb
//
// This endpoint is for express checkout for web.
//
// Endpoint: /web/express-checkout
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: ExpressCheckoutWebOpts
//     "$ref": "#/definitions/ExpressCheckoutWebOpts"
//   required: true
//
// parameters:
// + name: platform
//   in: query
//   description: Platform type for example android, web or ios.
//   schema:
//   type: string
//   required: true
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
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description: Invalid User
//  200: OrderInfo description: OK
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
	fmt.Println("isCod:", s.IsCOD)
	if s.IsCOD {
		if s.RequestID == "" {
			requestCTX.SetErr(goerror.New("request id is required for cod orders", &goerror.BadRequest), http.StatusBadRequest)
			return
		}
	}
	resp, err := a.App.ExpressCheckout.ExpressCheckoutWeb(&s, fullName)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  POST /app/express-checkout/check/cod CODviaGoKwik expressCheckoutRTO
// expressCheckoutRTO
//
// This endpoint will express checkout RTO.
//
// Endpoint: /app/express-checkout/check/cod
//
// Method: POST
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: ExpressCheckoutWebOpts
//     "$ref": "#/definitions/ExpressCheckoutWebOpts"
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
//  200: description: OK
func (a *API) expressCheckoutRTO(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
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
	email := requestCTX.UserClaim.(*auth.UserClaim).Email

	resp, err := a.App.ExpressCheckout.ExpressCheckoutRTO(&s, fullName, r.UserAgent(), r.RemoteAddr, email)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}
