package api

import (
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pasztorpisti/qs"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateInfluencerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.CreateInfluencer(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
	return
}

func (a *API) editInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.EditInfluencer(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencersByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByIDOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencersByID(s.IDs)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerByName(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByNameOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencerByName(s.Name)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /app/influencer/basic AppInfluencer getInfluencersBasic
// getInfluencersBasic
//
// This endpoint will return influencer basic information.
//
// Endpoint: /app/influencer/basic
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetInfluencersByIDBasicOpts
//     "$ref": "#/definitions/GetInfluencersByIDBasicOpts"
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
//  200: GetInfluencerBasicESEesp description: OK
func (a *API) getInfluencersBasic(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByIDBasicOpts
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
	res, err := a.App.Elasticsearch.GetInfluencersByIDBasic(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /app/influencer/{influencerID} AppInfluencer getInfluencerInfo
// getInfluencerInfo
//
// This endpoint will return influencer information.
//
// Endpoint: /app/influencer/{influencerID}
//
// Method: GET
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetInfluencerInfoByIDOpts
//     "$ref": "#/definitions/GetInfluencerInfoByIDOpts"
//   required: true
//
// parameters:
// + name: influencerID
//   in: path
//   schema:
//   type: string
//     "$ref": "#/definitions/ObjectID"
//   enum: 60b50277a97a2d73b211aec7
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
//  200: GetInfluencerInfoEsResp description: OK
func (a *API) getInfluencerInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["influencerID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid influencer id:%s in url", mux.Vars(r)["influencerID"]), http.StatusBadRequest)
		return
	}
	var userID primitive.ObjectID
	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetInfluencerInfoByID(&schema.GetInfluencerInfoByIDOpts{ID: id, CustomerID: userID})
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /app/user/influencer-request InfluencerRequest claimInfluencerRequest
// claimInfluencerRequest
//
// This endpoint will post the request claim by the influencer.
//
// Endpoint: /app/user/influencer-request
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: InfluencerAccountRequestOpts
//     "$ref": "#/definitions/InfluencerAccountRequestOpts"
//   required: true
//
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
func (a *API) claimInfluencerRequest(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.InfluencerAccountRequestOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	var source map[string]string
	if err := qs.Unmarshal(&source, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if len(source) != 0 {
		s.Source = source
	}
	if requestCTX.UserClaim != nil {
		s.UserID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if err := a.App.Influencer.InfluencerAccountRequest(&s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

// swagger:route  GET /app/user/influencer-request/status InfluencerRequest checkClaimInfluencerRequestStatus
// checkClaimInfluencerRequestStatus
//
// This endpoint will return the status of the influencer request.
//
// Endpoint: /app/user/influencer-request/status
//
// Method: GET
//
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
func (a *API) checkClaimInfluencerRequestStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var userID primitive.ObjectID
	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	}
	res, err := a.App.Influencer.GetInfluencerAccountRequestStatus(userID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerClaimRequests(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	res, err := a.App.Influencer.GetInfluencerAccountRequest()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) updateClaimInfluencerRequestStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateInfluencerAccountRequestStatusOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.GranteeID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	}
	if err := a.App.Influencer.UpdateInfluencerAccountRequestStatus(&s); err != nil {
		a.Logger.Err(err).Msg("failed to update status request")
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

// swagger:route  GET /influencer/check/username CheckUsername checkInfluencerUsernameExists
// checkInfluencerUsernameExists
//
// This endpoint will check the influencer username exists or not.
//
// Endpoint: /influencer/check/username
//
// Method: GET
//
// parameters:
// + name: username
//   in: query
//   schema:
//   enum: kartikay_sharma
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
//  200: description: true
func (a *API) checkInfluencerUsernameExists(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		requestCTX.SetErr(errors.Errorf("username cannot be empty nil"), http.StatusBadRequest)
		return
	}
	err := a.App.Influencer.CheckInfluencerUsernameExists(username, nil)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

// swagger:route  POST /app/influencer/username/basic AppInfluencer getInfluencersBasicByUsername
// getInfluencersBasicByUsername
//
// This endpoint will return influencer basic information by username.
//
// Endpoint: /app/influencer/username/basic
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetInfluencersByUsernameBasicOpts
//     "$ref": "#/definitions/GetInfluencersByUsernameBasicOpts"
//   required: true
//
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
//  200: GetInfluencerBasicESEesp description: OK
func (a *API) getInfluencersBasicByUsername(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByUsernameBasicOpts
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
	res, err := a.App.Elasticsearch.GetInfluencersByUserameBasic(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /app/influencer/username/{username} AppInfluencer getInfluencerInfoByUsername
// getInfluencerInfoByUsername
//
// This endpoint will return influencer information by username.
//
// Endpoint: /app/influencer/username/{username}
//
// Method: GET
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetInfluencerInfoByUsernameOpts
//     "$ref": "#/definitions/GetInfluencerInfoByUsernameOpts"
//   required: true
//
//
// parameters:
// + name: username
//   in: path
//   schema:
//   type: string
//   enum: kartikay_sharma
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
//  200: GetInfluencerInfoEsResp description: OK
func (a *API) getInfluencerInfoByUsername(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	if username == "" {
		requestCTX.SetErr(errors.Errorf("invalid influencer id:%s in url", mux.Vars(r)["influencerID"]), http.StatusBadRequest)
		return
	}
	var userID primitive.ObjectID
	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetInfluencerInfoByUsername(&schema.GetInfluencerInfoByUsernameOpts{Username: username, CustomerID: userID})
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  PUT /app/influencer AppInfluencer editInfluencerApp
// editInfluencerApp
//
// This endpoint will edit the influencer details.
//
// Endpoint: /app/influencer
//
// Method: PUT
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: EditInfluencerAppOpts
//     "$ref": "#/definitions/EditInfluencerAppOpts"
//   required: true
//
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
//  403: AppErr description:Not Authorized
//  200: EditInfluencerResp description: OK
func (a *API) editInfluencerApp(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerAppOpts

	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID {
		requestCTX.SetErr(errors.New("not authorized"), http.StatusForbidden)
		return
	}

	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.EditInfluencerApp(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /app/creator/debit-request DebitRequest debitRequest
// debitRequest
//
// This endpoint will post the debit request.
//
// Endpoint: /app/creator/debit-request
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CommissionDebitRequest
//     "$ref": "#/definitions/CommissionDebitRequest"
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
func (a *API) debitRequest(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CommissionDebitRequest
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.ID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if err := a.App.Influencer.DebitRequest(&s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) updateDebitRequest(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCommissionDebitRequest
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.GranteeID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	}

	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if err := a.App.Influencer.UpdateDebitRequest(&s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getDebitRequest(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	res, err := a.App.Influencer.GetActiveDebitRequest()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /app/creator/dashboard Creator getInfluencerDashboard
// getInfluencerDashboard
//
// This endpoint will return dashaboard of the influencer.
//
// Endpoint: /app/creator/dashboard
//
// Method: POST
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetInfluencerDashboardOpts
//     "$ref": "#/definitions/GetInfluencerDashboardOpts"
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
//  200: GetInfluencerDashboardResp description: OK
func (a *API) getInfluencerDashboard(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencerDashboardOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencerDashboard(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /app/creator/ledger Creator getInfluencerLedger
// getInfluencerLedger
//
// This endpoint will return the ledger details of the influencer.
//
// Endpoint: /app/creator/ledger
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetInfluencerLedgerOpts
//     "$ref": "#/definitions/GetInfluencerLedgerOpts"
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
//  200: GetInfluencerLedgerResp description: OK
func (a *API) getInfluencerLedger(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencerLedgerOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.ID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencerLedger(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /app/creator/payout-info Creator getInfluencerPayoutInfo
// getInfluencerPayoutInfo
//
// This endpoint will return influencer payout information.
//
// Endpoint: /app/creator/payout-info
//
// Method: GET
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
//  200: GetPayoutInfoResp description: OK
func (a *API) getInfluencerPayoutInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencerPayoutInfo(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /app/creator/commission Creator getCommissionAndRevenue
// getCommissionAndRevenue
//
// This endpoint will return commision and revenue.
//
// Endpoint: /app/creator/commission
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetCommissionAndRevenueOpts
//     "$ref": "#/definitions/GetCommissionAndRevenueOpts"
//   required: true
//
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
//  200: GetCommissionAndRevenueResp description: OK
func (a *API) getCommissionAndRevenue(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetCommissionAndRevenueOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.ID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetCommissionAndRevenue(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  PUT /v2/app/influencer v2Influnencer editInfluencerAppV2
// editInfluencerAppV2
//
// This endpoint edit the influencer information.
//
// Endpoint: /v2/app/influencer
//
// Method: PUT
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: EditInfluencerAppV2Opts
//     "$ref": "#/definitions/EditInfluencerAppV2Opts"
//   required: true
//
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
//  200: EditInfluencerResp description: OK
func (a *API) editInfluencerAppV2(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerAppV2Opts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	var err error
	s.ID, err = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(errors.Wrapf(err, "influencer id invalid"), http.StatusBadRequest)
		return
	}
	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID {
		requestCTX.SetErr(errors.New("not authorized"), http.StatusForbidden)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.EditInfluencerAppV2(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /v2/app/user/influencer-request v2Influnencer claimInfluencerRequestV2
// claimInfluencerRequestV2
//
// This endpoint will claim the influencer request.
//
// Endpoint: /v2/app/user/influencer-request
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: InfluencerAccountRequestV2Opts
//     "$ref": "#/definitions/InfluencerAccountRequestV2Opts"
//   required: true
//
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
func (a *API) claimInfluencerRequestV2(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.InfluencerAccountRequestV2Opts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
		s.CustomerID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	var source map[string]string
	if err := qs.Unmarshal(&source, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if len(source) != 0 {
		s.Source = source
	}
	if s.Email != "" {
		err := a.App.User.UpdateUserEmail(&schema.UpdateUserEmailOpts{
			ID:    s.UserID,
			Email: s.Email,
		})
		if err != nil {
			requestCTX.SetErr(err, http.StatusBadRequest)
			return
		}
	}
	if s.Phone != nil {
		err := a.App.User.UpdateUserPhoneNo(&schema.UpdateUserPhoneNoOpts{
			ID:      s.UserID,
			PhoneNo: s.Phone,
		})
		if err != nil {
			requestCTX.SetErr(err, http.StatusBadRequest)
			return
		}
	}
	if err := a.App.Influencer.InfluencerAccountRequestV2(&s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
