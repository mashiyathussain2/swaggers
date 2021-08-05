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

func (a *API) createDiscount(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s []schema.CreateDiscountOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	for _, discount := range s {
		if errs := a.Validator.Validate(&discount); errs != nil {
			requestCTX.SetErrs(errs, http.StatusBadRequest)
			return
		}
	}
	var errs []error
	var results []interface{}
	for _, discount := range s {
		res, err := a.App.Discount.CreateDiscount(&discount)
		if err != nil {
			errs = append(errs, err)
		} else {
			results = append(results, res)
		}
	}
	requestCTX.SetCustomResponse(true, results, errs, http.StatusAccepted)
}
func (a *API) deactivateDiscount(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["discountID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["discountID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	err = a.App.Discount.DeactivateDiscount(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)

}
func (a *API) createSale(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateSaleOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Discount.CreateSale(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}
func (a *API) editSale(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditSaleOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Discount.EditSale(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) editSaleStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditSaleStatusOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Discount.EditSaleStatus(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getSales(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetSalesOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Discount.GetSales(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getDiscountInfoBySaleID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["saleID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["saleID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	res, err := a.App.Discount.GetDiscountAndCatalogInfoBySaleID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getAppActiveSale(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetAppActiveSaleOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	} else {
		if requestCTX.UserClaim != nil {
			s = schema.GetAppActiveSaleOpts{
				Genders: []string{requestCTX.UserClaim.(*auth.UserClaim).Gender},
			}
		}
	}
	res, err := a.App.Discount.GetAppActiveSale(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getSaleCatalogs(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetCatalogBySaleIDOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.GetCatalogBySaleID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) removeDiscountFromSale(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.RemoveDiscountFromSaleOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Discount.RemoveDiscountFromSale(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) changeSaleTime(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.ChangeSaleTimeOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Discount.ChangeSaleTime(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
