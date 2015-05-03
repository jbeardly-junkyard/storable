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

// AddCriteria adds a new mathing expression to the query, all the expressions
// are merged on a $and expression.
//
// Use operators package instead of build expresion by hand:
//
//  import . "github.com/tyba/storable/operators"
//
//  func (q *YourQuery) FindNonZeroRecords() {
//      // All the Fields are defined on the Schema generated variable
//      size := NewField("size", "int")
//      q.AddCriteria(Gt(size, 0))
//  }
func (q *BaseQuery) AddCriteria(expr bson.M) {
	q.clauses = append(q.clauses, expr)
}

// GetCriteria returns a valid bson.M used internally by Store.
func (q *BaseQuery) GetCriteria() bson.M {
	if len(q.clauses) == 0 {
		return nil
	}

	return operators.And(q.clauses...)
}

// Sort sets the sorting cristeria of the query.
func (q *BaseQuery) Sort(s Sort) {
	q.sort = s
}

// Limit sets the limit of the query.
func (q *BaseQuery) Limit(l int) {
	q.limit = l
}

// Skip sets the skip of the query.
func (q *BaseQuery) Skip(s int) {
	q.skip = s
}

// GetSort return the current sorting preferences of the query.
func (q *BaseQuery) GetSort() Sort {
	return q.sort
}

// GetSort return the current limit preferences of the query.
func (q *BaseQuery) GetLimit() int {
	return q.limit
}

// GetSort return the current skip preferences of the query.
func (q *BaseQuery) GetSkip() int {
	return q.skip
}

// Strings return a json representation of the criteria. Sorry but this is not
// fully compatible with the MongoDb CLI.
func (q *BaseQuery) String() string {
	j, _ := json.Marshal(q.GetCriteria())

	return string(j)
}
