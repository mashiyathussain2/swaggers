package api

import (
	"go-app/schema"
	"go-app/server/handler"
	"net/http"
)

func (a *API) expressCheckout(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ExpressCheckoutOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.ExpressCheckout.ExpressCheckoutComplete(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}
