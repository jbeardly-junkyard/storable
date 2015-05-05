package example

import (
	"time"

	"github.com/tyba/storable"

	"gopkg.in/mgo.v2"
)

//go:generate storable gen

type Alias string
type MyModel struct {
	storable.Document `bson:",inline" collection:"my_model"`

	String        string
	Int           int `bson:"bla2"`
	Bytes         []byte
	Slice         []string
	NestedRef     *SomeType
	Nested        SomeType
	NestedSlice   []*SomeType
	AliasOfString Alias
	Time          time.Time
	MapsOfString  map[string]string
	InlineStruct  struct {
		MapOfString   map[string]string
		MapOfSomeType map[string]SomeType
	}
}

func (m *MyModel) IrrelevantFunction() {
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
	storable.Document `bson:",inline" collection:"another_model"`
	Foo               float64
	Bar               string
}

type AnotherModelStore struct {
	storable.Store
	Foo bool
}

func NewAnotherModelStore(db *mgo.Database, foo bool) *AnotherModelStore {
	return &AnotherModelStore{
		*storable.NewStore(db, "another_model"), foo,
	}
}

func NewAnotherModel() *AnotherModel {
	return &AnotherModel{}
}
