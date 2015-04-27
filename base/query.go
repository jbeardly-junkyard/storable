package base

import (
	"gopkg.in/mgo.v2/bson"
)

type Query interface {
	GetCriteria() bson.M
	GetSort() Sort
	GetLimit() int
	GetSkip() int
}

type BaseQuery struct {
	criteria    bson.M
	Limit, Skip int
	Sort        Sort
}

func NewBaseQuery() *BaseQuery {
	return &BaseQuery{criteria: make(bson.M, 0)}
}

func (q *BaseQuery) AddCriteria(key Field, val interface{}) {
	q.criteria[key.String()] = val
}

func (q *BaseQuery) GetCriteria() bson.M {
	return q.criteria
}

func (q *BaseQuery) GetSort() Sort {
	return q.Sort
}

func (q *BaseQuery) GetLimit() int {
	return q.Limit
}

func (q *BaseQuery) GetSkip() int {
	return q.Skip
}
