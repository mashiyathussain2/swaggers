package api

import (
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/pasztorpisti/qs"
	"github.com/pkg/errors"
)

// swagger:route POST /app/review AppReview createReview
// createReview
//
// This endpoint will post the app review.
//
// Endpoint: /app/review
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CreateReviewStoryOpts
//     "$ref": "#/definitions/CreateReviewStoryOpts"
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
//  200: CreateReviewStoryResp description: OK
func (a *API) createReview(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateReviewStoryOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.UserID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.Errorf("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Review.CreateReviewStory(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

// swagger:route GET /app/review/catalog AppReview getReviewsByCatalogID
// getReviewsByCatalogID
//
// This endpoint will return the review by catalog ID.
//
// Endpoint: /app/review/catalog
//
// Method: GET
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetReviewsByCatalogIDFilter
//     "$ref": "#/definitions/GetReviewsByCatalogIDFilter"
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
//  200: GetReviewsByCatalogIDResp description: OK
func (a *API) getReviewsByCatalogID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetReviewsByCatalogIDFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.GetReviewsByCatalogID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
