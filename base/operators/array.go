package operators

import (
	"github.com/tyba/mongogen/base"

	"gopkg.in/mgo.v2/bson"
)

// All Matches arrays that contain all elements specified in the query.
func All(field base.Field, values ...interface{}) bson.M {
	return bson.M{field.String(): bson.M{"$all": values}}
}

// Size Selects documents if the array field is a specified size.
func Size(field base.Field, count int) bson.M {
	return bson.M{field.String(): bson.M{"$size": count}}
}

//TODO: $elemMatch
