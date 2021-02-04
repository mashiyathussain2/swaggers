package schema

import (
	"fmt"
	"time"

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
		FeaturedImage: img{
			SRC: faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 300, 300),
		},
		Thumbnail: img{
			SRC: faker.Avatar().Url(faker.RandomChoice([]string{"jpg", "jpeg", "png"}), 100, 100),
		},
	}
	return &c
}
