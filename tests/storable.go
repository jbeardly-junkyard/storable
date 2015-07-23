package example

import (
	"github.com/tyba/storable"
	"github.com/tyba/storable/operators"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (s *AnotherModelStore) New() (doc *AnotherModel) {
	doc = &AnotherModel{}
	doc.SetIsNew(true)
	return
}

func (s *AnotherModelStore) Query() *AnotherModelQuery {
	return &AnotherModelQuery{*storable.NewBaseQuery()}
}

func (s *AnotherModelStore) Find(query *AnotherModelQuery) (*AnotherModelResultSet, error) {
	resultSet, err := s.Store.Find(query)
	if err != nil {
		return nil, err
	}

	return &AnotherModelResultSet{*resultSet}, nil
}

func (s *AnotherModelStore) MustFind(query *AnotherModelQuery) *AnotherModelResultSet {
	resultSet, err := s.Find(query)
	if err != nil {
		panic(err)
	}

	return resultSet
}

func (s *AnotherModelStore) FindOne(query *AnotherModelQuery) (*AnotherModel, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

func (s *AnotherModelStore) MustFindOne(query *AnotherModelQuery) *AnotherModel {
	doc, err := s.FindOne(query)
	if err != nil {
		panic(err)
	}

	return doc
}

func (s *AnotherModelStore) Insert(doc *AnotherModel) error {

	err := s.Store.Insert(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *AnotherModelStore) Update(doc *AnotherModel) error {

	err := s.Store.Update(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *AnotherModelStore) Save(doc *AnotherModel) (updated bool, err error) {

	updated, err = s.Store.Save(doc)
	if err != nil {
		return false, err
	}

	if updated {

	} else {

	}

	return
}

type AnotherModelQuery struct {
	storable.BaseQuery
}

func (q *AnotherModelQuery) FindById(ids ...bson.ObjectId) {
	var vs []interface{}
	for _, id := range ids {
		vs = append(vs, id)
	}
	q.AddCriteria(operators.In(storable.IdField, vs...))
}

type AnotherModelResultSet struct {
	storable.ResultSet
}

func (r *AnotherModelResultSet) All() ([]*AnotherModel, error) {
	var result []*AnotherModel
	err := r.ResultSet.All(&result)

	return result, err
}

func (r *AnotherModelResultSet) One() (*AnotherModel, error) {
	var result *AnotherModel
	_, err := r.ResultSet.One(&result)

	return result, err
}

func (r *AnotherModelResultSet) Next() (*AnotherModel, error) {
	var result *AnotherModel
	_, err := r.ResultSet.Next(&result)

	return result, err
}

type MyModelStore struct {
	storable.Store
}

func NewMyModelStore(db *mgo.Database) *MyModelStore {
	return &MyModelStore{*storable.NewStore(db, "my_model")}
}

func (s *MyModelStore) New() (doc *MyModel) {
	doc = &MyModel{}
	doc.SetIsNew(true)
	return
}

func (s *MyModelStore) Query() *MyModelQuery {
	return &MyModelQuery{*storable.NewBaseQuery()}
}

func (s *MyModelStore) Find(query *MyModelQuery) (*MyModelResultSet, error) {
	resultSet, err := s.Store.Find(query)
	if err != nil {
		return nil, err
	}

	return &MyModelResultSet{*resultSet}, nil
}

func (s *MyModelStore) MustFind(query *MyModelQuery) *MyModelResultSet {
	resultSet, err := s.Find(query)
	if err != nil {
		panic(err)
	}

	return resultSet
}

func (s *MyModelStore) FindOne(query *MyModelQuery) (*MyModel, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

func (s *MyModelStore) MustFindOne(query *MyModelQuery) *MyModel {
	doc, err := s.FindOne(query)
	if err != nil {
		panic(err)
	}

	return doc
}

func (s *MyModelStore) Insert(doc *MyModel) error {

	err := s.Store.Insert(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *MyModelStore) Update(doc *MyModel) error {

	err := s.Store.Update(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *MyModelStore) Save(doc *MyModel) (updated bool, err error) {

	updated, err = s.Store.Save(doc)
	if err != nil {
		return false, err
	}

	if updated {

	} else {

	}

	return
}

type MyModelQuery struct {
	storable.BaseQuery
}

func (q *MyModelQuery) FindById(ids ...bson.ObjectId) {
	var vs []interface{}
	for _, id := range ids {
		vs = append(vs, id)
	}
	q.AddCriteria(operators.In(storable.IdField, vs...))
}

type MyModelResultSet struct {
	storable.ResultSet
}

func (r *MyModelResultSet) All() ([]*MyModel, error) {
	var result []*MyModel
	err := r.ResultSet.All(&result)

	return result, err
}

func (r *MyModelResultSet) One() (*MyModel, error) {
	var result *MyModel
	_, err := r.ResultSet.One(&result)

	return result, err
}

func (r *MyModelResultSet) Next() (*MyModel, error) {
	var result *MyModel
	_, err := r.ResultSet.Next(&result)

	return result, err
}

type schema struct {
	AnotherModel *schemaAnotherModel
	MyModel      *schemaMyModel
}

type schemaAnotherModel struct {
	Foo storable.Field
	Bar storable.Field
}

type schemaMyModel struct {
	String        storable.Field
	Int           storable.Field
	Slice         storable.Field
	SliceAlias    storable.Field
	NestedRef     *schemaMyModelNestedRef
	Nested        *schemaMyModelNested
	NestedSlice   *schemaMyModelNestedSlice
	AliasOfString storable.Field
	Time          storable.Field
	MapsOfString  storable.Map
	InlineStruct  *schemaMyModelInlineStruct
}

type schemaMyModelNestedRef struct {
	X       storable.Field
	Y       storable.Field
	Another *schemaMyModelNestedRefAnother
}

type schemaMyModelNested struct {
	X       storable.Field
	Y       storable.Field
	Another *schemaMyModelNestedRefAnother
}

type schemaMyModelNestedSlice struct {
	X       storable.Field
	Y       storable.Field
	Another *schemaMyModelNestedRefAnother
}

type schemaMyModelInlineStruct struct {
	MapOfString   storable.Map
	MapOfSomeType *schemaMyModelInlineStructMapOfSomeType
}

type schemaMyModelNestedRefAnother struct {
	X storable.Field
	Y storable.Field
}

type schemaMyModelInlineStructMapOfSomeType struct {
	X       storable.Field
	Y       storable.Field
	Another *schemaMyModelNestedRefAnother
}

var Schema = schema{
	AnotherModel: &schemaAnotherModel{
		Foo: storable.NewField("foo", "float64"),
		Bar: storable.NewField("bar", "string"),
	},
	MyModel: &schemaMyModel{
		String:     storable.NewField("string", "string"),
		Int:        storable.NewField("bla2", "int"),
		Slice:      storable.NewField("slice", "string"),
		SliceAlias: storable.NewField("slicealias", "string"),
		NestedRef: &schemaMyModelNestedRef{
			X: storable.NewField("nestedslice.x", "int"),
			Y: storable.NewField("nestedslice.y", "int"),
			Another: &schemaMyModelNestedRefAnother{
				X: storable.NewField("nestedslice.another.x", "int"),
				Y: storable.NewField("nestedslice.another.y", "int"),
			},
		},
		Nested: &schemaMyModelNested{
			X:       storable.NewField("nestedslice.x", "int"),
			Y:       storable.NewField("nestedslice.y", "int"),
			Another: nil,
		},
		NestedSlice: &schemaMyModelNestedSlice{
			X:       storable.NewField("nestedslice.x", "int"),
			Y:       storable.NewField("nestedslice.y", "int"),
			Another: nil,
		},
		AliasOfString: storable.NewField("aliasofstring", "string"),
		Time:          storable.NewField("time", "time.Time"),
		MapsOfString:  storable.NewMap("mapsofstring.[map]", "string"),
		InlineStruct: &schemaMyModelInlineStruct{
			MapOfString: storable.NewMap("inlinestruct.mapofstring.[map]", "string"),
			MapOfSomeType: &schemaMyModelInlineStructMapOfSomeType{
				X:       storable.NewField("nestedslice.x", "int"),
				Y:       storable.NewField("nestedslice.y", "int"),
				Another: nil,
			},
		},
	},
}
