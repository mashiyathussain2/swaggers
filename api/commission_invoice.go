package api

import (
	"go-app/server/handler"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) generateCommissionInvoice(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	_, err = a.App.CommissionInvoice.CreateCommissionInvoice(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) downloadCommissionInvoice(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	invoiceNo := r.URL.Query().Get("invoiceNo")
	// userID, _ := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	file, name, err := a.App.CommissionInvoice.GetInvoicePDF(invoiceNo)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	w.Header().Set("Content-type", "application/pdf")
	io.Copy(w, file)
	return
}
