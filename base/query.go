package base

import (
	"gopkg.in/mgo.v2/bson"
)

type Query struct {
	criteria    bson.M
	Limit, Skip int
	Sort        Sort
}

func NewQuery() *Query {
	return &Query{criteria: make(bson.M, 0)}
}

func (q *Query) FindById(id bson.ObjectId) {
	q.AddCriteria(IdField, id)
}

func (q *Query) AddCriteria(key Field, val interface{}) {
	q.criteria[key.String()] = val
}

func (q *Query) GetCriteria() bson.M {
	return q.criteria
}
