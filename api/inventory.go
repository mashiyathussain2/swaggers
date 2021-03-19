package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/handler"
	"net/http"

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
