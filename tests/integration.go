package example

import (
	"time"

	"github.com/tyba/storable"

	"gopkg.in/mgo.v2"
)

//go:generate storable gen

type Alias string
type SliceAlias []string
type MyModel struct {
	storable.Document `bson:",inline" collection:"my_model"`

	String        string
	Int           int `bson:"bla2"`
	Bytes         []byte
	Slice         []string
	SliceAlias    SliceAlias
	NestedRef     *SomeType
	Nested        SomeType
	NestedSlice   []*SomeType
	AliasOfString Alias
	Time          time.Time
	MapsOfString  map[string]string
	InlineStruct  struct {
		MapOfString    map[string]string
		MapOfSomeType  map[string]SomeType
		MapOfInterface map[string]interface{}
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

type EventsTests struct {
	storable.Document `bson:",inline" collection:"event"`
	Checks            map[string]bool
	MustFailBefore    error
	MustFailAfter     error
}

func newEventsTests() *EventsTests {
	return &EventsTests{
		Checks: make(map[string]bool, 0),
	}
}

func (s *EventsTestsStore) BeforeInsert(doc *EventsTests) error {
	if doc.MustFailBefore != nil {
		return doc.MustFailBefore
	}

	doc.Checks["BeforeInsert"] = true
	return nil
}

func (s *EventsTestsStore) AfterInsert(doc *EventsTests) error {
	if doc.MustFailAfter != nil {
		return doc.MustFailAfter
	}

	doc.Checks["AfterInsert"] = true
	return nil
}

func (s *EventsTestsStore) BeforeUpdate(doc *EventsTests) error {
	if doc.MustFailBefore != nil {
		return doc.MustFailBefore
	}

	doc.Checks["BeforeUpdate"] = true
	return nil
}

func (s *EventsTestsStore) AfterUpdate(doc *EventsTests) error {
	if doc.MustFailAfter != nil {
		return doc.MustFailAfter
	}

	doc.Checks["AfterUpdate"] = true
	return nil
}
