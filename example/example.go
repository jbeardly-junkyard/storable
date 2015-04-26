package example

import "github.com/maxwellhealth/bongo"

//go:generate mongogen gen

type MyModel struct {
	bongo.DocumentBase `bson:",inline" collection:"my_model"`

	Foo       string
	Bar       int `bson:"bla2"`
	Bytes     []byte
	NestedRef *SomeType
	Nested    SomeType
}

type SomeType struct { // not generated
	X       int
	Y       int
	Another AnotherType
}

type AnotherType struct { // not generated
	X int
	Y int
}

type AnotherModel struct {
	bongo.DocumentBase `bson:",inline" collection:"another_model"`
	Foo                float64
	Bar                string
}

func (m *MyModel) IrrelevantFunction() {
}
