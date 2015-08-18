package tests

import (
	"errors"

	"gopkg.in/tyba/storable.v1"
	. "gopkg.in/check.v1"
)

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

func (s *MongoSuite) TestResultSetNextEmpty(c *C) {
	store := NewResultSetFixtureStore(s.db)
	rs := store.MustFind(store.Query())
	returned := rs.Next()
	c.Assert(returned, Equals, false)

	doc, err := rs.Get()
	c.Assert(err, IsNil)
	c.Assert(doc, IsNil)
}

func (s *MongoSuite) TestResultSetNext(c *C) {
	store := NewResultSetFixtureStore(s.db)
	c.Assert(store.Insert(store.New("bar")), IsNil)

	rs := store.MustFind(store.Query())
	returned := rs.Next()
	c.Assert(returned, Equals, true)

	doc, err := rs.Get()
	c.Assert(err, IsNil)
	c.Assert(doc.Foo, Equals, "bar")

	returned = rs.Next()
	c.Assert(returned, Equals, false)

	doc, err = rs.Get()
	c.Assert(err, IsNil)
	c.Assert(doc, IsNil)
}

func (s *MongoSuite) TestResultSetForEach(c *C) {
	store := NewResultSetFixtureStore(s.db)
	c.Assert(store.Insert(store.New("bar")), IsNil)
	c.Assert(store.Insert(store.New("foo")), IsNil)

	count := 0
	err := store.MustFind(store.Query()).ForEach(func(*ResultSetFixture) error {
		count++
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(count, Equals, 2)
}

func (s *MongoSuite) TestResultSetForEachStop(c *C) {
	store := NewResultSetFixtureStore(s.db)
	c.Assert(store.Insert(store.New("bar")), IsNil)
	c.Assert(store.Insert(store.New("foo")), IsNil)

	count := 0
	err := store.MustFind(store.Query()).ForEach(func(*ResultSetFixture) error {
		count++
		return storable.ErrStop
	})

	c.Assert(err, IsNil)
	c.Assert(count, Equals, 1)
}

func (s *MongoSuite) TestResultSetForEachError(c *C) {
	store := NewResultSetFixtureStore(s.db)
	c.Assert(store.Insert(store.New("bar")), IsNil)
	c.Assert(store.Insert(store.New("foo")), IsNil)

	fail := errors.New("foo")
	err := store.MustFind(store.Query()).ForEach(func(*ResultSetFixture) error {
		return fail
	})

	c.Assert(err, Equals, fail)
}
