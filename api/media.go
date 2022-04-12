package api

import (
	"go-app/model"
	"go-app/schema"
	"go-app/server/handler"
	"net/http"
)

func (a *API) uploadImage(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateImageMediaOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Media.CreateImageMedia(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) uploadImageV2(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(10 * model.MB); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	file, header, err := r.FormFile("image")
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	s := schema.CreateImageMediaV2Opts{
		FileName: header.Filename,
		File:     file,
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Media.CreateImageMediaV2(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}
