package api

import (
	"errors"
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:parameters error commonError
type CommonError struct {
	// Status of the error
	Status int `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:parameters error commonError
type AddErrorBody struct {
	// - name: body
	//  in: body
	//  description: erroror
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/CommonError"
	//  required: true
	Body CommonError `json:"body"`
}

// swagger:route  POST /customer/email/login login loginViaEmail
// User login via email
//
// parameters:
// + name: body
//   in: body
//   description: Login Via Email
//   schema:
//   type: EmailLoginCustomerOpts
//     "$ref": "#/definitions/EmailLoginCustomerOpts"
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
//  400: CommonError description: Error
//  200: SuccessfulLogin
func (a *API) loginViaEmail(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EmailLoginCustomerOpts
	var isWeb bool
	isWeb, _ = strconv.ParseBool(r.URL.Query().Get("isWeb"))
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Customer.Login(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	token, err := a.TokenAuth.SignToken(res)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if isWeb {
		if err := a.SessionAuth.Create(token, w); err != nil {
			requestCTX.SetErr(fmt.Errorf("failed to login user: %s", err), http.StatusInternalServerError)
			return
		}
		requestCTX.SetAppResponse(true, http.StatusOK)
		return
	}
	requestCTX.SetAppResponse(token, http.StatusOK)
}

// swagger:route  POST /customer/email/signup signup signUpViaEmail
// signUpViaEmail
//
// User Signup via email
//
// Endpoint: /customer/email/signup
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   description: Signup Via Email
//   schema:
//   type: CreateUserOpts
//     "$ref": "#/definitions/CreateUserOpts"
//   required: true
//
// parameters:
// + name: isWeb
//   in: query
//   description: If value is set to True, token is omitted from response
//   schema:
//   type: boolean
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
//  400: CommonError description: Error
//  200: SuccessfulLogin description: Success
func (a *API) signUpViaEmail(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateUserOpts
	var isWeb bool
	isWeb, _ = strconv.ParseBool(r.URL.Query().Get("isWeb"))
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Customer.SignUp(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	token, err := a.TokenAuth.SignToken(res)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if isWeb {
		if err := a.SessionAuth.Create(token, w); err != nil {
			requestCTX.SetErr(fmt.Errorf("failed to login user: %s", err), http.StatusInternalServerError)
			return
		}
		requestCTX.SetAppResponse(true, http.StatusOK)
		return
	}

	requestCTX.SetAppResponse(token, http.StatusOK)
}

// swagger:route  PUT /customer customer updateCustomerInfo
// updateCustomerInfo
//
// This endpoint will update the customer information.
//
// Endpoint: /customer
//
// Method: PUT
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: UpdateCustomerOpts
//     "$ref": "#/definitions/UpdateCustomerOpts"
//   required: true
//
// parameters:
// + name: isWeb
//   in: query
//   description: If value is set to True, token is omitted from response
//   schema:
//   type: boolean
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
//  200: SuccessfulLogin description: Success
func (a *API) updateCustomerInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCustomerOpts
	var isWeb bool
	isWeb, _ = strconv.ParseBool(r.URL.Query().Get("isWeb"))
	resp := make(map[string]interface{})
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	s.ID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	claim, err := a.App.Customer.UpdateCustomer(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	resp["user"] = claim
	token, err := a.TokenAuth.SignToken(claim)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if isWeb {
		if err := a.SessionAuth.Create(token, w); err != nil {
			requestCTX.SetErr(fmt.Errorf("failed to login user: %s", err), http.StatusInternalServerError)
			return
		}
		requestCTX.SetAppResponse(true, http.StatusOK)
		return
	}
	resp["token"] = token
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  POST /app/customer/influencer/follow followUnfollow followInfluencer
// followInfluencer
//
// This endpoint will follow the influencer.
//
// Endpoint: /app/customer/influencer/follow
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: AddInfluencerFollowerOpts
//     "$ref": "#/definitions/AddInfluencerFollowerOpts"
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
//  200: description: true
func (a *API) followInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddInfluencerFollowerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.CustomerID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Influencer.AddFollower(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /app/customer/brand/follow followUnfollow followBrand
// followBrand
//
// This endpoint will follow the brand.
//
// Endpoint: /app/customer/brand/follow
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: AddBrandFollowerOpts
//     "$ref": "#/definitions/AddBrandFollowerOpts"
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
//  200: description: true
func (a *API) followBrand(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddBrandFollowerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.CustomerID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Brand.AddFollower(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /app/customer/influencer/unfollow followUnfollow unFollowInfluencer
// unFollowInfluencer
//
// This endpoint will unfollow the influencer.
//
// Endpoint: /app/customer/influencer/unfollow
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: AddInfluencerFollowerOpts
//     "$ref": "#/definitions/AddInfluencerFollowerOpts"
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
//  200: description: true
func (a *API) unFollowInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddInfluencerFollowerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.CustomerID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Influencer.RemoveFollower(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /app/customer/brand/unfollow followUnfollow unFollowBrand
// unFollowBrand
//
// This endpoint will unfollow the brand.
//
// Endpoint: /app/customer/brand/unfollow
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: AddBrandFollowerOpts
//     "$ref": "#/definitions/AddBrandFollowerOpts"
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
//  200: description: true
func (a *API) unFollowBrand(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddBrandFollowerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.CustomerID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Brand.RemoveFollower(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  PUT /customer/address customer addAddress
// addAddress
//
// This endpoint will add the address of the customer.
//
// Endpoint: /customer/address
//
// Method: PUT
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: AddAddressOpts
//     "$ref": "#/definitions/AddAddressOpts"
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
//  400: CommonError description: Error
//  200: AddAddressResp description: Success
func (a *API) addAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddAddressOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.UserID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.Customer.AddAddress(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  GET /customer/{userID}/address customer GetAddress
// GetAddress
//
// This endpoint will return the address of the user.
//
// Endpoint: /customer/{userID}/address
//
// Method: GET
//
// parameters:
// + name: userID
//   in: path
//   schema:
//   type: ObjectID
//     "$ref": "#/definitions/ObjectID"
//   enum: 60b50277a97a2d73b211aec8
//   required: true
//
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
//
// responses:
//  400: CommonError description: Error
//  200: getAddress description: Success
func (a *API) getAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if id.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	resp, err := a.App.Customer.GetAddresses(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// swagger:route  POST /app/customer/{customerID} getCustomerInfo getCustomerInfo
//
// This endpoint will return the address of the user.
//
// Endpoint: /app/customer/{customerID}
//
// Method: POST
//
// parameters:
// + name: customerID
//   in: path
//   schema:
//   type: ObjectID
//     "$ref": "#/definitions/ObjectID"
//   enum: 60b50277a97a2d73b211aec7
//   required: true
//   examples:
//       customerId:
//         summary: Example of a customer ID
//         value: [60b50277a97a2d73b211aec7]
//
//
// consumes:
//         - application/json
//
// produces:
//         - application/json
//
// responses:
//  400: CommonError description: Error
//  200: getCustomerInfo description: Success
func (a *API) getCustomerInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["customerID"])
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if id.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CustomerID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	res, err := a.App.Customer.GetAppCustomerInfo(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  DELETE /customer/address customer removeAddress
// removeAddress
//
// This endpoint will delete the address of the customer.
//
// Endpoint: /customer/address
//
// Method: DELETE
//
// parameters:
// + name: user_id
//   in: body
//   name: address_id
//   type: string
//   enum: address_id:60b50277a97a2d73b211aec7
//   enum: user_id:60b50277a97a2d73b211aec7
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
//  400: CommonError description: Error
//  200: description: Success
func (a *API) removeAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	userID, err := primitive.ObjectIDFromHex(r.FormValue("user_id"))
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if userID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	addressID, err := primitive.ObjectIDFromHex(r.FormValue("address_id"))
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	err = a.App.Customer.RemoveAddress(userID, addressID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

// swagger:route  PUT /customer/address/edit customer editAddress
// editAddress
//
// This endpoint will edit the address.
//
// Endpoint : /customer/address/edit
// Method: PUT
//
// parameters:
// + name: body
//   in: body
//   description: Edit Address
//   schema:
//   type: EditAddressOpts
//     "$ref": "#/definitions/EditAddressOpts"
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
//  400: CommonError description: Error
//  200: description: true
func (a *API) editAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditAddressOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	//TODO:WHY this not just check?
	s.UserID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Customer.EditAddress(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
