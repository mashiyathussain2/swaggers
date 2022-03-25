package app

// InitService this initializes all the busines logic services
func InitService(a *App) {
	a.KeeperCatalog = InitKeeperCatalog(&KeeperCatalogOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.KeeperCatalogConfig.DBName),
		Logger: a.Logger,
	})
	a.Brand = InitBrand(&BrandOpts{
		App:    a,
		Logger: a.Logger,
	})
	a.Category = InitCategory(&CategoryOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.CategoryConfig.DBName),
		Logger: a.Logger,
	})
	a.Discount = InitDiscount(&DiscountOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.DiscountConfig.DBName),
		Logger: a.Logger,
	})
	a.Group = InitGroup(&GroupOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.GroupConfig.DBName),
		Logger: a.Logger,
	})
	a.Collection = InitCollection(&CollectionOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.CollectionConfig.DBName),
		Logger: a.Logger,
	})
	a.Inventory = InitInventory(&InventoryOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.InventoryConfig.DBName),
		Logger: a.Logger,
	})

	a.Review = InitReview(&ReviewOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.ReviewConfig.DBName),
		Logger: a.Logger,
	})
	a.InfluencerCollection = InitInfluencerCollection(&InfluencerCollectionOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.InfluencerCollectionConfig.DBName),
		Logger: a.Logger,
	})
	a.InfluencerProducts = InitInfluencerProducts(&InfluencerProductsOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.InfluencerProductsConfig.DBName),
		Logger: a.Logger,
	})
}
