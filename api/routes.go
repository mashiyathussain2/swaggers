package api

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {

	//KEEPER CATEGORY
	a.Router.Root.Handle("/", a.requestHandler(a.home)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/category", a.requestHandler(a.createCategory)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/category", a.requestHandler(a.editCategory)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/category", a.requestHandler(a.getCategory)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/category/main", a.requestHandler(a.getMainCategoryMap)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/category/{categoryID}/path", a.requestHandler(a.getCategoryPath)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/category/{categoryID}/ancestors", a.requestHandler(a.getAncestorsByID)).Methods("GET")

	//KEEPER CATALOG
	a.Router.APIRoot.Handle("/keeper/catalog", a.requestHandler(a.createCatalog)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog", a.requestHandler(a.editCatalog)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/catalog/variant", a.requestHandler(a.addVariants)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog/status", a.requestHandler(a.updateCatalogStatus)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog/variant", a.requestHandler(a.deleteVariant)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/catalog/search", a.requestHandler(a.keeperSearchCatalog)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/catalog/content/video", a.requestHandler(a.addCatalogContent)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog/content/image", a.requestHandler(a.addCatalogContentImage)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/catalog", a.requestHandler(a.getCatalogsByFilter)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/catalog/{slug}", a.requestHandler(a.getCatalogBySlug)).Methods("GET")

	//KEEPER GROUP
	a.Router.APIRoot.Handle("/keeper/group", a.requestHandler(a.createCatalogGroup)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/group", a.requestHandler(a.getCatalogsByGroupID)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/groups", a.requestHandler(a.getGroups)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/groups/catalog", a.requestHandler(a.keeperGetGroupsByCatalogID)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/group/catalogs", a.requestHandler(a.addCatalogsInTheGroup)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/group", a.requestHandler(a.editGroup)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/groups/catalog-name", a.requestHandler(a.getGroupsByCatalogName)).Methods("GET")

	//KEEPER INVENTORY
	a.Router.APIRoot.Handle("/keeper/inventory", a.requestHandler(a.updateInventory)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/inventory/outofstock", a.requestHandler(a.setOutOfStock)).Methods("POST")

	//KEEPER COLLECTION
	a.Router.APIRoot.Handle("/keeper/collection", a.requestHandler(a.createCollection)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/collection", a.requestHandler(a.editCollection)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/collection/{collectionID}", a.requestHandler(a.deleteCollection)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/collection/subcollection", a.requestHandler(a.addSubCollection)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/collection/{collectionID}/subcollection/{subCollectionID}", a.requestHandler(a.deleteSubCollection)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/collection/subcollection/image", a.requestHandler(a.updateSubCollectionImage)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/collection/subcollection/catalog", a.requestHandler(a.addCatalogsToSubCollection)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/collection/subcollection/catalog", a.requestHandler(a.removeCatalogsFromSubCollection)).Methods("DELETE")

	//KEEPER DISCOUNT
	a.Router.APIRoot.Handle("/keeper/discount", a.requestHandler(a.createDiscount)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/discount/{discountID}/deactivate", a.requestHandler(a.deactivateDiscount)).Methods("POST")

	//KEEPER SALE
	a.Router.APIRoot.Handle("/keeper/sale", a.requestHandler(a.createSale)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/sale", a.requestHandler(a.editSale)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/sale/status", a.requestHandler(a.editSaleStatus)).Methods("POST")

	//APP CATALOG
	a.Router.APIRoot.Handle("/app/groups/catalog", a.requestHandler(a.getGroupsByCatalogID)).Methods("GET")

	//APP CATEGORY
	a.Router.APIRoot.Handle("/app/category/lvl1", a.requestHandler(a.getParentCategory)).Methods("GET")
	a.Router.APIRoot.Handle("/app/category/{categoryID}/lvl2", a.requestHandler(a.getMainCategoryByParentID)).Methods("GET")
	a.Router.APIRoot.Handle("/app/category/{categoryID}/lvl3", a.requestHandler(a.getSubCatergoryByParentID)).Methods("GET")
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
