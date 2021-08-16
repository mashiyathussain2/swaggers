package schema

import (
	"go-app/model"

	"syreclabs.com/go/faker"
)

// GetRandomCreateUserOpts returns CreateUserOpts with random data
func GetRandomCreateUserOpts() *CreateUserOpts {
	u := CreateUserOpts{
		Type:  faker.RandomChoice([]string{model.InfluencerType, model.CustomerType, model.BrandType}),
		Email: faker.Internet().FreeEmail(),
		MobileNo: &PhoneNoOpts{
			Prefix: faker.PhoneNumber().AreaCode(),
			Number: faker.PhoneNumber().SubscriberNumber(10),
		},
		Password: faker.Internet().Password(6, 10),
	}
	u.ConfirmPassword = u.Password
	return &u
}

// GetRandomEmailLoginCustomerOpts returns EmailLoginCustomerOpts with random data
func GetRandomEmailLoginCustomerOpts() *EmailLoginCustomerOpts {
	s := EmailLoginCustomerOpts{
		Email:    faker.Internet().FreeEmail(),
		Password: faker.Internet().Password(6, 10),
	}
	return &s
}

// GetRandomCreateBrandOpts returns CreateBrandOpts with random data
func GetRandomCreateBrandOpts() *CreateBrandOpts {
	b := CreateBrandOpts{
		Name:             faker.Company().String(),
		RegisteredName:   faker.Company().Name(),
		FulfillmentEmail: faker.Internet().FreeEmail(),
		FulfillmentCCEmail: []string{
			faker.Internet().FreeEmail(),
			faker.Internet().FreeEmail(),
		},
		Domain:  faker.Internet().DomainName(),
		Website: faker.Internet().Url(),
		Logo: &Img{
			SRC: faker.Avatar().Url("png", 200, 200),
		},
		CoverImg: &Img{
			SRC: faker.Avatar().Url("png", 200, 400),
		},
		Bio: faker.Lorem().Sentence(3),
	}

	if faker.RandomInt(0, 1) == 1 {
		b.SocialAccount = &SocialAccountBrandOpts{
			Facebook: &SocialMediaBrandOpts{
				FollowersCount: faker.RandomInt(0, 10000000),
			},
			Youtube: &SocialMediaBrandOpts{
				FollowersCount: faker.RandomInt(0, 10000000),
			},
			Instagram: &SocialMediaBrandOpts{
				FollowersCount: faker.RandomInt(0, 10000000),
			},
			Twitter: &SocialMediaBrandOpts{
				FollowersCount: faker.RandomInt(0, 10000000),
			},
		}
	}
	return &b
}

// GetRandomCreateInfluencerOpts returns CreateInfluencerOpts with random data
func GetRandomCreateInfluencerOpts() *CreateInfluencerOpts {
	c := CreateInfluencerOpts{
		Name: faker.Company().String(),
		CoverImg: &Img{
			SRC: faker.Avatar().Url("png", 200, 400),
		},
		ProfileImage: &Img{
			SRC: faker.Avatar().Url("png", 200, 400),
		},
		ExternalLinks: []string{faker.Internet().Url()},
		Bio:           faker.Lorem().Sentence(3),
	}

	if faker.RandomInt(0, 1) == 1 {
		c.SocialAccount = &SocialAccountOpts{
			Facebook: &SocialMediaOpts{
				FollowersCount: faker.RandomInt(0, 10000000),
			},
			Youtube: &SocialMediaOpts{
				FollowersCount: faker.RandomInt(0, 10000000),
			},
			Instagram: &SocialMediaOpts{
				FollowersCount: faker.RandomInt(0, 10000000),
			},
			Twitter: &SocialMediaOpts{
				FollowersCount: faker.RandomInt(0, 10000000),
			},
		}
	}
	return &c
}
