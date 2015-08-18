package tests

import "gopkg.in/tyba/storable.v1"

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
