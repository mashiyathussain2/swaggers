package api

import (
	"go-app/schema"
	"go-app/server/handler"
	"net/http"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createSizeProfile(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	var s schema.CreateSizeProfileOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.SizeProfile.CreateSizeProfile(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) addBrandToSizeProfile(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	var s schema.AddBrandToSizeProfileOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.SizeProfile.AddBrandToSizeProfile(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) getSizeProfilesForBrand(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("brandID"))
	if err != nil {
		requestCTX.SetErr(errors.Wrapf(err, "brandID provided is not in correct format"), http.StatusBadRequest)
	}

	resp, err := a.App.SizeProfile.GetSizeProfilesForBrand(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getSizeProfile(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
	if err != nil {
		requestCTX.SetErr(errors.Wrapf(err, "brandID provided is not in correct format"), http.StatusBadRequest)
	}

	resp, err := a.App.SizeProfile.GetSizeProfile(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getAllSizeProfiles(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	resp, err := a.App.SizeProfile.GetAllSizeProfiles()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}
