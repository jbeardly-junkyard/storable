package base

import (
	. "gopkg.in/check.v1"
)

func (s *BaseSuite) TestBaseQuery_AddCriteria(c *C) {
	q := NewBaseQuery()
	q.AddCriteria(NewField("foo", ""), "bar")

	c.Assert(q.GetCriteria()["foo"], Equals, "bar")
}
