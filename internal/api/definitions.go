package api

import "time"

type Admin struct {
	Id        int       `jsonapi:"id"`
	CreatedAt time.Time `jsonapi:"created_at"`
	UpdatedAt time.Time `jsonapi:"updated_at"`
	Name      string    `jsonapi:"name"`
	Secret    string    `jsonapi:"api_secret"`
}

type Merchant struct {
	Id        int       `jsonapi:"primary,merchants"`
	CreatedAt time.Time `jsonapi:"attr,created_at"`
	UpdatedAt time.Time `jsonapi:"attr,updated_at"`
	Name      string    `jsonapi:"attr,name"`
	Secret    string    `jsonapi:"attr,api_secret"`
}

type Member struct {
	Id         int       `jsonapi:"primary,members"`
	CreatedAt  time.Time `jsonapi:"attr,created_at"`
	UpdatedAt  time.Time `jsonapi:"attr,updated_at"`
	Name       string    `jsonapi:"attr,name"`
	Secret     string    `jsonapi:"attr,api_secret"`
	Email      string    `jsonapi:"attr,email"`
	MerchantId int       `jsonapi:"attr,merchant_id"`
}

type MerchantRequest struct {
	Data struct {
		Type       string   `jsonapi: "type"`
		Attributes Merchant `jsonapi: "attributes"`
	} `jsonapi: "data"`
}
