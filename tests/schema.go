package tests

import "gopkg.in/tyba/storable.v1"

type SchemaFixture struct {
	storable.Document `bson:",inline" collection:"schema"`

	String         string
	Int            int `bson:"foo"`
	Nested         *SchemaFixture
	MapOfString    map[string]string
	MapOfInterface map[string]interface{}
	MapOfSomeType  map[string]struct {
		Foo string
	}
}
