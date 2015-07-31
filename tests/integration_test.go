package example

import (
	"testing"

	"github.com/tyba/storable"
	"github.com/tyba/storable/operators"

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
}

func (s *MongoSuite) TearDownTest(c *C) {
	s.db.DropDatabase()
}
