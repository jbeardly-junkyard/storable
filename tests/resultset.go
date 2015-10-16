package tests

import "gopkg.in/src-d/storable.v1"

type ResultSetFixture struct {
	storable.Document `bson:",inline" collection:"resultset"`
	Foo               string
}

func newResultSetFixture(f string) *ResultSetFixture {
	return &ResultSetFixture{Foo: f}
}
