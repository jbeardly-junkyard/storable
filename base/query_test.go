package base

import (
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

func (s *BaseSuite) TestQuery_AddCriteria(c *C) {
	q := NewBaseQuery()
	q.AddCriteria("foo", "bar")

	c.Assert(q.GetCriteria()["foo"], Equals, "bar")
}

func (s *BaseSuite) TestQuery_FindById(c *C) {
	id := bson.NewObjectId()

	q := NewBaseQuery()
	q.FindById(id)

	c.Assert(q.GetCriteria()["_id"], Equals, id)
}
