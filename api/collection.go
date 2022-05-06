package api

import (
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

func (a *API) createCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, errs := a.App.Collection.CreateCollection(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}
func (a *API) deleteCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["collectionID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["collectionID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	err = a.App.Collection.DeleteCollection(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) addSubCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddSubCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, errs := a.App.Collection.AddSubCollection(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}
func (a *API) deleteSubCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	collID, err := primitive.ObjectIDFromHex(mux.Vars(r)["collectionID"])

	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["collectionID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	subCollID, err := primitive.ObjectIDFromHex(mux.Vars(r)["subCollectionID"])

	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid subcollection id:%s in url", mux.Vars(r)["subCollectionID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	err = a.App.Collection.DeleteSubCollection(collID, subCollID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) editCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Collection.EditCollection(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
func (a *API) updateSubCollectionImage(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateSubCollectionImageOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Collection.UpdateSubCollectionImage(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) addCatalogsToSubCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCatalogsInSubCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	errs := a.App.Collection.AddCatalogsToSubCollection(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) removeCatalogsFromSubCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCatalogsInSubCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	errs := a.App.Collection.RemoveCatalogsFromSubCollection(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getCollections(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetCollectionsKeeperFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	resp, err := a.App.Collection.GetCollections(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)

}

// swagger:route GET /app/collections AppCollectionCatalogV2 getActiveCollections
// getActiveCollections
//
// This endpoint will return the collections.
//
// Endpoint: /app/collections
//
// Method: GET
//
// parameters:
// + name: body
//   in: body
//   schema:
//   type: GetActiveCollectionsOpts
//     "$ref": "#/definitions/GetActiveCollectionsOpts"
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
//  200: GetCollectionESResp description: true
func (a *API) getActiveCollections(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetActiveCollectionsOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.Gender = requestCTX.UserClaim.(*auth.UserClaim).Gender
	}
	resp, err := a.App.Elasticsearch.GetActiveCollections(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)

}

func (a *API) getWebActiveCollections(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetActiveCollectionsOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.Gender = requestCTX.UserClaim.(*auth.UserClaim).Gender
	}
	resp, err := a.App.Elasticsearch.GetActiveCollections(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)

}

func (a *API) updateCollectionStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCollectionStatus
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Collection.UpdateCollectionStatus(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) setFeaturedCatalogs(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.SetFeaturedCatalogs
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Collection.SetFeaturedCatalogs(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

// swagger:route GET /app/subcollection AppCollectionCatalogV2 GetCatalogBySubCollectionID
// GetCatalogBySubCollectionID
//
// This endpoint return catalog by sub collection ID.
//
// Endpoint: /app/subcollection
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
//  200: GetCatalogsBySubCollectionResp description: OK
func (a *API) GetCatalogBySubCollectionID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	collID, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
	fmt.Println(collID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	resp, err := a.App.Collection.GetCatalogsBySubCollection(collID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}
