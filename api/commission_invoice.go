package api

import (
	"go-app/server/auth"
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
	err = a.App.CommissionInvoice.CreateCommissionInvoice(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) downloadCommissionInvoice(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	invoiceNo := r.URL.Query().Get("invoiceNo")
	userID, _ := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	file, name, err := a.App.CommissionInvoice.GetInvoicePDF(userID, invoiceNo)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+name+".pdf")
	w.Header().Set("Content-type", "application/zip")
	// z := zip.NewWriter(w)
	// defer z.Close()
	// for i, pdf := range files {
	// 	zf, _ := z.Create(names[i])
	// 	f, _ := os.Open(names[i])
	// 	defer f.Close()
	// 	io.Copy(zf, pdf)
	// }
	io.Copy(w, file)
	return
}
