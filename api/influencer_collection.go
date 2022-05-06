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

func (a *API) createInfluencerCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateInfluencerCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	// s.InfluencerID,err:=
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.InfluencerCollection.CreateInfluencerCollection(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

func (a *API) keeperGetInfluencerCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencerCollectionsOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	fmt.Println(s)
	res, err := a.App.InfluencerCollection.KeeperGetInfluencerCollections(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

func (a *API) editInfluencerCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.InfluencerCollection.EditInfluencerCollection(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

// swagger:route GET /app/influencer/collections/active InfluencerCollectionApp getActiveInfluencerCollections
// getActiveInfluencerCollections
//
// This endpoint will return active influencer collections.
//
// Endpoint: /app/influencer/collections/active
//
// Method: GET
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetActiveInfluencerCollectionsOpts
//     "$ref": "#/definitions/GetActiveInfluencerCollectionsOpts"
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
//  200: GetInfluencerCollectionESResp description: OK
func (a *API) getActiveInfluencerCollections(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetActiveInfluencerCollectionsOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	resp, err := a.App.Elasticsearch.GetActiveInfluencerCollections(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route GET /app/influencer/collections InfluencerCollectionApp appGetInfluencerCollections
// appGetInfluencerCollections
//
// This endpoint will return the influencer collections.
//
// Endpoint: /app/influencer/collections
//
// Method: GET
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetInfluencerCollectionsOpts
//     "$ref": "#/definitions/GetInfluencerCollectionsOpts"
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
//  200: GetInfluencerCollectionRespApp description: OK
func (a *API) appGetInfluencerCollections(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencerCollectionsOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	iid, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.InfluencerID = iid.Hex()
	resp, err := a.App.InfluencerCollection.GetInfluencerCollectionsByInfluencerIDApp(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route POST /app/influencer/collection InfluencerCollectionApp createInfluencerCollectionApp
// createInfluencerCollectionApp
//
// This endpoint will create influencer collection.
//
// Endpoint: /app/influencer/collection
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreateInfluencerCollectionOpts
//     "$ref": "#/definitions/CreateInfluencerCollectionOpts"
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
//  200: InfluencerCollectionResp description: OK
func (a *API) createInfluencerCollectionApp(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateInfluencerCollectionOpts
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
	res, err := a.App.InfluencerCollection.CreateInfluencerCollection(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

// swagger:route PUT /app/influencer/collection InfluencerCollectionApp editInfluencerCollectionApp
// editInfluencerCollectionApp
//
// This endpoint will edit the influencer collection app.
//
// Endpoint: /app/influencer/collection
//
// Method: PUT
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: EditInfluencerCollectionAppOpts
//     "$ref": "#/definitions/EditInfluencerCollectionAppOpts"
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
//  200: InfluencerCollectionResp description: OK
func (a *API) editInfluencerCollectionApp(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerCollectionAppOpts
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
	res, err := a.App.InfluencerCollection.EditInfluencerCollectionApp(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

// swagger:route GET /app/influencer/collection InfluencerCollectionApp getActiveInfluencerCollectionByID
// getActiveInfluencerCollectionByID
//
// This endpoint will return the active influencer collection by ID.
//
// Endpoint: /app/influencer/collection
//
// Method: GET
//
// parameters:
// + name: id
//   in: query
//   schema:
//   type: string
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
//  200: GetInfluencerCollectionESResp description: OK
func (a *API) getActiveInfluencerCollectionByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	resp, err := a.App.Elasticsearch.GetActiveInfluencerCollectionByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}
