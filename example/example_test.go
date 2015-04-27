package example

import (
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
)

const (
	testMongoHost = "localhost"
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

func (s *MongoSuite) TestQuery_FindByFoo(c *C) {
	store := NewMyModelStore(s.db)
	m := store.New()
	m.Foo = "foo"

	c.Assert(store.Insert(m), IsNil)

	q := store.Query()
	q.FindByFoo("foo")

	r, err := store.Find(q)
	c.Assert(err, IsNil)

	res, err := r.All()
	c.Assert(res, HasLen, 1)
	c.Assert(err, IsNil)

	q.FindByFoo("bar")
	r, err = store.Find(q)
	c.Assert(err, IsNil)

	one, err := r.One()
	c.Assert(one, IsNil)
	c.Assert(err, IsNil)
}

func (s *MongoSuite) TestSchema(c *C) {
	c.Assert(Schema.MyModel.Foo.String(), Equals, "foo")
	c.Assert(Schema.MyModel.Bar.String(), Equals, "bla2")
	c.Assert(Schema.MyModel.Nested.X.String(), Equals, "nested.x")
	c.Assert(Schema.MyModel.Nested.Another.X.String(), Equals, "nested.another.x")
}

func (s *MongoSuite) TearDownTest(c *C) {
	s.db.DropDatabase()
}
