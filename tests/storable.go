package tests

import (
	"github.com/tyba/storable"
	"github.com/tyba/storable/operators"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type EventsFixtureStore struct {
	storable.Store
}

func NewEventsFixtureStore(db *mgo.Database) *EventsFixtureStore {
	return &EventsFixtureStore{*storable.NewStore(db, "event")}
}

func (s *EventsFixtureStore) New() (doc *EventsFixture) {
	doc = newEventsFixture()
	doc.SetIsNew(true)
	doc.SetId(bson.NewObjectId())
	return
}

func (s *EventsFixtureStore) Query() *EventsFixtureQuery {
	return &EventsFixtureQuery{*storable.NewBaseQuery()}
}

func (s *EventsFixtureStore) Find(query *EventsFixtureQuery) (*EventsFixtureResultSet, error) {
	resultSet, err := s.Store.Find(query)
	if err != nil {
		return nil, err
	}

	return &EventsFixtureResultSet{*resultSet}, nil
}

func (s *EventsFixtureStore) MustFind(query *EventsFixtureQuery) *EventsFixtureResultSet {
	resultSet, err := s.Find(query)
	if err != nil {
		panic(err)
	}

	return resultSet
}

func (s *EventsFixtureStore) FindOne(query *EventsFixtureQuery) (*EventsFixture, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

func (s *EventsFixtureStore) MustFindOne(query *EventsFixtureQuery) *EventsFixture {
	doc, err := s.FindOne(query)
	if err != nil {
		panic(err)
	}

	return doc
}

func (s *EventsFixtureStore) Insert(doc *EventsFixture) error {
	if err := s.BeforeInsert(doc); err != nil {
		return err
	}

	err := s.Store.Insert(doc)
	if err != nil {
		return err
	}

	return s.AfterInsert(doc)
}

func (s *EventsFixtureStore) Update(doc *EventsFixture) error {
	if err := s.BeforeUpdate(doc); err != nil {
		return err
	}

	err := s.Store.Update(doc)
	if err != nil {
		return err
	}

	return s.AfterUpdate(doc)
}

func (s *EventsFixtureStore) Save(doc *EventsFixture) (updated bool, err error) {
	switch doc.IsNew() {
	case true:
		if err := s.BeforeInsert(doc); err != nil {
			return false, err
		}
	case false:
		if err := s.BeforeUpdate(doc); err != nil {
			return false, err
		}
	}

	updated, err = s.Store.Save(doc)
	if err != nil {
		return false, err
	}

	switch updated {
	case false:
		if err := s.AfterInsert(doc); err != nil {
			return false, err
		}
	case true:
		if err := s.AfterUpdate(doc); err != nil {
			return false, err
		}
	}

	return
}

type EventsFixtureQuery struct {
	storable.BaseQuery
}

func (q *EventsFixtureQuery) FindById(ids ...bson.ObjectId) *EventsFixtureQuery {
	var vs []interface{}
	for _, id := range ids {
		vs = append(vs, id)
	}
	q.AddCriteria(operators.In(storable.IdField, vs...))

	return q
}

type EventsFixtureResultSet struct {
	storable.ResultSet
}

func (r *EventsFixtureResultSet) All() ([]*EventsFixture, error) {
	var result []*EventsFixture
	err := r.ResultSet.All(&result)

	return result, err
}

func (r *EventsFixtureResultSet) One() (*EventsFixture, error) {
	var result *EventsFixture
	err := r.ResultSet.One(&result)

	return result, err
}

func (r *EventsFixtureResultSet) Next() (*EventsFixture, error) {
	var result *EventsFixture
	_, err := r.ResultSet.Next(&result)

	return result, err
}

func (r *EventsFixtureResultSet) ForEach(f func(*EventsFixture) error) error {
	for {
		var result *EventsFixture
		found, err := r.ResultSet.Next(&result)
		if err != nil {
			return err
		}

		if !found {
			break
		}

		err = f(result)
		if err == storable.ErrStop {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

type QueryFixtureStore struct {
	storable.Store
}

func NewQueryFixtureStore(db *mgo.Database) *QueryFixtureStore {
	return &QueryFixtureStore{*storable.NewStore(db, "query")}
}

func (s *QueryFixtureStore) New(f string) (doc *QueryFixture) {
	doc = newQueryFixture(f)
	doc.SetIsNew(true)
	doc.SetId(bson.NewObjectId())
	return
}

func (s *QueryFixtureStore) Query() *QueryFixtureQuery {
	return &QueryFixtureQuery{*storable.NewBaseQuery()}
}

func (s *QueryFixtureStore) Find(query *QueryFixtureQuery) (*QueryFixtureResultSet, error) {
	resultSet, err := s.Store.Find(query)
	if err != nil {
		return nil, err
	}

	return &QueryFixtureResultSet{*resultSet}, nil
}

func (s *QueryFixtureStore) MustFind(query *QueryFixtureQuery) *QueryFixtureResultSet {
	resultSet, err := s.Find(query)
	if err != nil {
		panic(err)
	}

	return resultSet
}

func (s *QueryFixtureStore) FindOne(query *QueryFixtureQuery) (*QueryFixture, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

func (s *QueryFixtureStore) MustFindOne(query *QueryFixtureQuery) *QueryFixture {
	doc, err := s.FindOne(query)
	if err != nil {
		panic(err)
	}

	return doc
}

func (s *QueryFixtureStore) Insert(doc *QueryFixture) error {

	err := s.Store.Insert(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *QueryFixtureStore) Update(doc *QueryFixture) error {

	err := s.Store.Update(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *QueryFixtureStore) Save(doc *QueryFixture) (updated bool, err error) {
	updated, err = s.Store.Save(doc)
	if err != nil {
		return false, err
	}

	return
}

type QueryFixtureQuery struct {
	storable.BaseQuery
}

func (q *QueryFixtureQuery) FindById(ids ...bson.ObjectId) *QueryFixtureQuery {
	var vs []interface{}
	for _, id := range ids {
		vs = append(vs, id)
	}
	q.AddCriteria(operators.In(storable.IdField, vs...))

	return q
}

type QueryFixtureResultSet struct {
	storable.ResultSet
}

func (r *QueryFixtureResultSet) All() ([]*QueryFixture, error) {
	var result []*QueryFixture
	err := r.ResultSet.All(&result)

	return result, err
}

func (r *QueryFixtureResultSet) One() (*QueryFixture, error) {
	var result *QueryFixture
	err := r.ResultSet.One(&result)

	return result, err
}

func (r *QueryFixtureResultSet) Next() (*QueryFixture, error) {
	var result *QueryFixture
	_, err := r.ResultSet.Next(&result)

	return result, err
}

func (r *QueryFixtureResultSet) ForEach(f func(*QueryFixture) error) error {
	for {
		var result *QueryFixture
		found, err := r.ResultSet.Next(&result)
		if err != nil {
			return err
		}

		if !found {
			break
		}

		err = f(result)
		if err == storable.ErrStop {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

type ResultSetFixtureStore struct {
	storable.Store
}

func NewResultSetFixtureStore(db *mgo.Database) *ResultSetFixtureStore {
	return &ResultSetFixtureStore{*storable.NewStore(db, "resultset")}
}

func (s *ResultSetFixtureStore) New(f string) (doc *ResultSetFixture) {
	doc = newResultSetFixture(f)
	doc.SetIsNew(true)
	doc.SetId(bson.NewObjectId())
	return
}

func (s *ResultSetFixtureStore) Query() *ResultSetFixtureQuery {
	return &ResultSetFixtureQuery{*storable.NewBaseQuery()}
}

func (s *ResultSetFixtureStore) Find(query *ResultSetFixtureQuery) (*ResultSetFixtureResultSet, error) {
	resultSet, err := s.Store.Find(query)
	if err != nil {
		return nil, err
	}

	return &ResultSetFixtureResultSet{*resultSet}, nil
}

func (s *ResultSetFixtureStore) MustFind(query *ResultSetFixtureQuery) *ResultSetFixtureResultSet {
	resultSet, err := s.Find(query)
	if err != nil {
		panic(err)
	}

	return resultSet
}

func (s *ResultSetFixtureStore) FindOne(query *ResultSetFixtureQuery) (*ResultSetFixture, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

func (s *ResultSetFixtureStore) MustFindOne(query *ResultSetFixtureQuery) *ResultSetFixture {
	doc, err := s.FindOne(query)
	if err != nil {
		panic(err)
	}

	return doc
}

func (s *ResultSetFixtureStore) Insert(doc *ResultSetFixture) error {

	err := s.Store.Insert(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *ResultSetFixtureStore) Update(doc *ResultSetFixture) error {

	err := s.Store.Update(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *ResultSetFixtureStore) Save(doc *ResultSetFixture) (updated bool, err error) {
	updated, err = s.Store.Save(doc)
	if err != nil {
		return false, err
	}

	return
}

type ResultSetFixtureQuery struct {
	storable.BaseQuery
}

func (q *ResultSetFixtureQuery) FindById(ids ...bson.ObjectId) *ResultSetFixtureQuery {
	var vs []interface{}
	for _, id := range ids {
		vs = append(vs, id)
	}
	q.AddCriteria(operators.In(storable.IdField, vs...))

	return q
}

type ResultSetFixtureResultSet struct {
	storable.ResultSet
}

func (r *ResultSetFixtureResultSet) All() ([]*ResultSetFixture, error) {
	var result []*ResultSetFixture
	err := r.ResultSet.All(&result)

	return result, err
}

func (r *ResultSetFixtureResultSet) One() (*ResultSetFixture, error) {
	var result *ResultSetFixture
	err := r.ResultSet.One(&result)

	return result, err
}

func (r *ResultSetFixtureResultSet) Next() (*ResultSetFixture, error) {
	var result *ResultSetFixture
	_, err := r.ResultSet.Next(&result)

	return result, err
}

func (r *ResultSetFixtureResultSet) ForEach(f func(*ResultSetFixture) error) error {
	for {
		var result *ResultSetFixture
		found, err := r.ResultSet.Next(&result)
		if err != nil {
			return err
		}

		if !found {
			break
		}

		err = f(result)
		if err == storable.ErrStop {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

type SchemaFixtureStore struct {
	storable.Store
}

func NewSchemaFixtureStore(db *mgo.Database) *SchemaFixtureStore {
	return &SchemaFixtureStore{*storable.NewStore(db, "schema")}
}

func (s *SchemaFixtureStore) New() (doc *SchemaFixture) {
	doc = &SchemaFixture{}
	doc.SetIsNew(true)
	doc.SetId(bson.NewObjectId())
	return
}

func (s *SchemaFixtureStore) Query() *SchemaFixtureQuery {
	return &SchemaFixtureQuery{*storable.NewBaseQuery()}
}

func (s *SchemaFixtureStore) Find(query *SchemaFixtureQuery) (*SchemaFixtureResultSet, error) {
	resultSet, err := s.Store.Find(query)
	if err != nil {
		return nil, err
	}

	return &SchemaFixtureResultSet{*resultSet}, nil
}

func (s *SchemaFixtureStore) MustFind(query *SchemaFixtureQuery) *SchemaFixtureResultSet {
	resultSet, err := s.Find(query)
	if err != nil {
		panic(err)
	}

	return resultSet
}

func (s *SchemaFixtureStore) FindOne(query *SchemaFixtureQuery) (*SchemaFixture, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

func (s *SchemaFixtureStore) MustFindOne(query *SchemaFixtureQuery) *SchemaFixture {
	doc, err := s.FindOne(query)
	if err != nil {
		panic(err)
	}

	return doc
}

func (s *SchemaFixtureStore) Insert(doc *SchemaFixture) error {

	err := s.Store.Insert(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *SchemaFixtureStore) Update(doc *SchemaFixture) error {

	err := s.Store.Update(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *SchemaFixtureStore) Save(doc *SchemaFixture) (updated bool, err error) {
	updated, err = s.Store.Save(doc)
	if err != nil {
		return false, err
	}

	return
}

type SchemaFixtureQuery struct {
	storable.BaseQuery
}

func (q *SchemaFixtureQuery) FindById(ids ...bson.ObjectId) *SchemaFixtureQuery {
	var vs []interface{}
	for _, id := range ids {
		vs = append(vs, id)
	}
	q.AddCriteria(operators.In(storable.IdField, vs...))

	return q
}

type SchemaFixtureResultSet struct {
	storable.ResultSet
}

func (r *SchemaFixtureResultSet) All() ([]*SchemaFixture, error) {
	var result []*SchemaFixture
	err := r.ResultSet.All(&result)

	return result, err
}

func (r *SchemaFixtureResultSet) One() (*SchemaFixture, error) {
	var result *SchemaFixture
	err := r.ResultSet.One(&result)

	return result, err
}

func (r *SchemaFixtureResultSet) Next() (*SchemaFixture, error) {
	var result *SchemaFixture
	_, err := r.ResultSet.Next(&result)

	return result, err
}

func (r *SchemaFixtureResultSet) ForEach(f func(*SchemaFixture) error) error {
	for {
		var result *SchemaFixture
		found, err := r.ResultSet.Next(&result)
		if err != nil {
			return err
		}

		if !found {
			break
		}

		err = f(result)
		if err == storable.ErrStop {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

type StoreFixtureStore struct {
	storable.Store
}

func NewStoreFixtureStore(db *mgo.Database) *StoreFixtureStore {
	return &StoreFixtureStore{*storable.NewStore(db, "store")}
}

func (s *StoreFixtureStore) New() (doc *StoreFixture) {
	doc = &StoreFixture{}
	doc.SetIsNew(true)
	doc.SetId(bson.NewObjectId())
	return
}

func (s *StoreFixtureStore) Query() *StoreFixtureQuery {
	return &StoreFixtureQuery{*storable.NewBaseQuery()}
}

func (s *StoreFixtureStore) Find(query *StoreFixtureQuery) (*StoreFixtureResultSet, error) {
	resultSet, err := s.Store.Find(query)
	if err != nil {
		return nil, err
	}

	return &StoreFixtureResultSet{*resultSet}, nil
}

func (s *StoreFixtureStore) MustFind(query *StoreFixtureQuery) *StoreFixtureResultSet {
	resultSet, err := s.Find(query)
	if err != nil {
		panic(err)
	}

	return resultSet
}

func (s *StoreFixtureStore) FindOne(query *StoreFixtureQuery) (*StoreFixture, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

func (s *StoreFixtureStore) MustFindOne(query *StoreFixtureQuery) *StoreFixture {
	doc, err := s.FindOne(query)
	if err != nil {
		panic(err)
	}

	return doc
}

func (s *StoreFixtureStore) Insert(doc *StoreFixture) error {

	err := s.Store.Insert(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreFixtureStore) Update(doc *StoreFixture) error {

	err := s.Store.Update(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreFixtureStore) Save(doc *StoreFixture) (updated bool, err error) {
	updated, err = s.Store.Save(doc)
	if err != nil {
		return false, err
	}

	return
}

type StoreFixtureQuery struct {
	storable.BaseQuery
}

func (q *StoreFixtureQuery) FindById(ids ...bson.ObjectId) *StoreFixtureQuery {
	var vs []interface{}
	for _, id := range ids {
		vs = append(vs, id)
	}
	q.AddCriteria(operators.In(storable.IdField, vs...))

	return q
}

type StoreFixtureResultSet struct {
	storable.ResultSet
}

func (r *StoreFixtureResultSet) All() ([]*StoreFixture, error) {
	var result []*StoreFixture
	err := r.ResultSet.All(&result)

	return result, err
}

func (r *StoreFixtureResultSet) One() (*StoreFixture, error) {
	var result *StoreFixture
	err := r.ResultSet.One(&result)

	return result, err
}

func (r *StoreFixtureResultSet) Next() (*StoreFixture, error) {
	var result *StoreFixture
	_, err := r.ResultSet.Next(&result)

	return result, err
}

func (r *StoreFixtureResultSet) ForEach(f func(*StoreFixture) error) error {
	for {
		var result *StoreFixture
		found, err := r.ResultSet.Next(&result)
		if err != nil {
			return err
		}

		if !found {
			break
		}

		err = f(result)
		if err == storable.ErrStop {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

type StoreWithConstructFixtureStore struct {
	storable.Store
}

func NewStoreWithConstructFixtureStore(db *mgo.Database) *StoreWithConstructFixtureStore {
	return &StoreWithConstructFixtureStore{*storable.NewStore(db, "store_construct")}
}

func (s *StoreWithConstructFixtureStore) New(f string) (doc *StoreWithConstructFixture) {
	doc = newStoreWithConstructFixture(f)
	doc.SetIsNew(true)
	doc.SetId(bson.NewObjectId())
	return
}

func (s *StoreWithConstructFixtureStore) Query() *StoreWithConstructFixtureQuery {
	return &StoreWithConstructFixtureQuery{*storable.NewBaseQuery()}
}

func (s *StoreWithConstructFixtureStore) Find(query *StoreWithConstructFixtureQuery) (*StoreWithConstructFixtureResultSet, error) {
	resultSet, err := s.Store.Find(query)
	if err != nil {
		return nil, err
	}

	return &StoreWithConstructFixtureResultSet{*resultSet}, nil
}

func (s *StoreWithConstructFixtureStore) MustFind(query *StoreWithConstructFixtureQuery) *StoreWithConstructFixtureResultSet {
	resultSet, err := s.Find(query)
	if err != nil {
		panic(err)
	}

	return resultSet
}

func (s *StoreWithConstructFixtureStore) FindOne(query *StoreWithConstructFixtureQuery) (*StoreWithConstructFixture, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

func (s *StoreWithConstructFixtureStore) MustFindOne(query *StoreWithConstructFixtureQuery) *StoreWithConstructFixture {
	doc, err := s.FindOne(query)
	if err != nil {
		panic(err)
	}

	return doc
}

func (s *StoreWithConstructFixtureStore) Insert(doc *StoreWithConstructFixture) error {

	err := s.Store.Insert(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreWithConstructFixtureStore) Update(doc *StoreWithConstructFixture) error {

	err := s.Store.Update(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreWithConstructFixtureStore) Save(doc *StoreWithConstructFixture) (updated bool, err error) {
	updated, err = s.Store.Save(doc)
	if err != nil {
		return false, err
	}

	return
}

type StoreWithConstructFixtureQuery struct {
	storable.BaseQuery
}

func (q *StoreWithConstructFixtureQuery) FindById(ids ...bson.ObjectId) *StoreWithConstructFixtureQuery {
	var vs []interface{}
	for _, id := range ids {
		vs = append(vs, id)
	}
	q.AddCriteria(operators.In(storable.IdField, vs...))

	return q
}

type StoreWithConstructFixtureResultSet struct {
	storable.ResultSet
}

func (r *StoreWithConstructFixtureResultSet) All() ([]*StoreWithConstructFixture, error) {
	var result []*StoreWithConstructFixture
	err := r.ResultSet.All(&result)

	return result, err
}

func (r *StoreWithConstructFixtureResultSet) One() (*StoreWithConstructFixture, error) {
	var result *StoreWithConstructFixture
	err := r.ResultSet.One(&result)

	return result, err
}

func (r *StoreWithConstructFixtureResultSet) Next() (*StoreWithConstructFixture, error) {
	var result *StoreWithConstructFixture
	_, err := r.ResultSet.Next(&result)

	return result, err
}

func (r *StoreWithConstructFixtureResultSet) ForEach(f func(*StoreWithConstructFixture) error) error {
	for {
		var result *StoreWithConstructFixture
		found, err := r.ResultSet.Next(&result)
		if err != nil {
			return err
		}

		if !found {
			break
		}

		err = f(result)
		if err == storable.ErrStop {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

type schema struct {
	EventsFixture             *schemaEventsFixture
	QueryFixture              *schemaQueryFixture
	ResultSetFixture          *schemaResultSetFixture
	SchemaFixture             *schemaSchemaFixture
	StoreFixture              *schemaStoreFixture
	StoreWithConstructFixture *schemaStoreWithConstructFixture
}

type schemaEventsFixture struct {
	Checks storable.Map
}

type schemaQueryFixture struct {
	Foo storable.Field
}

type schemaResultSetFixture struct {
	Foo storable.Field
}

type schemaSchemaFixture struct {
	String         storable.Field
	Int            storable.Field
	Nested         *schemaSchemaFixtureNested
	MapOfString    storable.Map
	MapOfInterface storable.Map
	MapOfSomeType  *schemaSchemaFixtureMapOfSomeType
}

type schemaStoreFixture struct {
	Foo storable.Field
}

type schemaStoreWithConstructFixture struct {
	Foo storable.Field
}

type schemaSchemaFixtureNested struct {
	String         storable.Field
	Int            storable.Field
	Nested         *schemaSchemaFixtureNestedNested
	MapOfString    storable.Map
	MapOfInterface storable.Map
	MapOfSomeType  *schemaSchemaFixtureNestedMapOfSomeType
}

type schemaSchemaFixtureMapOfSomeType struct {
	Foo storable.Map
}

type schemaSchemaFixtureNestedNested struct {
}

type schemaSchemaFixtureNestedMapOfSomeType struct {
	Foo storable.Map
}

var Schema = schema{
	EventsFixture: &schemaEventsFixture{
		Checks: storable.NewMap("checks.[map]", "bool"),
	},
	QueryFixture: &schemaQueryFixture{
		Foo: storable.NewField("foo", "string"),
	},
	ResultSetFixture: &schemaResultSetFixture{
		Foo: storable.NewField("foo", "string"),
	},
	SchemaFixture: &schemaSchemaFixture{
		String: storable.NewField("string", "string"),
		Int:    storable.NewField("foo", "int"),
		Nested: &schemaSchemaFixtureNested{
			String:         storable.NewField("nested.string", "string"),
			Int:            storable.NewField("nested.foo", "int"),
			Nested:         &schemaSchemaFixtureNestedNested{},
			MapOfString:    storable.NewMap("nested.mapofstring.[map]", "string"),
			MapOfInterface: storable.NewMap("nested.mapofinterface.[map]", "interface{}"),
			MapOfSomeType: &schemaSchemaFixtureNestedMapOfSomeType{
				Foo: storable.NewMap("nested.mapofsometype.[map].foo", "string"),
			},
		},
		MapOfString:    storable.NewMap("mapofstring.[map]", "string"),
		MapOfInterface: storable.NewMap("mapofinterface.[map]", "interface{}"),
		MapOfSomeType: &schemaSchemaFixtureMapOfSomeType{
			Foo: storable.NewMap("mapofsometype.[map].foo", "string"),
		},
	},
	StoreFixture: &schemaStoreFixture{
		Foo: storable.NewField("foo", "string"),
	},
	StoreWithConstructFixture: &schemaStoreWithConstructFixture{
		Foo: storable.NewField("foo", "string"),
	},
}
