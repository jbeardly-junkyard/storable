package storable

import (
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

func (s *BaseSuite) TestStore_Insert(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.db, "test")
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
	st := NewStore(s.db, "test")
	err := st.Insert(p)
	c.Assert(err, IsNil)

	err = st.Insert(p)
	c.Assert(err, Equals, NonNewDocumentErr)
}

func (s *BaseSuite) TestStore_Update(c *C) {
	p := &Person{FirstName: "foo"}

	st := NewStore(s.db, "test")
	st.Insert(p)
	st.Insert(&Person{FirstName: "bar"})

	p.FirstName = "qux"
	err := st.Update(p)
	c.Assert(err, IsNil)

	q := NewBaseQuery()
	q.AddCriteria(bson.M{"firstname": "qux"})

	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "qux")
}

func (s *BaseSuite) TestStore_UpdateNew(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.db, "test")

	err := st.Update(p)
	c.Assert(err, Equals, NewDocumentErr)
}

func (s *BaseSuite) TestStore_Delete(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.db, "test")
	st.Insert(p)

	err := st.Delete(p)
	c.Assert(err, IsNil)

	r, err := st.Find(NewBaseQuery())
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 0)
}

func (s *BaseSuite) TestStore_FindLimit(c *C) {
	st := NewStore(s.db, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	q := NewBaseQuery()
	q.Limit = 1
	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "foo")
}

func (s *BaseSuite) TestStore_FindSkip(c *C) {
	st := NewStore(s.db, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	q := NewBaseQuery()
	q.Skip = 1
	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "bar")
}

func (s *BaseSuite) TestStore_FindSort(c *C) {
	st := NewStore(s.db, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	q := NewBaseQuery()
	q.Sort = Sort{{IdField, Desc}}
	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 2)
	c.Assert(result[0].FirstName, Equals, "bar")
	c.Assert(result[1].FirstName, Equals, "foo")
}

func (s *BaseSuite) TestStore_RawUpdate(c *C) {
	st := NewStore(s.db, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	q := NewBaseQuery()
	q.AddCriteria(bson.M{"firstname": "foo"})

	err := st.RawUpdate(q, bson.M{"firstname": "qux"})
	c.Assert(err, IsNil)

	q = NewBaseQuery()
	q.AddCriteria(bson.M{"firstname": "qux"})

	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "qux")
}

func (s *BaseSuite) TestStore_RawUpdateEmpty(c *C) {
	st := NewStore(s.db, "test")
	q := NewBaseQuery()
	err := st.RawUpdate(q, bson.M{"firstname": "qux"})
	c.Assert(err, Equals, EmptyQueryInRawErr)
}
