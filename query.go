package storable

import (
	"encoding/json"

	"github.com/tyba/storable/operators"

	"gopkg.in/mgo.v2/bson"
)

type Query interface {
	GetCriteria() bson.M
	Sort(s Sort)
	Limit(l int)
	Skip(s int)
	GetSort() Sort
	GetLimit() int
	GetSkip() int
}

type BaseQuery struct {
	clauses     []bson.M
	limit, skip int
	sort        Sort
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

func (q *BaseQuery) Sort(s Sort) {
	q.sort = s
}

func (q *BaseQuery) Limit(l int) {
	q.limit = l
}

func (q *BaseQuery) Skip(s int) {
	q.skip = s
}

func (q *BaseQuery) GetSort() Sort {
	return q.sort
}

func (q *BaseQuery) GetLimit() int {
	return q.limit
}

func (q *BaseQuery) GetSkip() int {
	return q.skip
}

func (q *BaseQuery) String() string {
	j, _ := json.Marshal(q.GetCriteria())

	return string(j)
}
