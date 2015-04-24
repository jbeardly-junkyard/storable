package base

import (
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/maxwellhealth/bongo.v0"
	"gopkg.in/mgo.v2/bson"
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

func (s *BaseSuite) TestStore_Insert(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.conn, "test")
	err := st.Insert(p)
	c.Assert(err, IsNil)
	c.Assert(p.IsNew(), Equals, false)

	r, err := st.Find(NewBaseQuery())
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "foo")
}

func (s *BaseSuite) TestStore_InsertOld(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.conn, "test")
	err := st.Insert(p)
	c.Assert(err, IsNil)

	err = st.Insert(p)
	c.Assert(err, Equals, NonNewDocumentErr)
}

func (s *BaseSuite) TestStore_Update(c *C) {
	p := &Person{FirstName: "foo"}

	st := NewStore(s.conn, "test")
	st.Insert(p)
	st.Insert(&Person{FirstName: "bar"})

	p.FirstName = "qux"
	err := st.Update(p)
	c.Assert(err, IsNil)

	q := NewBaseQuery()
	q.AddCriteria("firstname", "qux")

	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "qux")
}

func (s *BaseSuite) TestStore_UpdateNew(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.conn, "test")

	err := st.Update(p)
	c.Assert(err, Equals, NewDocumentErr)
}

func (s *BaseSuite) TestStore_Delete(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.conn, "test")
	st.Insert(p)

	err := st.Delete(p)
	c.Assert(err, IsNil)

	r, err := st.Find(NewBaseQuery())
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 0)
}

func (s *BaseSuite) TestStore_RawUpdate(c *C) {
	st := NewStore(s.conn, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	q := NewBaseQuery()
	q.AddCriteria("firstname", "foo")

	err := st.RawUpdate(q, bson.M{"firstname": "qux"})
	c.Assert(err, IsNil)

	q = NewBaseQuery()
	q.AddCriteria("firstname", "qux")

	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "qux")
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
