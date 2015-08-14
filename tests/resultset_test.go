package tests

import . "gopkg.in/check.v1"

func (s *MongoSuite) TestResultSetAll(c *C) {
	store := NewResultSetFixtureStore(s.db)
	c.Assert(store.Insert(store.New("bar")), IsNil)
	c.Assert(store.Insert(store.New("foo")), IsNil)

	docs, err := store.MustFind(store.Query()).All()
	c.Assert(err, IsNil)
	c.Assert(docs, HasLen, 2)
}

func (s *MongoSuite) TestResultSetOne(c *C) {
	store := NewResultSetFixtureStore(s.db)
	c.Assert(store.Insert(store.New("bar")), IsNil)

	doc, err := store.MustFind(store.Query()).One()
	c.Assert(err, IsNil)
	c.Assert(doc.Foo, Equals, "bar")
}

func (s *MongoSuite) TestResultSetNext(c *C) {
	store := NewResultSetFixtureStore(s.db)
	c.Assert(store.Insert(store.New("bar")), IsNil)

	rs := store.MustFind(store.Query())
	doc, err := rs.Next()
	c.Assert(err, IsNil)
	c.Assert(doc.Foo, Equals, "bar")

	doc, err = rs.Next()
	c.Assert(err, IsNil)
	c.Assert(doc, IsNil)
}
