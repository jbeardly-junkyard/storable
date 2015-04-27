package operators

import (
	"testing"

	"github.com/tyba/storable"

	"gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

var (
	Foo = storable.NewField("foo", "")
)

func Test(t *testing.T) { check.TestingT(t) }

type OperatorsSuite struct{}

var _ = check.Suite(&OperatorsSuite{})

func (s *OperatorsSuite) TestLogical(c *check.C) {
	or := Or(bson.M{"foo": "qux"}, bson.M{"qux": "qux"})
	c.Assert(or, check.DeepEquals, bson.M{
		"$or": []bson.M{
			bson.M{"foo": "qux"},
			bson.M{"qux": "qux"},
		},
	})

	and := And(bson.M{"foo": "qux"}, bson.M{"qux": "qux"})
	c.Assert(and, check.DeepEquals, bson.M{
		"$and": []bson.M{
			bson.M{"foo": "qux"},
			bson.M{"qux": "qux"},
		},
	})

	not := Not(bson.M{"foo": "qux"})
	c.Assert(not, check.DeepEquals, bson.M{
		"$not": bson.M{"foo": "qux"},
	})

	nor := Nor(bson.M{"foo": "qux"}, bson.M{"qux": "qux"})
	c.Assert(nor, check.DeepEquals, bson.M{
		"$nor": []bson.M{
			bson.M{"foo": "qux"},
			bson.M{"qux": "qux"},
		},
	})
}

func (s *OperatorsSuite) TestComparsion(c *check.C) {
	eq := Eq(Foo, "bar")
	c.Assert(eq, check.DeepEquals, bson.M{"foo": bson.M{"$eq": "bar"}})

	gt := Gt(Foo, "bar")
	c.Assert(gt, check.DeepEquals, bson.M{"foo": bson.M{"$gt": "bar"}})

	gte := Gte(Foo, "bar")
	c.Assert(gte, check.DeepEquals, bson.M{"foo": bson.M{"$gte": "bar"}})

	lt := Lt(Foo, "bar")
	c.Assert(lt, check.DeepEquals, bson.M{"foo": bson.M{"$lt": "bar"}})

	lte := Lte(Foo, "bar")
	c.Assert(lte, check.DeepEquals, bson.M{"foo": bson.M{"$lte": "bar"}})

	ne := Ne(Foo, "bar")
	c.Assert(ne, check.DeepEquals, bson.M{"foo": bson.M{"$ne": "bar"}})

	in := In(Foo, "bar", "qux")
	c.Assert(in, check.DeepEquals, bson.M{"foo": bson.M{"$in": []interface{}{"bar", "qux"}}})

	nin := Nin(Foo, "bar", "qux")
	c.Assert(nin, check.DeepEquals, bson.M{"foo": bson.M{"$nin": []interface{}{"bar", "qux"}}})
}

func (s *OperatorsSuite) TestElement(c *check.C) {
	exists := Exists(Foo, true)
	c.Assert(exists, check.DeepEquals, bson.M{"foo": bson.M{"$exists": true}})

	t := Type(Foo, Double)
	c.Assert(t, check.DeepEquals, bson.M{"foo": bson.M{"$type": Double}})
}

func (s *OperatorsSuite) TestEvaluation(c *check.C) {
	mod := Mod(Foo, 42, 82)
	c.Assert(mod, check.DeepEquals, bson.M{"foo": bson.M{"$mod": []float64{42, 82}}})

	re := RegEx(Foo, ".*", "i")
	c.Assert(re, check.DeepEquals, bson.M{
		"foo": bson.M{"$regex": bson.RegEx{Pattern: ".*", Options: "i"}},
	})

	text := Text(Foo, "foo", "none")
	c.Assert(text, check.DeepEquals, bson.M{
		"foo": bson.M{"$text": bson.M{"$search": "foo", "$language": "none"}},
	})

	where := Where(Foo, "foo", nil)
	c.Assert(where, check.DeepEquals, bson.M{
		"foo": bson.M{"$where": bson.JavaScript{Code: "foo", Scope: interface{}(nil)}},
	})
}

func (s *OperatorsSuite) TestArray(c *check.C) {
	all := All(Foo, "qux", "bar")
	c.Assert(all, check.DeepEquals, bson.M{"foo": bson.M{"$all": []interface{}{"qux", "bar"}}})

	size := Size(Foo, 2)
	c.Assert(size, check.DeepEquals, bson.M{"foo": bson.M{"$size": 2}})
}

func (s *OperatorsSuite) TestComment(c *check.C) {
	comment := Comment("foo")
	c.Assert(comment, check.DeepEquals, bson.M{"$comment": "foo"})
}
