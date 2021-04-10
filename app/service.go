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

	a.Brand = InitBrand(&BrandImplOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.BrandConfig.DBName),
		Logger: a.Logger,
	})

	a.Influencer = InitInfluencer(&InfluencerImplOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.InfluencerConfig.DBName),
		Logger: a.Logger,
	})
	a.Cart = InitCart(&CartImplOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.CartConfig.DBName),
		Logger: a.Logger,
	})

	a.KeeperUser = InitKeeperUser(&KeeperUserOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.UserConfig.DBName),
		Logger: a.Logger,
	})
}
