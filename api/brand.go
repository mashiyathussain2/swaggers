package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pasztorpisti/qs"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createbrand(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateBrandOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.CreateBrand(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) createBrandAdminUser(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateBrandAdminUserOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.CreateBrandAdminUser(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getBrandsById(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetBrandsByIDOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.GetBrandsByID(s.IDs)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) editbrand(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditBrandOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.EditBrand(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /brand/{brandID} Brand getBrandByID
// getBrandByID
//
// This endpoint will return brand details.
//
// Endpoint: /brand/{brandID}
//
// Method: POST
//
//
// parameters:
// + name: brandID
//   in: path
//   schema:
//   type: string
//     "$ref": "#/definitions/ObjectID"
//   enum: 60b50277a97a2d73b211aec7
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
//  200: getBrandByID description: Success
func (a *API) getBrandByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["brandID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid brand id:%s in url", mux.Vars(r)["brandID"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.GetBrandByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) checkBrandByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["brandID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid brand id:%s in url", mux.Vars(r)["brandID"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.CheckBrandByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getBrands(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	res, err := a.App.Brand.GetBrands()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /app/brand/basic Brand getBrandsBasic
// getBrandsBasic
//
// This endpoint will return brand basic information.
//
// Endpoint: /app/brand/basic
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetBrandsByIDBasicOpts
//     "$ref": "#/definitions/GetBrandsByIDBasicOpts"
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
//  200: GetBrandBasicESEesp description: OK
func (a *API) getBrandsBasic(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetBrandsByIDBasicOpts
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
	res, err := a.App.Elasticsearch.GetBrandsByIDBasic(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /app/brand/{brandID} Brand getBrandInfo
// getBrandInfo
//
// This endpoint will return information of one brand.
//
// Endpoint: /app/brand/{brandID}
//
// Method: GET
//
// parameters:
// + name: brandID
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
//   type: GetBrandsInfoByIDOpts
//     "$ref": "#/definitions/GetBrandsInfoByIDOpts"
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
//  200: GetBrandInfoEsResp description: OK
func (a *API) getBrandInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["brandID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid brand id:%s in url", mux.Vars(r)["brandID"]), http.StatusBadRequest)
		return
	}
	var userID primitive.ObjectID
	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetBrandInfoByID(&schema.GetBrandsInfoByIDOpts{ID: id, CustomerID: userID})
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /app/brands Brand getActiveBrandsList
// getActiveBrandsList
//
// This endpoint will return active brand list.
//
// Endpoint: /app/brands
//
// Method: GET
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetBrandsListOpts
//     "$ref": "#/definitions/GetBrandsListOpts"
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
//  200: GetActiveBrandsListESEesp description: OK
func (a *API) getActiveBrandsList(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetBrandsListOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.GetBrandsList(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /brand/user/login BrandDashAPIs brandUserLogin
// brandUserLogin
//
// This endpoint will login the brand user.
//
// Endpoint: /brand/user/login
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: BrandUserLoginOpts
//     "$ref": "#/definitions/BrandUserLoginOpts"
//   required: true
//
//
// parameters:
// + name: returnToken
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
//  200: description: token
func (a *API) brandUserLogin(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.BrandUserLoginOpts
	var returnToken bool
	returnToken, _ = strconv.ParseBool(r.URL.Query().Get("returnToken"))

	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}

	res, err := a.App.Brand.BrandUserLogin(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	token, err := a.TokenAuth.SignToken(res)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if !returnToken {
		if err := a.SessionAuth.Create(token, w); err != nil {
			requestCTX.SetErr(fmt.Errorf("failed to login user: %s", err), http.StatusInternalServerError)
			return
		}
		requestCTX.SetAppResponse(true, http.StatusOK)
		return
	}
	requestCTX.SetAppResponse(token, http.StatusOK)
}

// swagger:route  POST /brand/user/forgot-password BrandDashAPIs brandUserForgotPassword
// brandUserForgotPassword
//
// This endpoint will help the brand to recover forget password.
//
// Endpoint: /brand/user/forgot-password
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: ForgotPasswordOpts
//     "$ref": "#/definitions/ForgotPasswordOpts"
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
func (a *API) brandUserForgotPassword(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ForgotPasswordOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.ForgotPassword(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  POST /brand/user/reset-password BrandDashAPIs brandUserResetPassword
// brandUserResetPassword
//
// This endpoint will reset the brand password.
//
// Endpoint: /brand/user/reset-password
//
// Method: POST
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: ResetPasswordOpts
//     "$ref": "#/definitions/ResetPasswordOpts"
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
func (a *API) brandUserResetPassword(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ResetPasswordOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Brand.ResetPassword(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /brand/check/username CheckUsername checkBrandUsernameExists
// checkBrandUsernameExists
//
// This endpoint will check the brand username exists or not.
//
// Endpoint: /brand/check/username
//
// Method: GET
//
// parameters:
// + name: username
//   in: query
//   schema:
//   enum: falthead
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
func (a *API) checkBrandUsernameExists(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		requestCTX.SetErr(errors.Errorf("username cannot be empty nil"), http.StatusBadRequest)
		return
	}
	err := a.App.Brand.CheckBrandUsernameExists(username, nil)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

// swagger:route  POST /app/brand/username/basic Brand getBrandsBasicByUsername
// getBrandsBasicByUsername
//
// This endpoint will return brand basic detail by username.
//
// Endpoint: /app/brand/username/basic
//
// Method: POST
//
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetBrandsByUsernameBasicOpts
//     "$ref": "#/definitions/GetBrandsByUsernameBasicOpts"
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
//  200: GetBrandBasicESEesp description: OK
func (a *API) getBrandsBasicByUsername(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetBrandsByUsernameBasicOpts
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
	res, err := a.App.Elasticsearch.GetBrandsByUsernameBasic(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route  GET /app/brand/username/{username} Brand getBrandInfoByUsername
// getBrandInfoByUsername
//
// This endpoint will return brand information by the username of the brand.
//
// Endpoint: /app/brand/username/{username}
//
// Method: GET
//
//
// parameters:
// + name: username
//   in: path
//   schema:
//   type: string
//   enum: vasu_pal_1
//   required: true
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetBrandsInfoByUsernameOpts
//     "$ref": "#/definitions/GetBrandsInfoByUsernameOpts"
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
//  200: GetBrandInfoEsResp description: OK
func (a *API) getBrandInfoByUsername(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	if username == "" {
		requestCTX.SetErr(errors.Errorf("invalid brand id:%s in url", mux.Vars(r)["username"]), http.StatusBadRequest)
		return
	}
	var userID primitive.ObjectID
	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetBrandInfoByUsername(&schema.GetBrandsInfoByUsernameOpts{Username: username, CustomerID: userID})
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
