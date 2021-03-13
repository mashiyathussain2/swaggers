package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"syreclabs.com/go/faker"
)

// GetRandomCustomer returns customer instance with random data
func GetRandomCustomer() *Customer {
	c := Customer{
		UserID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		CartID:   primitive.NewObjectIDFromTimestamp(time.Now()),
		FullName: faker.Name().Name(),
		ProfileImage: &IMG{
			SRC:    faker.Avatar().Url("png", 100, 200),
			Width:  100,
			Height: 200,
		},
		Addresses: []Address{
			{
				ID:          primitive.NewObjectIDFromTimestamp(time.Now()),
				DisplayName: faker.Name().Name(),
				Line1:       faker.Address().BuildingNumber(),
				Line2:       faker.Address().SecondaryAddress(),
				District:    faker.Address().City(),
				State: &State{
					ISOCode: faker.Address().StateAbbr(),
					Name:    faker.Address().State(),
				},
				PostalCode: faker.Address().Postcode(),
				Country: &Country{
					Name:    faker.Address().Country(),
					ISOCode: faker.Address().CountryCode(),
				},
				IsBillingAddress:  true,
				IsShippingAddress: true,
				IsDefaultAddress:  true,
				ContactNumber: &PhoneNumber{
					Prefix: faker.PhoneNumber().ExchangeCode(),
					Number: faker.PhoneNumber().SubscriberNumber(10),
				},
			},
		},
		CreatedAt: time.Now().UTC(),
	}

	dob, _ := time.Parse(time.RFC3339, "1996-10-21T00:00:00+00:00")
	c.DOB = dob
	gender := faker.RandomChoice([]Gender{Male, Female, Others})
	c.Gender = &gender
	return &c
}

// GetRandomCustomerUser returns user instance with random data
func GetRandomCustomerUser() *User {
	u := User{
		Type:  CustomerType,
		Role:  AdminRole,
		Email: faker.Internet().SafeEmail(),
		PhoneNo: &PhoneNumber{
			Prefix: faker.PhoneNumber().ExchangeCode(),
			Number: faker.PhoneNumber().SubscriberNumber(10),
		},
		Username:  faker.Internet().UserName(),
		Password:  faker.Internet().Password(6, 10),
		CreatedAt: time.Now().UTC(),
	}
	if faker.RandomInt(0, 1) == 1 {
		u.EmailVerifiedAt = time.Now().UTC()
	}
	if faker.RandomInt(0, 1) == 1 {
		u.PhoneVerifiedAt = time.Now().UTC()
	}
	return &u
}
