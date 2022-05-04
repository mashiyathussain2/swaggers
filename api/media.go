package api

import (
	"go-app/model"
	"go-app/schema"
	"go-app/server/handler"
	"net/http"
)

// swagger:route  POST /image/upload UploadImage uploadImage
// uploadImage
//
// This endpoint post image.
//
// Endpoint: /image/upload
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreateImageMediaOpts
//     "$ref": "#/definitions/CreateImageMediaOpts"
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
//  403: AppErr description:Invalid User
//  200: CreateImageMediaResp description: OK
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

// swagger:route  POST /v2/image/upload UploadImage uploadImageV2
// uploadImageV2
//
// This endpoint post image.
//
// Endpoint: /v2/image/upload
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreateImageMediaV2Opts
//     "$ref": "#/definitions/CreateImageMediaV2Opts"
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
//  403: AppErr description:Invalid User
//  200: CreateImageMediaResp description: OK
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
