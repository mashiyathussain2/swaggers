package model

// GetRandomCreateCatalogResp returns random response based on passed opts
// func GetRandomCreateCatalogResp(opts *schema.CreateCatalogOpts) *Catalog {
// 	res := Catalog{
// 		ID:      primitive.NewObjectIDFromTimestamp(time.Now()),
// 		Name:    opts.Name,
// 		LName:   strings.ToLower(opts.Name),
// 		Slug:    slugify.Slugify(opts.Name),
// 		BrandID: opts.BrandID,
// 		BasePrice: &Price{
// 			Value:       float32(opts.BasePrice),
// 			CurrencyISO: "inr",
// 		},
// 		RetailPrice: &Price{
// 			Value:       float32(opts.RetailPrice),
// 			CurrencyISO: "inr",
// 		},
// 		Description: opts.Description,
// 		Status: Unl,

// 	}
// 	return &res
// }
