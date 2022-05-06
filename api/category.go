package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vasupal1996/goerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createCategory(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateCategoryOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Category.CreateCategory(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

func (a *API) editCategory(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditCategoryOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Category.EditCategory(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getCategory(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	res, err := a.App.Category.GetCategories()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getMainCategoryMap(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	res, err := a.App.Category.GetMainCategoriesMap()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route GET /app/category/lvl1 AppCategoryCatalog getParentCategory
// getParentCategory
//
// This endpoint return the parent app category.
//
// Endpoint: /app/category/lvl1
//
// Method: GET
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
//  200: GetParentCategoriesResp description: true
func (a *API) getParentCategory(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	res, err := a.App.Category.GetMainParentCategories()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route GET /app/category/{categoryID}/lvl2 AppCategoryCatalog getMainCategoryByParentID
// getMainCategoryByParentID
//
// This endpoint return main catalog by parent ID.
//
// Endpoint: /app/category/{categoryID}/lvl2
//
// Method: GET
//
// parameters:
// + name: categoryID
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
//  200: GetMainCategoriesByParentIDResp description: OK
func (a *API) getMainCategoryByParentID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	parentID, err := primitive.ObjectIDFromHex(mux.Vars(r)["categoryID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["categoryID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	res, err := a.App.Category.GetMainCategoriesByParentID(parentID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

// swagger:route GET /app/category/{categoryID}/lvl3 AppCategoryCatalog getSubCatergoryByParentID
// getSubCatergoryByParentID
//
// This endpoint return sub category by parent ID.
//
// Endpoint: /app/category/{categoryID}/lvl3
//
// Method: GET
//
// parameters:
// + name: categoryID
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
//  200: GetSubCategoriesByParentIDResp description: OK
func (a *API) getSubCatergoryByParentID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	parentID, err := primitive.ObjectIDFromHex(mux.Vars(r)["categoryID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["categoryID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	res, err := a.App.Category.GetSubCategoriesByParentID(parentID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getCategoryPath(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	categoryID, err := primitive.ObjectIDFromHex(mux.Vars(r)["categoryID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["categoryID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	res, err := a.App.Category.GetCategoryPath(categoryID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getAncestorsByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	categoryID, err := primitive.ObjectIDFromHex(mux.Vars(r)["categoryID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["categoryID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	res, err := a.App.Category.GetAncestorsByID(categoryID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
