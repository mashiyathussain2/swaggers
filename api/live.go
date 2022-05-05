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

func (a *API) getLiveStreams(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetLiveStreamsFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.GetLiveStreams(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) getLiveStreamByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["liveID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid live id: %s in url", mux.Vars(r)["liveID"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.GetLiveStreamByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) createLiveStream(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateLiveStreamOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.CreateLiveStream(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route POST /live/{liveID}/comment live pushComment
// pushComment
//
// This endpoint will post the comment on live stream.
//
// Endpoint: /live/{liveID}/comment
//
// Method: POST
//
// parameters:
// + name: liveID
//   in: path
//   schema:
//   type: string
//     "$ref": "#/definitions/ObjectID"
//   enum: 6065d4503824bf77961c21ae
//   required: true
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreateLiveCommentOpts
//     "$ref": "#/definitions/CreateLiveCommentOpts"
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
//  200: description: OK
func (a *API) pushComment(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateLiveCommentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.UserID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	a.App.Live.PushComment(&s)
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

// swagger:route GET /app/live/{liveID}/start LiveStream startLiveStream
// startLiveStream
//
// This endpoint will start live stream.
//
// Endpoint: /app/live/{liveID}/start
//
// Method: GET
//
// parameters:
// + name: liveID
//   in: path
//   schema:
//   type: string
//     "$ref": "#/definitions/ObjectID"
//   enum: 6065d4503824bf77961c21ae
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
//  200: StartLiveStreamResp description: OK
func (a *API) startLiveStream(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["liveID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid live id: %s in url", mux.Vars(r)["content"]), http.StatusBadRequest)
		return
	}

	res, err := a.App.Live.StartLiveStream(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route  GET /live/{liveID}/join live joinLiveStream
// joinLiveStream
//
// This endpoint will join live stream.
//
// Endpoint: /live/{liveID}/join
//
// Method: GET
//
// parameters:
// + name: liveID
//   in: path
//   schema:
//   type: string
//     "$ref": "#/definitions/ObjectID"
//   enum: 6065d4503824bf77961c21ae
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
//  200: JoinLiveStreamResp description: OK
func (a *API) joinLiveStream(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["liveID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid live id: %s in url", mux.Vars(r)["content"]), http.StatusBadRequest)
		return
	}

	res, err := a.App.Live.JoinLiveStream(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route  POST /live/{liveID}/joined live joinedLiveStream
// joinedLiveStream
//
// This endpoint will joined the stream.
//
// Endpoint: /live/{liveID}/joined
//
// Method: POST
//
// parameters:
// + name: liveID
//   in: path
//   schema:
//   type: string
//     "$ref": "#/definitions/ObjectID"
//   enum: 6065d4503824bf77961c21ae
//   required: true
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: PushJoinOpts
//     "$ref": "#/definitions/PushJoinOpts"
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
//  200: description: OK
func (a *API) joinedLiveStream(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.PushJoinOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	a.App.Live.PushJoin(&s)
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

// swagger:route GET /app/live/{liveID}/stop LiveStream stopLiveStream
// stopLiveStream
//
// This endpoint will stop the live stream.
//
// Endpoint: /app/live/{liveID}/stop
//
// Method: GET
//
//
// parameters:
// + name: liveID
//   in: path
//   schema:
//   type: string
//     "$ref": "#/definitions/ObjectID"
//   enum: 6065d4503824bf77961c21ae
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
//  200: description: OK
func (a *API) stopLiveStream(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["liveID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid live id: %s in url", mux.Vars(r)["content"]), http.StatusBadRequest)
		return
	}

	if err := a.App.Live.EndLiveStream(id); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

// swagger:route POST /app/live/{liveID}/catalog LiveStream pushCatalog
// pushCatalog
//
// This endpoint push catalog.
//
// Endpoint: /app/live/{liveID}/catalog
//
// Method: POST
//
// parameters:
// + name: liveID
//   in: path
//   schema:
//   type: string
//     "$ref": "#/definitions/ObjectID"
//   enum: 6065d4503824bf77961c21ae
//   required: true
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: PushCatalogOpts
//     "$ref": "#/definitions/PushCatalogOpts"
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
//  200: description: true
func (a *API) pushCatalog(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.PushCatalogOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	a.App.Live.PushCatalog(&s)
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

func (a *API) pushOrder(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.PushNewOrderOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	a.App.Live.PushOrder(&s)
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

// swagger:route  GET /live live getAppLiveStreams
// getAppLiveStreams
//
// This endpoint get the app live streams.
//
// Endpoint: /live
//
// Method: GET
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetAppLiveStreamsFilter
//     "$ref": "#/definitions/GetAppLiveStreamsFilter"
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
//  200: GetAppLiveStreamResp description: OK
func (a *API) getAppLiveStreams(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetAppLiveStreamsFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.GetAppLiveStreams(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route  GET /live/{liveID} live getAppLiveStreamByID
// getAppLiveStreamByID
//
// This endpoint get the app live streams by ID.
//
// Endpoint: /live/{liveID}
//
// Method: GET
//
// parameters:
// + name: liveID
//   in: path
//   schema:
//   type: string
//     "$ref": "#/definitions/ObjectID"
//   enum: 6065d4503824bf77961c21ae
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
//  200: GetAppLiveStreamResp description: OK
func (a *API) getAppLiveStreamByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["liveID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid live id: %s in url", mux.Vars(r)["liveID"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.GetAppLiveStreamByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route GET /app/influencer/live LiveStream getAppLiveStreamsByInfluencerID
// getAppLiveStreamsByInfluencerID
//
// This endpoint will return app live streams by influencerID.
//
// Endpoint: /app/influencer/live
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetAppLiveStreamsFilter
//     "$ref": "#/definitions/GetAppLiveStreamsFilter"
//   required: true
//
// parameters:
// + name: cookie
//   type: string
//   in: header
//   description: Login required for successful response.
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
//  200: GetAppLiveStreamInfluencerResp description: OK
func (a *API) getAppLiveStreamsByInfluencerID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetAppLiveStreamsFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	influencer_id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if influencer_id == primitive.NilObjectID {
		requestCTX.SetErr(errors.New("influencer id missing"), http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.GetAppLiveStreamsByInfluencerID(influencer_id, &s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route POST /app/influencer/live LiveStream createLiveStreamByApp
// createLiveStreamByApp
//
// This endpoint will create live stream.
//
// Endpoint: /app/influencer/live
//
// Method: POST
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreateLiveStreamOpts
//     "$ref": "#/definitions/CreateLiveStreamOpts"
//   required: true
//
// parameters:
// + name: cookie
//   type: string
//   in: header
//   description: Login required for successful response.
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
//  200: CreateLiveStreamResp description: OK
func (a *API) createLiveStreamByApp(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateLiveStreamOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	influencer_id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if influencer_id == primitive.NilObjectID {
		requestCTX.SetErr(errors.New("influencer id missing from profile"), http.StatusBadRequest)
		return
	}
	s.InfluencerIDs = append(s.InfluencerIDs, influencer_id)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}

	res, err := a.App.Live.CreateLiveStream(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route GET /v2/app/influencer/live LiveStream v2GetAppLiveStreamsByInfluencerID
// v2GetAppLiveStreamsByInfluencerID
//
// This endpoint return app live stream by influencer ID.
//
// Endpoint: /v2/app/influencer/live
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetAppLiveStreamsFilter
//     "$ref": "#/definitions/GetAppLiveStreamsFilter"
//   required: true
//
// parameters:
// + name: cookie
//   type: string
//   in: header
//   description: Login required for successful response.
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
//  200: GetLiveByInfluencerID description: OK
func (a *API) v2GetAppLiveStreamsByInfluencerID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetAppLiveStreamsFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	influencer_id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if influencer_id == primitive.NilObjectID {
		requestCTX.SetErr(errors.New("influencer id missing"), http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.GetAppLiveStreamsByInfluencerIDV2(influencer_id, &s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}
