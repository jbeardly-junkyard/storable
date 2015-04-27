package operators

import (
	"github.com/tyba/mongogen/base"

	"gopkg.in/mgo.v2/bson"
)

// Eq Matches values that are equal to a specified value.
func Eq(field base.Field, value interface{}) bson.M {
	return bson.M{field.String(): bson.M{"$eq": value}}
}

// Gt Matches values that are greater than a specified value.
func Gt(field base.Field, value interface{}) bson.M {
	return bson.M{field.String(): bson.M{"$gt": value}}
}

// Gte Matches values that are greater than or equal to a specified value.
func Gte(field base.Field, value interface{}) bson.M {
	return bson.M{field.String(): bson.M{"$gte": value}}
}

// Lt Matches values that are less than a specified value.
func Lt(field base.Field, value interface{}) bson.M {
	return bson.M{field.String(): bson.M{"$lt": value}}
}

// Lte Matches values that are less than or equal to a specified value.
func Lte(field base.Field, value interface{}) bson.M {
	return bson.M{field.String(): bson.M{"$lte": value}}
}

// Ne Matches all values that are not equal to a specified value.
func Ne(field base.Field, value interface{}) bson.M {
	return bson.M{field.String(): bson.M{"$ne": value}}
}

// In Matches any of the values.
func In(field base.Field, values ...interface{}) bson.M {
	return bson.M{field.String(): bson.M{"$in": values}}
}

// Nin Matches none of the values.
func Nin(field base.Field, values ...interface{}) bson.M {
	return bson.M{field.String(): bson.M{"$nin": values}}
}
