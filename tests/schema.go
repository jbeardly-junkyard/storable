package tests

import "github.com/tyba/storable"

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
