package base

import (
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/maxwellhealth/bongo.v0"
)

const (
	testDatabase = "mongogen-test"
)

func Test(t *testing.T) { TestingT(t) }

type BaseSuite struct {
	conn *bongo.Connection
}

var _ = Suite(&BaseSuite{})

func (s *BaseSuite) SetUpTest(c *C) {
	s.conn, _ = bongo.Connect(&bongo.Config{
		ConnectionString: "localhost",
		Database:         testDatabase,
	})
}

func (s *BaseSuite) TestSort_String(c *C) {
	sort := Sort{{Field{"foo"}, Asc}}
	c.Assert(sort.String(), Equals, "foo")

	sort = Sort{{Field{"foo"}, Desc}}
	c.Assert(sort.String(), Equals, "-foo")

	sort = Sort{{Field{"foo"}, Asc}, {Field{"qux"}, Desc}}
	c.Assert(sort.String(), Equals, "foo,-qux")
}

func (s *BaseSuite) TestSort_IsEmpty(c *C) {
	sort := Sort{{Field{"foo"}, Asc}}
	c.Assert(sort.IsEmpty(), Equals, false)

	sort = Sort{}
	c.Assert(sort.IsEmpty(), Equals, true)
}

func (s *BaseSuite) TearDownTest(c *C) {
	s.conn.Session.DB(testDatabase).DropDatabase()
}

type Person struct {
	bongo.DocumentBase `bson:",inline"`
	FirstName          string
	LastName           string
	Gender             string
}
