package schema

import (
	"fmt"
	"go-app/model"
	"time"

	"github.com/avelino/slugify"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/icrowley/fake"
	"go.mongodb.org/mongo-driver/bson/primitive"
	faker "syreclabs.com/go/faker"
)

// GetRandomCreateCatalogOpts returns CreateCatalogOpts with random data
func GetRandomCreateCatalogOpts() *CreateCatalogOpts {
	s := CreateCatalogOpts{
		Name:        faker.Commerce().ProductName(),
		Description: faker.Lorem().Paragraph(gofakeit.Number(1, 5)),
		HSNCode:     faker.Code().Ean8(),
		Keywords:    faker.Lorem().Words(gofakeit.Number(1, 5)),
		RetailPrice: uint32(gofakeit.Price(1, 10000)),
	}

	s.BrandID = primitive.NewObjectIDFromTimestamp(time.Now())

	var sp []specsOpts
	var fa []filterAttribute
	var cid []primitive.ObjectID
	for i := 0; i < gofakeit.Number(1, 5); i++ {
		sp = append(sp, specsOpts{Name: fake.Title(), Value: fake.Sentence()})
		cid = append(cid, primitive.NewObjectIDFromTimestamp(time.Now()))
		fa = append(fa, filterAttribute{Name: fake.Title(), Value: fake.Sentence()})
	}

	s.Specifications = sp
	s.CategoryID = cid
	s.BasePrice = uint32(gofakeit.Price(1, 5000)) + s.RetailPrice

	return &s
}

// GetRandomCreateVariantOpts returns CreateVariantOpts with random data
func GetRandomCreateVariantOpts() *CreateVariantOpts {
	s := CreateVariantOpts{
		SKU:       faker.Internet().Slug(),
		Attribute: faker.RandomChoice([]string{"", "color", "size"}),
	}
	return &s
}

// GetRandomCreateBrandOpts returns CreateBrandOpts with random data in it
func GetRandomCreateBrandOpts() *CreateBrandOpts {
	b := CreateBrandOpts{
		Name:             faker.Company().Name(),
		Description:      fake.Sentences(),
		WebsiteLink:      faker.Internet().Url(),
		FulfillmentEmail: faker.Internet().Email(),
	}
	b.RegisteredName = fmt.Sprintf("%s pvt ltd", b.Name)
	return &b
}

// GetRandomCreateCategoryOpts returns CreateCategoryOpts with random data
func GetRandomCreateCategoryOpts() *CreateCategoryOpts {
	c := CreateCategoryOpts{
		Name:     faker.Commerce().Department(),
		ParentID: primitive.NewObjectIDFromTimestamp(time.Now()),
		FeaturedImage: Img{
			SRC: faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 300, 300),
		},
		Thumbnail: Img{
			SRC: faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 100, 100),
		},
	}
	return &c
}

// GetRandomCreateCategoryResp returns random response based on passed opts
func GetRandomCreateCategoryResp(opts *CreateCategoryOpts) *CreateCategoryResp {
	res := CreateCategoryResp{
		ID:         primitive.NewObjectIDFromTimestamp(time.Now()),
		Name:       opts.Name,
		Slug:       slugify.Slugify(opts.Name),
		ParentID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		AncestorID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now()), primitive.NewObjectIDFromTimestamp(time.Now())},
		FeaturedImage: &model.IMG{
			SRC:    faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 300, 300),
			Width:  300,
			Height: 300,
		},
		Thumbnail: &model.IMG{
			SRC:    faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 100, 100),
			Width:  100,
			Height: 100,
		},
		IsMain: false,
	}
	if faker.RandomInt(0, 1) == 1 {
		res.IsMain = true
	}
	return &res
}

// GetRandomEditCategoryResp returns random response based on passed opts
func GetRandomEditCategoryResp(opts *EditCategoryOpts) *EditCategoryResp {
	res := EditCategoryResp{
		ID:         primitive.NewObjectIDFromTimestamp(time.Now()),
		Name:       opts.Name,
		Slug:       slugify.Slugify(opts.Name),
		ParentID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		AncestorID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now()), primitive.NewObjectIDFromTimestamp(time.Now())},
		FeaturedImage: &model.IMG{
			SRC:    faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 300, 300),
			Width:  300,
			Height: 300,
		},
		Thumbnail: &model.IMG{
			SRC:    faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 100, 100),
			Width:  100,
			Height: 100,
		},
		IsMain: false,
	}
	if faker.RandomInt(0, 1) == 1 {
		res.IsMain = true
	}
	return &res
}

// GetRandomEditCategoryOpts masking random edit category with create category
func GetRandomEditCategoryOpts() *EditCategoryOpts {
	t := true
	c := EditCategoryOpts{
		ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		Name: faker.Commerce().Department(),
		FeaturedImage: &Img{
			SRC: faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 300, 300),
		},
		Thumbnail: &Img{
			SRC: faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 100, 100),
		},
		IsMain: &t,
	}
	if faker.RandomInt(0, 1) == 1 {
		f := false
		c.IsMain = &f
	}
	return &c
}

// GetRandomGetCategoriesResp returns random data into GetCategoriesResp struct
func GetRandomGetCategoriesResp() *GetCategoriesResp {
	c := GetCategoriesResp{
		ID:         primitive.NewObjectIDFromTimestamp(time.Now()),
		ParentID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		AncestorID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now()), primitive.NewObjectIDFromTimestamp(time.Now()), primitive.NewObjectIDFromTimestamp(time.Now())},
		Thumbnail: &model.IMG{
			SRC:    faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 300, 300),
			Width:  300,
			Height: 300,
		},
		FeaturedImage: &model.IMG{
			SRC:    faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 100, 100),
			Width:  100,
			Height: 100,
		},
		IsMain: true,
	}
	if faker.RandomInt(0, 1) == 1 {
		c.IsMain = true
	}
	return &c
}

// GetRandomGetMainCategoriesMapResp fills random data into struct
func GetRandomGetMainCategoriesMapResp() *GetMainCategoriesMapResp {
	c := GetMainCategoriesMapResp{
		ID:         primitive.NewObjectIDFromTimestamp(time.Now()),
		ParentID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		AncestorID: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now()), primitive.NewObjectIDFromTimestamp(time.Now()), primitive.NewObjectIDFromTimestamp(time.Now())},
		Thumbnail: &model.IMG{
			SRC:    faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 100, 100),
			Width:  100,
			Height: 100,
		},
		FeaturedImage: &model.IMG{
			SRC:    faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 300, 300),
			Width:  300,
			Height: 300,
		},
	}
	return &c
}

// GetRandomGetParentCategoriesResp fills random data into struct
func GetRandomGetParentCategoriesResp() *GetParentCategoriesResp {
	c := GetParentCategoriesResp{
		ID: primitive.NewObjectIDFromTimestamp(time.Now()),
		Thumbnail: &model.IMG{
			SRC:    faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 100, 100),
			Width:  100,
			Height: 100,
		},
	}
	return &c
}

// GetRandomGetMainCategoriesByParentIDResp fills random data into struct
func GetRandomGetMainCategoriesByParentIDResp() *GetMainCategoriesByParentIDResp {
	c := GetMainCategoriesByParentIDResp{
		ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		Name: faker.Commerce().ProductName(),
		FeaturedImage: &model.IMG{
			SRC:    faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 300, 300),
			Width:  300,
			Height: 300,
		},
	}
	return &c
}

// GetRandomGetSubCategoriesByParentIDResp fills random data into struct
func GetRandomGetSubCategoriesByParentIDResp() *GetSubCategoriesByParentIDResp {
	c := GetSubCategoriesByParentIDResp{
		ID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		Name: faker.Commerce().ProductName(),
	}
	return &c
}
