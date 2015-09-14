package tests

import (
	"time"

	"gopkg.in/tyba/storable.v1"
)

type StoreFixture struct {
	storable.Document `bson:",inline" collection:"store"`
	Foo               string
}

type StoreWithConstructFixture struct {
	storable.Document `bson:",inline" collection:"store_construct"`
	Foo               string
}

func newStoreWithConstructFixture(f string) *StoreWithConstructFixture {
	if f == "" {
		return nil
	}
	return &StoreWithConstructFixture{Foo: f}
}

type MultiKeySortFixture struct {
	storable.Document `bson:",inline" collection:"query"`
	Name              string
	Start             time.Time
	End               time.Time
}
