package tests

import "github.com/tyba/storable"

type ResultSetFixture struct {
	storable.Document `bson:",inline" collection:"resultset"`
	Foo               string
}

func newResultSetFixture(f string) *ResultSetFixture {
	return &ResultSetFixture{Foo: f}
}
