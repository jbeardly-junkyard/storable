package base

import (
	. "gopkg.in/check.v1"
)

func (s *BaseSuite) TestResultSet_All(c *C) {
	st := NewStore(s.conn, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	r, err := st.Find(NewBaseQuery())
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 2)

	c.Assert(r.IsClosed, Equals, true)
}

func (s *BaseSuite) TestResultSet_One(c *C) {
	st := NewStore(s.conn, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	r, err := st.Find(NewBaseQuery())
	c.Assert(err, IsNil)

	var result *Person
	f, err := r.One(&result)

	c.Assert(err, IsNil)
	c.Assert(f, Equals, true)
	c.Assert(r.IsClosed, Equals, true)
}

func (s *BaseSuite) TestResultSet_Next(c *C) {
	st := NewStore(s.conn, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	r, err := st.Find(NewBaseQuery())
	c.Assert(err, IsNil)

	var result *Person
	f, err := r.Next(&result)

	c.Assert(err, IsNil)
	c.Assert(f, Equals, true)
	c.Assert(r.IsClosed, Equals, false)
}

func (s *BaseSuite) TestResultSet_Close(c *C) {
	st := NewStore(s.conn, "test")
	r, _ := st.Find(NewBaseQuery())

	c.Assert(r.Close(), IsNil)
	c.Assert(r.Close(), Equals, ResultSetClosed)
}
