package app

// InitService this initializes all the busines logic services
func InitService(a *App) {
	a.User = InitUser(&UserImplOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.UserConfig.DBName),
		Logger: a.Logger,
	})

	a.Customer = InitCustomer(&CustomerImplOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.CustomerConfig.DBName),
		Logger: a.Logger,
	})
}
