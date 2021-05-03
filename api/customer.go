package api

import (
	"errors"
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) loginViaEmail(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EmailLoginCustomerOpts
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
	res, err := a.App.Customer.Login(&s)
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

func (a *API) signUpViaEmail(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateUserOpts
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
	res, err := a.App.Customer.SignUp(&s)
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

func (a *API) updateCustomerInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCustomerOpts
	var isWeb bool
	isWeb, _ = strconv.ParseBool(r.URL.Query().Get("isWeb"))
	resp := make(map[string]interface{})
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	s.ID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	claim, err := a.App.Customer.UpdateCustomer(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	resp["user"] = claim
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
	resp["token"] = token
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) followInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddInfluencerFollowerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.CustomerID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Influencer.AddFollower(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) followBrand(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddBrandFollowerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.CustomerID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Brand.AddFollower(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) unFollowInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddInfluencerFollowerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.CustomerID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Influencer.RemoveFollower(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) unFollowBrand(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddBrandFollowerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.CustomerID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Brand.RemoveFollower(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) addAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddAddressOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.UserID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.Customer.AddAddress(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if id.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	resp, err := a.App.Customer.GetAddresses(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCustomerInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["customerID"])
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if id.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Customer.GetAppCustomerInfo(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) removeAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	userID, err := primitive.ObjectIDFromHex(r.FormValue("user_id"))
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if userID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	addressID, err := primitive.ObjectIDFromHex(r.FormValue("address_id"))
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	err = a.App.Customer.RemoveAddress(userID, addressID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) editAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditAddressOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	//TODO:WHY this not just check?
	s.UserID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Customer.EditAddress(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
