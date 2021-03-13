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
