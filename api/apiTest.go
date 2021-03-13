package api

import (
	"go-app/app"
	"go-app/server/config"
	"go-app/server/logger"
	"go-app/server/validator"

	"github.com/gorilla/mux"
)

// NewTestAPI returns api struct for unit testing
func NewTestAPI(c *config.APIConfig) *API {
	l := logger.NewLogger(nil, logger.NewZeroLogConsoleWriter(logger.NewStandardConsoleWriter()), nil)
	api := &API{
		MainRouter: &mux.Router{},
		Router:     &Router{},
		Config:     c,
		Logger:     l,
		Validator:  validator.NewValidation(),
	}
	api.setupRoutes()
	api.App = &app.App{}
	return api
}

// func getTestConfig() *config.APIConfig {
// 	c := config.GetConfigFromFile("test")
// 	return &c.APIConfig
// }
