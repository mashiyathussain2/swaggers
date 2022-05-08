package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/pasztorpisti/qs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:route POST /app/influencer/products InfluencerCollectionKEEPER addInfluencerProducts
// addInfluencerProducts
//
// This endpoint will add influencer products.
//
// Endpoint: /app/influencer/products
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: AddInfluencerProductsOpts
//     "$ref": "#/definitions/AddInfluencerProductsOpts"
//   required: true
//
// parameters:
// + name: cookie
//   type: string
//   in: header
//   description: Customer login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
//   type: string
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
//  200: description: true
func (a *API) addInfluencerProducts(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddInfluencerProductsOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	var err error
	s.InfluencerID, err = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err = a.App.InfluencerProducts.AddInfluencerProductsOpts(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
}

// swagger:route DELETE /app/influencer/products InfluencerCollectionKEEPER removeInfluencerProducts
// removeInfluencerProducts
//
// This endpoint will remove influencer products.
//
// Endpoint: /app/influencer/products
//
// Method: DELETE
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: RemoveInfluencerProductsOpts
//     "$ref": "#/definitions/RemoveInfluencerProductsOpts"
//   required: true
//
// parameters:
// + name: cookie
//   type: string
//   in: header
//   description: Customer login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
//   type: string
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
//  200: description: true
func (a *API) removeInfluencerProducts(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.RemoveInfluencerProductsOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	var err error
	s.InfluencerID, err = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err = a.App.InfluencerProducts.RemoveInfluencerProductsOpts(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
}

// swagger:route GET /app/influencer/products InfluencerCollectionKEEPER getInfluencerProducts
// getInfluencerProducts
//
// This endpoint will return influencer products.
//
// Endpoint: /app/influencer/products
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetInfluencerProducts
//     "$ref": "#/definitions/GetInfluencerProducts"
//   required: true
//
// parameters:
// + name: cookie
//   type: string
//   in: header
//   description:Login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
//   type: string
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
//  200: GetInfluencerProductESResp description: OK
func (a *API) getInfluencerProducts(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencerProducts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	var resp *schema.GetInfluencerProductESResp
	var err error
	if s.Type == "self" {
		s.InfluencerID = requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID
		resp, err = a.App.InfluencerProducts.GetInfluencerProductsOpts(s.InfluencerID)
		fmt.Println("resp from db", resp)
	} else {
		if errs := a.Validator.Validate(&s); errs != nil {
			requestCTX.SetErrs(errs, http.StatusBadRequest)
			return
		}
		resp, err = a.App.Elasticsearch.GetInfluencerProducts(s.InfluencerID, s.Page)
	}
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusCreated)
}
