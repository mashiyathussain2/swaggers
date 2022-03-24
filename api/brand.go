package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pasztorpisti/qs"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createbrand(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateBrandOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.CreateBrand(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) createBrandAdminUser(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateBrandAdminUserOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.CreateBrandAdminUser(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getBrandsById(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetBrandsByIDOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.GetBrandsByID(s.IDs)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) editbrand(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditBrandOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.EditBrand(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getBrandByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["brandID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid brand id:%s in url", mux.Vars(r)["brandID"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.GetBrandByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) checkBrandByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["brandID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid brand id:%s in url", mux.Vars(r)["brandID"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.CheckBrandByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getBrands(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	res, err := a.App.Brand.GetBrands()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getBrandsBasic(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetBrandsByIDBasicOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.CustomerID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetBrandsByIDBasic(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getBrandInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["brandID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid brand id:%s in url", mux.Vars(r)["brandID"]), http.StatusBadRequest)
		return
	}
	var userID primitive.ObjectID
	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetBrandInfoByID(&schema.GetBrandsInfoByIDOpts{ID: id, CustomerID: userID})
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getActiveBrandsList(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetBrandsListOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.GetBrandsList(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) brandUserLogin(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.BrandUserLoginOpts
	var returnToken bool
	returnToken, _ = strconv.ParseBool(r.URL.Query().Get("returnToken"))

	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}

	res, err := a.App.Brand.BrandUserLogin(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	token, err := a.TokenAuth.SignToken(res)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if !returnToken {
		if err := a.SessionAuth.Create(token, w); err != nil {
			requestCTX.SetErr(fmt.Errorf("failed to login user: %s", err), http.StatusInternalServerError)
			return
		}
		requestCTX.SetAppResponse(true, http.StatusOK)
		return
	}
	requestCTX.SetAppResponse(token, http.StatusOK)
}

func (a *API) brandUserForgotPassword(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ForgotPasswordOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.ForgotPassword(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) brandUserResetPassword(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ResetPasswordOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.ResetPassword(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) checkBrandUsernameExists(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		requestCTX.SetErr(errors.Errorf("username cannot be empty nil"), http.StatusBadRequest)
		return
	}
	err := a.App.Brand.CheckBrandUsernameExists(username, nil)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getBrandsBasicByUsername(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetBrandsByUsernameBasicOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.CustomerID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetBrandsByUsernameBasic(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getBrandInfoByUsername(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	if username == "" {
		requestCTX.SetErr(errors.Errorf("invalid brand id:%s in url", mux.Vars(r)["username"]), http.StatusBadRequest)
		return
	}
	var userID primitive.ObjectID
	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetBrandInfoByUsername(&schema.GetBrandsInfoByUsernameOpts{Username: username, CustomerID: userID})
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
