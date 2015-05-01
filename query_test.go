package storable

import (
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

func (s *BaseSuite) TestBaseQuery_AddCriteria(c *C) {
	foo := bson.M{"foo": "foo"}
	qux := bson.M{"qux": "qux"}

	q := NewBaseQuery()
	q.AddCriteria(foo)
	q.AddCriteria(qux)

	c.Assert(q.GetCriteria(), DeepEquals, bson.M{
		"$and": []bson.M{
			bson.M{"foo": "foo"},
			bson.M{"qux": "qux"},
		},
	})
}
