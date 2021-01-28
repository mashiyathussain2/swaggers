package app

// InitService this initializes all the busines logic services
func InitService(a *App) {
	a.KeeperCatalog = InitKeeperCatalog(&KeeperCatalogOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.KeeperCatalogConfig.DBName),
		Logger: a.Logger,
	})
}
