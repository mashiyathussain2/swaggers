package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/handler"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vasupal1996/goerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) updateInventory(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateInventoryOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Inventory.UpdateInventory(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) unicommerceUpdateInventory(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UnicommerceUpdateInventoryByInventoryIDsOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Inventory.UnicommerceUpdateInventoryByVariantIDs(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) setOutOfStock(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	invID, err := GetObjectID(r)

	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", r.URL.Query().Get("id")), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	if invID == primitive.NilObjectID {
		requestCTX.SetErr(goerror.New("id cannot be empty", &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	err = a.App.Inventory.SetOutOfStock(invID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) checkInventoryExists(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	cat_id, err := primitive.ObjectIDFromHex(mux.Vars(r)["catalogID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid catalog id:%s in url", mux.Vars(r)["catalogID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	var_id, err := primitive.ObjectIDFromHex(mux.Vars(r)["variantID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid variant id:%s in url", mux.Vars(r)["variantID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}

	qty, err := strconv.Atoi(mux.Vars(r)["quantity"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid quantity :%s in url", mux.Vars(r)["quantity"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}

	resp, err := a.App.Inventory.CheckInventoryExists(cat_id, var_id, qty)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) updateInventoryInternal(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	var s []schema.UpdateInventoryCVOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	for _, inventory := range s {
		if errs := a.Validator.Validate(&inventory); errs != nil {
			requestCTX.SetErrs(errs, http.StatusBadRequest)
			return
		}
	}

	err := a.App.Inventory.UpdateInventoryInternal(s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) updateInventoryBySKU(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s []schema.UpdateInventoryBySKUOpt
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	for _, inv := range s {
		if errs := a.Validator.Validate(&inv); errs != nil {
			requestCTX.SetErrs(errs, http.StatusBadRequest)
			return
		}
	}

	notUniqueSku, err := a.App.Inventory.UpdateInventorybySKUs(s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(notUniqueSku, http.StatusOK)
}
