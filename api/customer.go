package api

import (
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
	requestCTX.SetAppResponse(token, http.StatusOK)
}

func (a *API) signUpViaEmail(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateUserOpts
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
	requestCTX.SetAppResponse(token, http.StatusOK)
}

func (a *API) updateCustomerInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCustomerOpts
	var returnToken bool
	resp := make(map[string]interface{})

	returnToken, _ = strconv.ParseBool(r.URL.Query().Get("returnToken"))

	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Customer.UpdateCustomer(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	if returnToken {
		claim := requestCTX.UserClaim.(*auth.UserClaim)
		claim.FullName = res.FullName
		if res.Gender != nil {
			claim.Gender = *res.Gender
		}
		claim.ProfileImage = res.ProfileImage
		claim.DOB = res.DOB
		token, err := a.TokenAuth.SignToken(claim)
		if err != nil {
			requestCTX.SetErr(err, http.StatusBadRequest)
			return
		}
		resp["token"] = token
	}

	resp["user"] = res
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
	res, err := a.App.Influencer.AddFollower(&s)
	fmt.Println(res, err)
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
	res, err := a.App.Brand.AddFollower(&s)
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
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Customer.AddAddress(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	resp, err := a.App.Customer.GetAddresses(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}
