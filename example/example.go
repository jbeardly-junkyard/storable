package example

import "github.com/maxwellhealth/bongo"

//go:generate mongogen -input=$GOFILE

type MyModel struct {
	bongo.DocumentBase `bson:",inline" collection:"my_model"`

	Bla   string
	Ble   int `bson:"bla2"`
	Bytes []byte
}

type SomeType struct { // not generated
	X int
	Y int
}

type AnotherModel struct {
	bongo.DocumentBase `bson:",inline" collection:"another_model"`
	Foo float64
	Bar string
}

func (m *MyModel) IrrelevantFunction() {
}
