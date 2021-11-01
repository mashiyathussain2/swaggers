package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pasztorpisti/qs"
	"github.com/vasupal1996/goerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createCatalog(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateCatalogOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.KeeperCatalog.CreateCatalog(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}

func (a *API) editCatalog(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditCatalogOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.KeeperCatalog.EditCatalog(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) addVariants(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddVariantOpts

	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.KeeperCatalog.AddVariant(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) keeperSearchCatalog(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.KeeperSearchCatalogOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.KeeperCatalog.KeeperSearchCatalog(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) deleteVariant(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.DeleteVariantOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.KeeperCatalog.DeleteVariant(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) updateCatalogStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCatalogStatusOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.UpdateCatalogStatus(&s)
	if err != nil {
		if resp != nil {
			requestCTX.SetCustomResponse(false, resp, err, http.StatusBadRequest)
			return
		}
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) editVariantSKU(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditVariantSKU
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.EditVariantSKU(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) addCatalogContent(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddCatalogContentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, errs := a.App.KeeperCatalog.AddCatalogContent(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) addCatalogContentImage(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddCatalogContentImageOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	errs := a.App.KeeperCatalog.AddCatalogContentImage(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getCatalogsByFilter(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetCatalogsByFilterOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.GetCatalogsByFilter(&s)
	if err != nil {

		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCatalogBySlug(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	resp, err := a.App.KeeperCatalog.GetCatalogBySlug(slug)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getBasicCatalogInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetBasicCatalogFilter
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.GetBasicCatalogInfo(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusInternalServerError)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCatalogFilter(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	resp, err := a.App.KeeperCatalog.GetCatalogFilter()
	if err != nil {
		requestCTX.SetErr(err, http.StatusInternalServerError)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCatalogVariant(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	cat_id, err := primitive.ObjectIDFromHex(mux.Vars(r)["catalogID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["catalogID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}

	var_id, err := primitive.ObjectIDFromHex(mux.Vars(r)["variantID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["variantID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}

	resp, err := a.App.KeeperCatalog.GetCatalogVariant(cat_id, var_id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)

}

func (a *API) getAllCatalogInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	cat_id, err := primitive.ObjectIDFromHex(mux.Vars(r)["catalogID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["catalogID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.GetAllCatalogInfo(cat_id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCatalogBasicByIds(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetCatalogByIDFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	resp, err := a.App.Elasticsearch.GetCatalogByIDs(s.IDs)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getSimilarProducts(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetSimilarProducts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	fmt.Println(s.Query)
	resp, err := a.App.Elasticsearch.GetSimilarProducts(s.Query)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCatalogInfoById(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["catalogID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["catalogID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	resp, err := a.App.Elasticsearch.GetCatalogInfoByID(id.Hex())
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCatalogByCategoryID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["categoryID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["catalogID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	var s schema.GetCatalogByCategoryIDOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.CategoryID = id.Hex()
	resp, err := a.App.Elasticsearch.GetCatalogInfoByCategoryID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) removeContentfromCatalog(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.RemoveContentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.KeeperCatalog.RemoveContent(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getCatalogContent(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["catalogID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["catalogID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.GetKeeperCatalogContent(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getPebbleCatalogInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebbleCatalogInfoByIDs
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.GetPebbleCatalogInfo(s.IDs)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) search(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.SearchOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.SearchBrandCatalogInfluencerContent(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) searchCatalog(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.SearchOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.SearchCatalog(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) searchDiscover(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.SearchOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.SearchDiscover(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) searchBrand(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.SearchOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.SearchBrand(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) searchInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.SearchOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.SearchInfluencer(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) searchSeries(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.SearchOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Elasticsearch.SearchSeries(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getCatalogInfoByBrandId(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetCatalogByBrandIDOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	resp, err := a.App.Elasticsearch.GetCatalogByBrandID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

// func (a *API) bulkAddCatalogCSV(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
// 	file, _, err := r.FormFile("myFile")
// 	if err != nil {
// 		requestCTX.SetErr(err, http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()
// 	returnFile, err := a.App.KeeperCatalog.BulkAddCatalogsCSV(file)
// 	if err != nil {
// 		requestCTX.SetErr(err, http.StatusBadRequest)
// 		return
// 	}
// 	w.Header().Set("Content-Disposition", "attachment; filename=WHATEVER_YOU_WANT.xlsx")
// 	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
// 	// w.Header().Set("Content-Length", fileSize)
// 	// t := bytes.NewReader(fileContents)
// 	// t.Seek(0, 0)
// 	io.Copy(w, returnFile)
// 	return
// }

func (a *API) bulkAddCatalogJSON(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	var s []schema.BulkUploadCatalogJSONOpts
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
	resp, err := a.App.KeeperCatalog.BulkAddCatalogsJSON(s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
	return
}
func (a *API) getCollectionCatalogByIDs(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetCollectionCatalogByIDs
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	resp, err := a.App.Elasticsearch.GetCollectionCatalogByIDs(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
	return
}
