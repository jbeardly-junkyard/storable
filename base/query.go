package base

import (
	"gopkg.in/mgo.v2/bson"
)

type Query interface {
	GetCriteria() bson.M
}

type BaseQuery struct {
	criteria bson.M
}

func NewBaseQuery() *BaseQuery {
	return &BaseQuery{criteria: make(bson.M, 0)}
}

func (q *BaseQuery) FindById(id bson.ObjectId) {
	q.AddCriteria("_id", id)
}

func (q *BaseQuery) AddCriteria(key string, val interface{}) {
	q.criteria[key] = val
}

func (q *BaseQuery) GetCriteria() bson.M {
	return q.criteria
}
