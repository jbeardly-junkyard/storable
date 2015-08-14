package tests

import "github.com/tyba/storable"

type QueryFixture struct {
	storable.Document `bson:",inline" collection:"query"`
	Foo               string
}

func newQueryFixture(f string) *QueryFixture {
	return &QueryFixture{Foo: f}
}
