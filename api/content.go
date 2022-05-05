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

func (a *API) getContent(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetContentFilter
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.GetContent(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) getContentByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["pebbleID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid content id: %s in url", mux.Vars(r)["content"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.GetContentByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) createPebble(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreatePebbleOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.CreatePebble(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) editPebble(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditPebbleOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.EditPebble(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) deleteContent(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["contentID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid pebble id: %s in url", mux.Vars(r)["contentID"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.DeleteContent(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
	return
}

func (a *API) processPebble(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ProcessVideoContentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		a.Logger.Err(err).Msg("invalid json body")
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		a.Logger.Err(errors.New("invalid json request")).Errs("errs", errs).Msg("failed to validate ProcessVideoContentOpts")
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.ProcessVideoContent(&s)
	if err != nil {
		a.Logger.Err(err).Interface("opts", s).Msg("failed to ProcessVideoContent")
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
	return
}

func (a *API) createVideoCatalogContent(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateVideoCatalogContentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.CreateCatalogVideoContent(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) createVideoReviewContent(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateVideoReviewContentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.CreateVideoReviewContent(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) createImageCatalogContent(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateImageCatalogContentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.CreateCatalogImageContent(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) editCatalogContent(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditCatalogContentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.EditCatalogContent(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route POST /content/comment content createContentComment
// createContentComment
//
// This endpoint will post the comment on the content.
//
// Endpoint: /content/comment
//
// Method: POST
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreateCommentOpts
//     "$ref": "#/definitions/CreateCommentOpts"
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
//  200: CreateCommentResp description: OK
func (a *API) createContentComment(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateCommentOpts
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
	res, err := a.App.Content.CreateComment(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route POST /content/like content createLike
// createLike
//
// This endpoint will create like on content.
//
// Endpoint: /content/like
//
// Method: POST
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreateLikeOpts
//     "$ref": "#/definitions/CreateLikeOpts"
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
func (a *API) createLike(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateLikeOpts
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
	err := a.App.Content.CreateLike(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

// swagger:route POST /content/view content createView
// createView
//
// This endpoint will create view on content.
//
// Endpoint: /content/view
//
// Method: POST
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreateViewOpts
//     "$ref": "#/definitions/CreateViewOpts"
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
func (a *API) createView(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateViewOpts
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
	err := a.App.Content.CreateView(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

func (a *API) getPebble(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebbleFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID = requestCTX.UserClaim.(*auth.UserClaim).ID
	}
	res, err := a.App.Elasticsearch.GetPebble(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route GET /pebble Pebble getPebbleV2
// getPebbleV2
//
// This endpoint return the pebbles.
//
// Endpoint: /pebble
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetPebbleFilter
//     "$ref": "#/definitions/GetPebbleFilter"
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
//  200: GetPebbleESResp description: OK
func (a *API) getPebbleV2(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebbleFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID = requestCTX.UserClaim.(*auth.UserClaim).ID
	}
	res, err := a.App.Elasticsearch.GetPebbleV2(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route GET /pebble/id Pebble getPebbleByID
// getPebbleByID
//
// This endpoint return the pebbles by id.
//
// Endpoint: /pebble/id
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetPebbleByIDFilter
//     "$ref": "#/definitions/GetPebbleByIDFilter"
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
//  200: GetPebbleESResp description: OK
func (a *API) getPebbleByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebbleByIDFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID = requestCTX.UserClaim.(*auth.UserClaim).ID
	}
	res, err := a.App.Elasticsearch.GetPebbleByID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) getContents(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebblesKeeperFilter
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.GetPebbles(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) changeContentStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ChangeContentStatusOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.ChangeContentStatus(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route GET /pebble/brand Pebble getPebblesByBrandID
// getPebblesByBrandID
//
// This endpoint return the pebbles of brand by ID.
//
// Endpoint: /pebble/brand
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetPebbleByBrandID
//     "$ref": "#/definitions/GetPebbleByBrandID"
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
//  200: GetPebbleESResp description: OK
func (a *API) getPebblesByBrandID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebbleByBrandID
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID = requestCTX.UserClaim.(*auth.UserClaim).ID
	}
	res, err := a.App.Elasticsearch.GetPebblesByBrandID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route GET /pebble/influencer Pebble getPebblesByInfluencerID
// getPebblesByInfluencerID
//
// This endpoint return the pebbles by influencer ID.
//
// Endpoint: /pebble/influencer
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetPebbleByInfluencerID
//     "$ref": "#/definitions/GetPebbleByInfluencerID"
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
//  200: GetPebbleESResp description: OK
func (a *API) getPebblesByInfluencerID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebbleByInfluencerID
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID = requestCTX.UserClaim.(*auth.UserClaim).ID
	}
	res, err := a.App.Elasticsearch.GetPebblesByInfluencerID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route GET /catalog/influencer Pebble getCatalogsByInfluencerID
// getCatalogsByInfluencerID
//
// This endpoint return the catalogs by influencer ID.
//
// Endpoint: /catalog/influencer
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetCatalogsByInfluencerID
//     "$ref": "#/definitions/GetCatalogsByInfluencerID"
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
//  200: description: ObjectID
func (a *API) getCatalogsByInfluencerID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetCatalogsByInfluencerID
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID = requestCTX.UserClaim.(*auth.UserClaim).ID
	}
	res, err := a.App.Elasticsearch.GetCatalogsByInfluencerID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route GET /pebble/hashtag Pebble getPebblesByHashtag
// getPebblesByHashtag
//
// This endpoint return the pebbles by respective hashtags.
//
// Endpoint: /pebble/hashtag
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetPebbleByHashtag
//     "$ref": "#/definitions/GetPebbleByHashtag"
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
//  200: GetPebbleESResp description: OK
func (a *API) getPebblesByHashtag(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebbleByHashtag
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID = requestCTX.UserClaim.(*auth.UserClaim).ID
	}
	res, err := a.App.Elasticsearch.GetPebblesByHashtag(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route GET /pebble/category Pebble getPebbleByCategoryID
// getPebbleByCategoryID
//
// This endpoint return the pebble by category IDs.
//
// Endpoint: /pebble/category
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetPebbleByCategoryIDOpts
//     "$ref": "#/definitions/GetPebbleByCategoryIDOpts"
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
//  200: GetPebbleESResp description: OK
func (a *API) getPebbleByCategoryID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebbleByCategoryIDOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID = requestCTX.UserClaim.(*auth.UserClaim).ID
	}
	res, err := a.App.Elasticsearch.GetPebblesInfoByCategoryID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) searchPebbleByCaption(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.SearchPebbleByCaption
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	res, err := a.App.Content.SearchPebbleByCaption(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

// swagger:route POST /app/influencer/pebble Pebble createPebbleApp
// createPebbleApp
//
// This endpoint create pebble in app.
//
// Endpoint: /app/influencer/pebble
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreatePebbleAppOpts
//     "$ref": "#/definitions/CreatePebbleAppOpts"
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
//  200: CreatePebbleResp description: OK
func (a *API) createPebbleApp(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreatePebbleAppOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.CreatorID = id
	// s.InfluencerIDs = append(s.InfluencerIDs, id)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.CreatePebbleApp(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

// swagger:route POST /v2/app/influencer/pebble Pebble createPebbleAppV2
// createPebbleAppV2
//
// This endpoint create pebble app.
//
// Endpoint: /v2/app/influencer/pebble
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreatePebbleAppV2Opts
//     "$ref": "#/definitions/CreatePebbleAppV2Opts"
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
//  200: CreatePebbleResp description: OK
func (a *API) createPebbleAppV2(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreatePebbleAppV2Opts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.CreatorID = id
	// s.InfluencerIDs = append(s.InfluencerIDs, id)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.CreatePebbleAppV2(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

// swagger:route PUT /app/influencer/pebble Pebble editPebbleApp
// editPebbleApp
//
// This endpoint create pebble app.
//
// Endpoint: /app/influencer/pebble
//
// Method: PUT
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: EditPebbleAppOpts
//     "$ref": "#/definitions/EditPebbleAppOpts"
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
//  200: EditPebbleAppResp description: OK
func (a *API) editPebbleApp(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditPebbleAppOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.CreatorID = id
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Content.EditPebbleApp(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route GET /app/influencer/pebble Pebble getPebblesForCreator
// getPebblesForCreator
//
// This endpoint return pebbles for the creator.
//
// Endpoint: /app/influencer/pebble
//
// Method: GET
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetPebbleByInfluencerID
//     "$ref": "#/definitions/GetPebbleByInfluencerID"
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
//  200: GetPebbleESResp description: OK
func (a *API) getPebblesForCreator(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebbleByInfluencerID
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.InfluencerID = id.Hex()

	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.GetPebblesForCreator(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
