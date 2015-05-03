package example

import (
	"time"

	"github.com/tyba/storable"
)

//go:generate storable gen

type Status int

const (
	Draft Status = iota
	Published
)

type Product struct {
	storable.Document `bson:",inline" collection:"products"`

	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Price     Price
	Discount  float64
	Url       string
	Tags      []string
}

type Price struct {
	Amount   float64
	Discount float64
}
