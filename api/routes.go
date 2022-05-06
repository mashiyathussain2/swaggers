package api

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {

	//KEEPER CATEGORY
	a.Router.Root.Handle("/", a.requestWithSudoHandler(a.home)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/category", a.requestWithSudoHandler(a.createCategory)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/category", a.requestWithSudoHandler(a.editCategory)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/category", a.requestWithSudoHandler(a.getCategory)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/category/main", a.requestWithSudoHandler(a.getMainCategoryMap)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/category/{categoryID}/path", a.requestWithSudoHandler(a.getCategoryPath)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/category/{categoryID}/ancestors", a.requestWithSudoHandler(a.getAncestorsByID)).Methods("GET")

	//KEEPER CATALOG
	a.Router.APIRoot.Handle("/keeper/catalog", a.requestWithSudoHandler(a.getCatalogsByFilter)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/catalog", a.requestWithSudoHandler(a.createCatalog)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog", a.requestWithSudoHandler(a.editCatalog)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/catalog/variant", a.requestWithSudoHandler(a.addVariants)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog/status", a.requestWithSudoHandler(a.updateCatalogStatus)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog/variant", a.requestWithSudoHandler(a.deleteVariant)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/catalog/variant", a.requestWithSudoHandler(a.editVariantSKU)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/catalog/search", a.requestWithSudoHandler(a.keeperSearchCatalog)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/catalog/get", a.requestWithSudoHandler(a.getCatalogsByFilter)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog/content/video", a.requestWithSudoHandler(a.addCatalogContent)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog/content/image", a.requestWithSudoHandler(a.addCatalogContentImage)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog/content/{catalogID}", a.requestWithSudoHandler(a.getCatalogContent)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/catalog/content", a.requestWithSudoHandler(a.removeContentfromCatalog)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/catalog/slug/{slug}", a.requestWithSudoHandler(a.getCatalogBySlug)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/catalog", a.requestWithSudoHandler(a.getCatalogsByFilter)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/catalog/{catalogID}/get", a.requestWithSudoHandler(a.getAllCatalogInfo)).Methods("GET")

	// UNICOMMERCE APIS
	a.Router.APIRoot.Handle("/unicommerce/catalog/count", a.requestWithInternalHandler(a.getCatalogCount)).Methods("POST")
	a.Router.APIRoot.Handle("/unicommerce/catalog", a.requestWithInternalHandler(a.getCatalogs)).Methods("POST")

	//INTERNAL APIS
	a.Router.APIRoot.Handle("/keeper/catalog/{catalogID}", a.requestWithInternalHandler(a.getAllCatalogInfo)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/catalog/{catalogID}/variant/{variantID}", a.requestWithInternalHandler(a.getCatalogVariant)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/catalog/get/ids", a.requestWithInternalHandler(a.getPebbleCatalogInfo)).Methods("POST")

	//KEEPER GROUP
	a.Router.APIRoot.Handle("/keeper/group", a.requestWithSudoHandler(a.createCatalogGroup)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/group", a.requestWithSudoHandler(a.getCatalogsByGroupID)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/groups", a.requestWithSudoHandler(a.getGroups)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/groups/catalog", a.requestWithSudoHandler(a.keeperGetGroupsByCatalogID)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/group/catalogs", a.requestWithSudoHandler(a.addCatalogsInTheGroup)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/group", a.requestWithSudoHandler(a.editGroup)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/group/status", a.requestWithSudoHandler(a.updateGroupStatus)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/groups/catalog-name", a.requestWithSudoHandler(a.getGroupsByCatalogName)).Methods("GET")

	//KEEPER INVENTORY
	a.Router.APIRoot.Handle("/keeper/inventory", a.requestWithSudoHandler(a.updateInventory)).Methods("POST")
	a.Router.APIRoot.Handle("/unicommerce/inventory", a.requestWithInternalHandler(a.unicommerceUpdateInventory)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/inventory/outofstock", a.requestWithSudoHandler(a.setOutOfStock)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/inventory/catalog/{catalogID}/variant/{variantID}/quantity/{quantity}", a.requestWithSudoHandler(a.checkInventoryExists)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/inventory/sku", a.requestWithSudoHandler(a.updateInventoryBySKU)).Methods("POST")

	//Internal API for Inventory
	a.Router.APIRoot.Handle("/inventory", a.requestWithInternalHandler(a.updateInventoryInternal)).Methods("POST")

	//KEEPER COLLECTION
	a.Router.APIRoot.Handle("/keeper/collection", a.requestWithSudoHandler(a.createCollection)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/collection", a.requestWithSudoHandler(a.editCollection)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/collections", a.requestWithSudoHandler(a.getCollections)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/collection/status", a.requestWithSudoHandler(a.updateCollectionStatus)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/collection/{collectionID}", a.requestWithSudoHandler(a.deleteCollection)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/collection/subcollection", a.requestWithSudoHandler(a.addSubCollection)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/collection/{collectionID}/subcollection/{subCollectionID}", a.requestWithSudoHandler(a.deleteSubCollection)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/collection/subcollection/image", a.requestWithSudoHandler(a.updateSubCollectionImage)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/collection/subcollection/catalog", a.requestWithSudoHandler(a.addCatalogsToSubCollection)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/collection/subcollection/catalog", a.requestWithSudoHandler(a.removeCatalogsFromSubCollection)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/collection/subcollection/featured", a.requestWithSudoHandler(a.setFeaturedCatalogs)).Methods("POST")

	//KEEPER DISCOUNT
	a.Router.APIRoot.Handle("/keeper/discount", a.requestWithSudoHandler(a.createDiscount)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/discount/{discountID}/deactivate", a.requestWithSudoHandler(a.deactivateDiscount)).Methods("POST")

	//KEEPER COMMISSION
	a.Router.APIRoot.Handle("/keeper/catalogs/brand", a.requestWithSudoHandler(a.getCatalogsByBrandID)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/catalogs/commission", a.requestWithSudoHandler(a.bulkUpdateCommission)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalogs/brand/commission", a.requestWithSudoHandler(a.addCommissionRateBasedonBrandID)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalogs/brand/commission", a.requestWithSudoHandler(a.getCommissionRateUsingBrandID)).Methods("GET")

	//KEEPER SALE
	a.Router.APIRoot.Handle("/keeper/sale", a.requestWithSudoHandler(a.getSales)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/sale/{saleID}/discount", a.requestWithSudoHandler(a.getDiscountInfoBySaleID)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/sale/create", a.requestWithSudoHandler(a.createSale)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/sale/edit", a.requestWithSudoHandler(a.editSale)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/sale/status", a.requestWithSudoHandler(a.editSaleStatus)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/sale/discount", a.requestWithSudoHandler(a.removeDiscountFromSale)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/sale/time", a.requestWithSudoHandler(a.changeSaleTime)).Methods("PUT")

	//APP CATALOG
	a.Router.APIRoot.Handle("/app/groups/catalog", a.requestHandler(a.getGroupsByCatalogID)).Methods("GET")
	a.Router.APIRoot.Handle("/app/brand/catalog", a.requestHandler(a.getCatalogInfoByBrandId)).Methods("GET")

	//APP CATEGORY
	a.Router.APIRoot.Handle("/app/category/lvl1", a.requestHandler(a.getParentCategory)).Methods("GET")
	a.Router.APIRoot.Handle("/app/category/{categoryID}/lvl2", a.requestHandler(a.getMainCategoryByParentID)).Methods("GET")
	a.Router.APIRoot.Handle("/app/category/{categoryID}/lvl3", a.requestHandler(a.getSubCatergoryByParentID)).Methods("GET")

	//APP COLLECTION
	a.Router.APIRoot.Handle("/app/collections", a.requestHandler(a.getActiveCollections)).Methods("GET")
	a.Router.APIRoot.Handle("/app/catalog/basic", a.requestHandler(a.getCatalogBasicByIds)).Methods("GET")
	a.Router.APIRoot.Handle("/app/catalog/similar", a.requestHandler(a.getSimilarProducts)).Methods("GET")
	a.Router.APIRoot.Handle("/app/catalog/{catalogID}", a.requestHandler(a.getCatalogInfoById)).Methods("GET")
	a.Router.APIRoot.Handle("/app/catalog/category/{categoryID}", a.requestHandler(a.getCatalogByCategoryID)).Methods("GET")
	a.Router.APIRoot.Handle("/v2/app/catalog/basic", a.requestHandler(a.getCollectionCatalogByIDs)).Methods("GET")
	a.Router.APIRoot.Handle("/app/subcollection", a.requestHandler(a.GetCatalogBySubCollectionID)).Methods("GET")

	//APP SALE
	a.Router.APIRoot.Handle("/app/sale", a.requestHandler(a.getAppActiveSale)).Methods("GET")
	a.Router.APIRoot.Handle("/app/sale/items", a.requestHandler(a.getSaleCatalogs)).Methods("GET")

	//APP REVIEW
	a.Router.APIRoot.Handle("/app/review", a.requestWithAuthHandler(a.createReview)).Methods("POST")
	a.Router.APIRoot.Handle("/app/review/catalog", a.requestWithAuthHandler(a.getReviewsByCatalogID)).Methods("GET")

	//SEARCH
	// legacy search api
	a.Router.APIRoot.Handle("/app/search", a.requestHandler(a.search)).Methods("GET")

	a.Router.APIRoot.Handle("/app/search/shop", a.requestHandler(a.searchShop)).Methods("GET")
	a.Router.APIRoot.Handle("/app/search/discover", a.requestHandler(a.searchDiscover)).Methods("GET")

	a.Router.APIRoot.Handle("/app/search/catalog", a.requestHandler(a.searchCatalog)).Methods("GET")
	a.Router.APIRoot.Handle("/app/search/brand", a.requestHandler(a.searchBrand)).Methods("GET")
	a.Router.APIRoot.Handle("/app/search/influencer", a.requestHandler(a.searchInfluencer)).Methods("GET")
	a.Router.APIRoot.Handle("/app/search/series", a.requestHandler(a.searchSeries)).Methods("GET")
	a.Router.APIRoot.Handle("/app/search/hashtag", a.requestHandler(a.searchHashtag)).Methods("GET")

	//BULK UPLOAD
	// a.Router.APIRoot.Handle("/keeper/bulk/catalogs/insert", a.requestHandler(a.bulkAddCatalogCSV)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/bulk/catalogs/insert", a.requestHandler(a.bulkAddCatalogJSON)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/bulk/catalogs/price", a.requestWithSudoHandler(a.bulkUpdatePrice)).Methods("PUT")

	//COMMISSION

	//Influencer Collection - KEEPER
	a.Router.APIRoot.Handle("/keeper/influencer/collection", a.requestWithSudoHandler(a.createInfluencerCollection)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/influencer/collections", a.requestWithSudoHandler(a.keeperGetInfluencerCollection)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/influencer/collection", a.requestWithSudoHandler(a.editInfluencerCollection)).Methods("PUT")

	//Influencer Collection - APP
	a.Router.APIRoot.Handle("/app/influencer/collection", a.requestWithAuthHandler(a.createInfluencerCollectionApp)).Methods("POST")
	a.Router.APIRoot.Handle("/app/influencer/collection", a.requestWithAuthHandler(a.editInfluencerCollectionApp)).Methods("PUT")
	a.Router.APIRoot.Handle("/app/influencer/collections/active", a.requestHandler(a.getActiveInfluencerCollections)).Methods("GET")
	a.Router.APIRoot.Handle("/app/influencer/collections", a.requestWithAuthHandler(a.appGetInfluencerCollections)).Methods("GET")
	a.Router.APIRoot.Handle("/app/influencer/collection", a.requestHandler(a.getActiveInfluencerCollectionByID)).Methods("GET")

	//Influencer Collection - KEEPER
	a.Router.APIRoot.Handle("/app/influencer/products", a.requestWithAuthHandler(a.addInfluencerProducts)).Methods("POST")
	a.Router.APIRoot.Handle("/app/influencer/products", a.requestWithAuthHandler(a.removeInfluencerProducts)).Methods("DELETE")
	a.Router.APIRoot.Handle("/app/influencer/products", a.requestWithAuthHandler(a.getInfluencerProducts)).Methods("GET")

}

// InitTestRoutes := intializing all the testing and development endpoints
func (a *API) InitTestRoutes() {
	// a.Router.APIRoot.Handle("/test/add-category", a.requestHandler(a.addSampleCategories)).Methods("GET")
	// a.Router.APIRoot.Handle("/test/add-catalog-and-variants", a.requestHandler(a.imageUpload)).Methods("GET")
}

// InitMediaRoutes initialize media urls
func (a *API) InitMediaRoutes() {
	// fs := http.FileServer(http.Dir("./.media/"))
	// a.Router.Root.PathPrefix("/media/").Handler(http.StripPrefix("/media/", fs))
}
