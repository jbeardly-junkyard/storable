package tests

import "github.com/tyba/storable"

type StoreFixture struct {
	storable.Document `bson:",inline" collection:"store"`
	Foo               string
}

type StoreWithConstructFixture struct {
	storable.Document `bson:",inline" collection:"store_construct"`
	Foo               string
}

func newStoreWithConstructFixture(f string) *StoreWithConstructFixture {
	return &StoreWithConstructFixture{Foo: f}
}
