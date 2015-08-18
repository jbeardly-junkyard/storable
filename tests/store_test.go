package tests

import . "gopkg.in/check.v1"

func (s *MongoSuite) TestStoreNew(c *C) {
	store := NewStoreFixtureStore(s.db)
	doc := store.New()

	c.Assert(doc.IsNew(), Equals, true)
	c.Assert(doc.GetId().Hex(), HasLen, 24)
}

func (s *MongoSuite) TestStoreQuery(c *C) {
	store := NewStoreFixtureStore(s.db)
	q := store.Query()
	c.Assert(q, Not(IsNil))
}

func (s *MongoSuite) TestStoreFind(c *C) {
	store := NewStoreFixtureStore(s.db)
	c.Assert(store.Insert(store.New()), IsNil)
	c.Assert(store.Insert(store.New()), IsNil)

	rs, err := store.Find(store.Query())
	c.Assert(err, IsNil)

	count, err := rs.Count()
	c.Assert(err, IsNil)
	c.Assert(count, Equals, 2)
}

func (s *MongoSuite) TestStoreMustFind(c *C) {
	store := NewStoreFixtureStore(s.db)
	c.Assert(store.Insert(store.New()), IsNil)
	c.Assert(store.Insert(store.New()), IsNil)

	count, err := store.MustFind(store.Query()).Count()
	c.Assert(err, IsNil)
	c.Assert(count, Equals, 2)
}

func (s *MongoSuite) TestStoreFindOne(c *C) {
	store := NewStoreWithConstructFixtureStore(s.db)
	c.Assert(store.Insert(store.New("bar")), IsNil)

	doc, err := store.FindOne(store.Query())
	c.Assert(err, IsNil)
	c.Assert(doc.Foo, Equals, "bar")
}

func (s *MongoSuite) TestStoreMustFindOne(c *C) {
	store := NewStoreWithConstructFixtureStore(s.db)
	c.Assert(store.Insert(store.New("foo")), IsNil)
	c.Assert(store.MustFindOne(store.Query()).Foo, Equals, "foo")
}

func (s *MongoSuite) TestStoreInsertUpdate(c *C) {
	store := NewStoreWithConstructFixtureStore(s.db)

	doc := store.New("foo")
	err := store.Insert(doc)
	c.Assert(err, IsNil)
	c.Assert(store.MustFindOne(store.Query()).Foo, Equals, "foo")

	doc.Foo = "bar"
	err = store.Update(doc)
	c.Assert(err, IsNil)
	c.Assert(store.MustFindOne(store.Query()).Foo, Equals, "bar")
}

func (s *MongoSuite) TestStoreSave(c *C) {
	store := NewStoreWithConstructFixtureStore(s.db)

	doc := store.New("foo")
	updated, err := store.Save(doc)
	c.Assert(err, IsNil)
	c.Assert(updated, Equals, false)
	c.Assert(doc.IsNew(), Equals, false)
	c.Assert(store.MustFindOne(store.Query()).Foo, Equals, "foo")

	doc.Foo = "bar"
	updated, err = store.Save(doc)
	c.Assert(err, IsNil)
	c.Assert(updated, Equals, true)
	c.Assert(store.MustFindOne(store.Query()).Foo, Equals, "bar")
}
