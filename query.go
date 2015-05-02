package storable

import (
	"encoding/json"

	"github.com/tyba/storable/operators"

	"gopkg.in/mgo.v2/bson"
)

type Query interface {
	GetCriteria() bson.M
	GetSort() Sort
	GetLimit() int
	GetSkip() int
}

type BaseQuery struct {
	clauses     []bson.M
	Limit, Skip int
	Sort        Sort
}

func NewBaseQuery() *BaseQuery {
	return &BaseQuery{clauses: make([]bson.M, 0)}
}

func (q *BaseQuery) AddCriteria(expr bson.M) {
	q.clauses = append(q.clauses, expr)
}

func (q *BaseQuery) GetCriteria() bson.M {
	if len(q.clauses) == 0 {
		return nil
	}

	return operators.And(q.clauses...)
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

func (q *BaseQuery) String() string {
	j, _ := json.Marshal(q.GetCriteria())

	return string(j)
}
