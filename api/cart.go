package api

import (
	"errors"
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pasztorpisti/qs"
	"github.com/vasupal1996/goerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createCart(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	user_id, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}

	resp, err := a.App.Cart.CreateCart(user_id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)

}

// swagger:route  POST /app/cart cart addToCart
// addToCart
//
// This endpoint will add product ot the cart.
//
// Endpoint: app/cart
//
// Method: POST
//
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// parameters:
// + name: body
//   in: body
//   description: Add to cart
//   schema:
//   type: AddToCartOpts
//     "$ref": "#/definitions/AddToCartOpts"
//   required: true
//
// responses:
//  400: AppErr description:Bad Request
//  403: AppErr description:Invalid User
//  200: addToCart description:OK
func (a *API) addToCart(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddToCartOpts
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
	resp, err := a.App.Cart.AddToCart(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  PUT /app/cart/item cart updateItemQty
// updateItemQty
//
// This endpoint will update the item quantity in the cart.
//
// Endpoint : /app/cart/item
// Method: PUT
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: UpdateItemQtyOpts
//     "$ref": "#/definitions/UpdateItemQtyOpts"
//   required: true
//
// parameters:
// + name: cookie
//   in: header
//   description: Customer login required for successful response.
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description: invalid user
//  200: addToCart description: OK
func (a *API) updateItemQty(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateItemQtyOpts
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
	resp, err := a.App.Cart.UpdateItemQty(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  GET /app/cart/{userID} cart getCartInfo
// getCartInfo
//
// This endpoint will get the cart information according to the user id.
//
// Endpoint : /app/cart/{userID}
//
// Method: GET
//
// parameters:
// + name: userID
//   in: path
//   schema:
//   type: ObjectID
//     "$ref": "#/definitions/ObjectID"
//   enum: 60b50277a97a2d73b211aec7
//   required: true
//
// parameters:
// + name: cookie
//   in: header
//   description: Customer login required for successful response.
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
//
// responses:
//  400: AppErr description: Bad Request
//  403: AppErr description: invalid user
//  200: GetCartInfoResp description: OK
func (a *API) getCartInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	if userID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	resp, err := a.App.Cart.GetCartInfo(userID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  POST /app/cart/address cart setCartAddress
// setCartAddress
//
// This endpoint will set the cart address of the user..
//
// Endpoint : /app/cart/address
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: AddressOpts
//     "$ref": "#/definitions/AddressOpts"
//   required: true
//
// parameters:
// + name: cookie
//   in: header
//   description: Customer login required for successful response.
//   required: true
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description: invalid user
//  200: description: true
func (a *API) setCartAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddressOpts
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
	err := a.App.Cart.SetCartAddress(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

// swagger:route  GET /app/cart/{userID}/checkout cart checkoutCart
// checkoutCart
//
// This endpoint will return checkout cart.
//
// Endpoint: /app/cart/{userID}/checkout
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
// parameters:
// + name: source
//   in: query
//   schema:
//   type: string
//   required: true
//
// parameters:
// + name: platform
//   in: query
//   schema:
//   type: string
//   enum: web, android, ios
//   required: true
//
// parameters:
// + name: isCOD
//   in: query
//   schema:
//   type: boolean
//   required: true
//
// parameters:
// + name: request_id
//   in: query
//   schema:
//   type: ObjectID
//   required: true
//
// parameters:
// + name: cookie
//   in: header
//   description: Login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
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
//  200: OrderInfo description: OK
func (a *API) checkoutCart(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	source := r.URL.Query().Get("source")
	if source == "" {
		requestCTX.SetErr(goerror.New("empty source in query", &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	platform := r.URL.Query().Get("platform")

	if platform != "web" && platform != "android" && platform != "ios" {
		requestCTX.SetErr(goerror.New("platform incorrect", &goerror.BadRequest), http.StatusBadRequest)
	}

	if id.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	cod := r.URL.Query().Get("cod")
	requestID := ""

	isCOD := true
	if cod != "true" {
		isCOD = false
	}
	if isCOD {
		requestID = r.URL.Query().Get("request_id")
		if requestID == "" {
			requestCTX.SetErr(goerror.New("request id is required for cod orders", &goerror.BadRequest), http.StatusBadRequest)
			return
		}
	}

	fullName := requestCTX.UserClaim.(*auth.UserClaim).FullName

	resp, err := a.App.Cart.CheckoutCart(id, source, platform, fullName, isCOD, requestID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  DELETE /cart/{userID} cart clearCart
// clearCart
//
// This endpoint will clear the cart of the user.
//
// Endpoint: /cart/{userID}
//
// Method: DELETE
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
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: Bad Request
//  200: description: true
func (a *API) clearCart(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	err = a.App.Cart.ClearCart(userID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)

}

// swagger:route  POST /app/cart/{userID}/coupon coupon applyCoupon
// applyCoupon
//
// This endpoint will successful apply the coupon on the product.
//
// Endpoint: /app/cart/{userID}/coupon
//
// Method: POST
//
// parameters:
// + name: userID
//   in: path
//   schema:
//   type: ObjectID
//     "$ref": "#/definitions/ObjectID"
//   enum: 61dd5c77c69b0de021ce1810
//   required: true
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: ApplyCouponOpts
//     "$ref": "#/definitions/ApplyCouponOpts"
//   required: true
//
// parameters:
// + name: cookie
//   in: header
//   description: Customer login required for successful response.
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
//  200: GetCartInfoResp description: OK
func (a *API) applyCoupon(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	userID, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	var s schema.ApplyCouponOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		fmt.Println(err, "input error")

		return
	}

	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if userID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	resp, err := a.App.Cart.ApplyCoupon(userID, &s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		fmt.Println(err, "app error")
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  DELETE /app/cart/{userID}/coupon coupon removeCoupon
// removeCoupon
//
// This endpoint will successful remove the coupon on the product.
//
// Endpoint: /app/cart/{userID}/coupon
//
// Method: DELETE
//
// parameters:
// + name: userID
//   in: path
//   schema:
//   type: ObjectID
//     "$ref": "#/definitions/ObjectID"
//   enum: 611ca3d6c2b96106c6c9ee47
//   required: true
//
//
//
// parameters:
// + name: cookie
//   in: header
//   description: Customer login required for successful response.
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
func (a *API) removeCoupon(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	if id.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	err = a.App.Cart.RemoveCoupon(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)

}

// swagger:route  GET /app/check/cod Check CODviaGoKwik checkCODEligiblity
// checkCODEligiblity
//
// This endpoint will check the COD eligiblity.
//
// Endpoint: /app/check/cod
//
// Method: GET
//
//
// parameters:
// + name: cookie
//   in: header
//   description:Login required for successful response.
//   required: true
//
//
// parameters:
// + name: auth token
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
func (a *API) checkCODEligiblity(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	email := requestCTX.UserClaim.(*auth.UserClaim).Email
	resp, err := a.App.Cart.CheckCODEligiblity(id, r.UserAgent(), r.RemoteAddr, email)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)

}

// swagger:route  GET /v2/app/cart/checkout CODviaGoKwik checkoutCartV2
// checkoutCartV2
//
// This endpoint will check the chekcout cart.
//
// Endpoint: /v2/app/cart/checkout
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: CheckoutOpts
//     "$ref": "#/definitions/CheckoutOpts"
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
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: AppErr description: BadRequest
//  403: AppErr description:Invalid User
//  200: OrderInfo description: OK
func (a *API) checkoutCartV2(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CheckoutOpts
	var err error

	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.ID, err = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.FullName = requestCTX.UserClaim.(*auth.UserClaim).FullName
	if s.IsCOD {
		if s.RequestID == "" {
			requestCTX.SetErr(goerror.New("request id is required for cod orders", &goerror.BadRequest), http.StatusBadRequest)
			return
		}
	}
	resp, err := a.App.Cart.CheckoutCartV2(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}
