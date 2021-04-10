package api

import (
	"errors"
	"fmt"
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vasupal1996/goerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createCart(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	cat_id, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}

	resp, err := a.App.Cart.CreateCart(cat_id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)

}

func (a *API) addToCart(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddToCartOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CartID {
		requestCTX.SetErr(errors.New("invalid cart"), http.StatusForbidden)
		return
	}
	resp, err := a.App.Cart.AddToCart(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) updateItemQty(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateItemQtyOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CartID {
		requestCTX.SetErr(errors.New("invalid cart"), http.StatusForbidden)
		return
	}
	resp, err := a.App.Cart.UpdateItemQty(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCartInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	if userID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	resp, err := a.App.Cart.GetCartInfo(userID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) setCartAddress(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddressOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CartID {
		requestCTX.SetErr(errors.New("invalid cart"), http.StatusForbidden)
		return
	}
	err := a.App.Cart.SetCartAddress(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) checkoutCart(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["cartID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["cartID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	source := r.URL.Query().Get("source")
	if source == "" {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid source :%s in url", mux.Vars(r)["source"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	if id.Hex() != requestCTX.UserClaim.(*auth.UserClaim).CartID {
		requestCTX.SetErr(errors.New("invalid cart"), http.StatusForbidden)
		return
	}
	resp, err := a.App.Cart.CheckoutCart(id, source)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)

}
