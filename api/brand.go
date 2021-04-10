package api

import (
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
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
	if s.UserID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
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
	userID, _ := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	res, err := a.App.Elasticsearch.GetBrandInfoByID(&schema.GetBrandsInfoByIDOpts{ID: id, UserID: userID})
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
