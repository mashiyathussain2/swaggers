package api

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {
	a.Router.Root.Handle("/", a.requestHandler(a.home)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/category", a.requestHandler(a.createCategory)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/category", a.requestHandler(a.editCategory)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/category", a.requestHandler(a.getCategory)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/category/main", a.requestHandler(a.getMainCategoryMap)).Methods("GET")

	a.Router.APIRoot.Handle("/category/lvl1", a.requestHandler(a.getParentCategory)).Methods("GET")
	a.Router.APIRoot.Handle("/category/{categoryID}/lvl2", a.requestHandler(a.getMainCategoryByParentID)).Methods("GET")
	a.Router.APIRoot.Handle("/category/{categoryID}/lvl3", a.requestHandler(a.getSubCatergoryByParentID)).Methods("GET")
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
