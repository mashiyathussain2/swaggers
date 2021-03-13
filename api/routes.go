package api

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {
	a.Router.Root.Handle("/", a.requestHandler(a.home)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/live/create", a.requestHandler(a.createLiveStream)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/live/start", a.requestHandler(a.home)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/live/stop", a.requestHandler(a.home)).Methods("POST")
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
