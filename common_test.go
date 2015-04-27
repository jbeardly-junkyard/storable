package storable

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

type BaseSuite struct {
	db *mgo.Database
}

var _ = Suite(&BaseSuite{})

func (s *BaseSuite) SetUpTest(c *C) {
	conn, _ := mgo.Dial(testMongoHost)
	s.db = conn.DB(testDatabase)
}

func (s *BaseSuite) TestSort_String(c *C) {
	sort := Sort{{Field{"foo", ""}, Asc}}
	c.Assert(sort.String(), Equals, "foo")

	sort = Sort{{Field{"foo", ""}, Desc}}
	c.Assert(sort.String(), Equals, "-foo")

	sort = Sort{{Field{"foo", ""}, Asc}, {Field{"qux", ""}, Desc}}
	c.Assert(sort.String(), Equals, "foo,-qux")
}

func (s *BaseSuite) TestSort_IsEmpty(c *C) {
	sort := Sort{{Field{"foo", ""}, Asc}}
	c.Assert(sort.IsEmpty(), Equals, false)

	sort = Sort{}
	c.Assert(sort.IsEmpty(), Equals, true)
}

func (s *BaseSuite) TearDownTest(c *C) {
	s.db.DropDatabase()
}

type Person struct {
	Document  `bson:",inline"`
	FirstName string
	LastName  string
	Gender    string
}
