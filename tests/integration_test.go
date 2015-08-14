package example

import (
	"errors"
	"testing"

	"github.com/tyba/storable"
	"github.com/tyba/storable/operators"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

const (
	testMongoHost = "127.0.0.1:27017"
	testDatabase  = "storable-test"
)

func Test(t *testing.T) { TestingT(t) }

type MongoSuite struct {
	db *mgo.Database
}

var _ = Suite(&MongoSuite{})

func (s *MongoSuite) SetUpTest(c *C) {
	conn, _ := mgo.Dial(testMongoHost)
	s.db = conn.DB(testDatabase)
}

func (s *MongoSuite) TestStore_New(c *C) {
	store := NewMyModelStore(s.db)
	m := store.New()

	c.Assert(m.IsNew(), Equals, true)
}

func (s *MongoSuite) TestQuery_FindByFoo(c *C) {
	store := NewMyModelStore(s.db)
	m := store.New()
	m.String = "foo"

	c.Assert(store.Insert(m), IsNil)

	q := store.Query()
	q.AddCriteria(operators.Eq(Schema.MyModel.String, "foo"))

	r, err := store.Find(q)
	c.Assert(err, IsNil)

	res, err := r.All()
	c.Assert(res, HasLen, 1)
	c.Assert(err, IsNil)

	q.AddCriteria(operators.Eq(Schema.MyModel.String, "bar"))
	one, err := store.MustFind(q).One()
	c.Assert(one, IsNil)
	c.Assert(err, Equals, storable.ErrNotFound)
}

func (s *MongoSuite) TestSchema(c *C) {
	c.Assert(Schema.MyModel.String.String(), Equals, "string")
	c.Assert(Schema.MyModel.Int.String(), Equals, "bla2")
	c.Assert(Schema.MyModel.Nested.X.String(), Equals, "nested.x")

	key := Schema.MyModel.Nested.Another.X.String()
	c.Assert(key, Equals, "nested.another.x")

	key = Schema.MyModel.MapsOfString.Key("foo").String()
	c.Assert(key, Equals, "mapsofstring.foo")

	key = Schema.MyModel.InlineStruct.MapOfString.Key("qux").String()
	c.Assert(key, Equals, "inlinestruct.mapofstring.qux")

	key = Schema.MyModel.InlineStruct.MapOfSomeType.X.Key("foo").String()
	c.Assert(key, Equals, "inlinestruct.mapofsometype.foo.x")

	key = Schema.MyModel.InlineStruct.MapOfInterface.Key("foo").String()
	c.Assert(key, Equals, "inlinestruct.mapofinterface.foo")
}

func (s *MongoSuite) TestEventsInsert(c *C) {
	store := NewEventsTestsStore(s.db)

	doc := store.New()
	err := store.Insert(doc)
	c.Assert(err, IsNil)
	c.Assert(doc.Checks, DeepEquals, map[string]bool{
		"BeforeInsert": true,
		"AfterInsert":  true,
	})
}

func (s *MongoSuite) TestEventsUpdate(c *C) {
	store := NewEventsTestsStore(s.db)

	doc := store.New()
	err := store.Insert(doc)
	c.Assert(err, IsNil)

	doc.Checks = make(map[string]bool, 0)
	err = store.Update(doc)
	c.Assert(err, IsNil)
	c.Assert(doc.Checks, DeepEquals, map[string]bool{
		"BeforeUpdate": true,
		"AfterUpdate":  true,
	})
}

func (s *MongoSuite) TestEventsUpdateError(c *C) {
	store := NewEventsTestsStore(s.db)

	doc := store.New()
	err := store.Insert(doc)
	doc.Checks = make(map[string]bool, 0)

	doc.MustFailAfter = errors.New("after")
	err = store.Update(doc)
	c.Assert(err, Equals, doc.MustFailAfter)

	doc.MustFailBefore = errors.New("before")
	err = store.Update(doc)
	c.Assert(err, Equals, doc.MustFailBefore)
}

func (s *MongoSuite) TestEventsSaveInsert(c *C) {
	store := NewEventsTestsStore(s.db)

	doc := store.New()
	updated, err := store.Save(doc)
	c.Assert(err, IsNil)
	c.Assert(updated, Equals, false)
	c.Assert(doc.Checks, DeepEquals, map[string]bool{
		"BeforeInsert": true,
		"AfterInsert":  true,
	})
}

func (s *MongoSuite) TestEventsSaveUpdate(c *C) {
	store := NewEventsTestsStore(s.db)

	doc := store.New()
	err := store.Insert(doc)
	doc.Checks = make(map[string]bool, 0)

	updated, err := store.Save(doc)
	c.Assert(err, IsNil)
	c.Assert(updated, Equals, true)
	c.Assert(doc.Checks, DeepEquals, map[string]bool{
		"BeforeUpdate": true,
		"AfterUpdate":  true,
	})
}

func (s *MongoSuite) TearDownTest(c *C) {
	s.db.DropDatabase()
}
