package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/vasupal1996/goerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:route  PUT /app/wishlist Wishlist addToWishlist
// addToWishlist
//
// This endpoint will add the product to the wishlist.
//
// Endpoint: /app/wishlist
//
// Method: PUT
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: AddToWishlistOpts
//     "$ref": "#/definitions/AddToWishlistOpts"
//   required: true
//
// parameters:
// + name: cookie
//   in: header
//   description: Customer login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
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
//  403: AppErr description: Invalid User
//  200: description: OK
func (a *API) addToWishlist(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddToWishlistOpts
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
	err := a.App.Wishlist.AddToWishlist(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

// swagger:route  GET /app/wishlist/{userID} Wishlist getWishlist
// getWishlist
//
// This endpoint will return the wishlist product of the user.
//
// Endpoint: /app/wishlist/{userID}
//
// Method: GET
//
// parameters:
// + name: userID
//   in: path
//   schema:
//   type: ObjectID
//     "$ref": "#/definitions/ObjectID"
//   enum: 6065d4503824bf77961c21ae
//   required: true
//
//
// parameters:
// + name: cookie
//   in: header
//   description: Customer login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
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
//  403: AppErr description: Invalid User
//  200: description: true
func (a *API) getWishlist(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	if id.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Wishlist.GetWishlistMap(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  DELETE /app/wishlist Wishlist removeFromWishlist
// removeFromWishlist
//
// This endpoint will delete the product from the wishlist.
//
// Endpoint: /app/wishlist
//
// Method: DELETE
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: RemoveFromWishlistOpts
//     "$ref": "#/definitions/RemoveFromWishlistOpts"
//   required: true
//
// parameters:
// + name: cookie
//   in: header
//   description: Customer login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
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
//  403: AppErr description: Invalid User
//  200: description: true
func (a *API) removeFromWishlist(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.RemoveFromWishlistOpts
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
	err := a.App.Wishlist.RemoveFromWishlist(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
