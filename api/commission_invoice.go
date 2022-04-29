package api

import (
	"go-app/server/handler"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) generateCommissionInvoice(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	err = a.App.CommissionInvoice.GenerateCommissionInvoice(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
