package storable

import (
	"errors"

	"gopkg.in/mgo.v2"
)

var (
	ResultSetClosed = errors.New("Cannot close a closed resultset.")
)

type ResultSet struct {
	IsClosed bool
	mgoQuery *mgo.Query
	mgoIter  *mgo.Iter
}

func (r *ResultSet) Count() (int, error) {
	return r.mgoQuery.Count()
}

func (r *ResultSet) All(result interface{}) error {
	defer r.Close()
	return r.mgoQuery.All(result)
}

func (r *ResultSet) One(doc interface{}) (bool, error) {
	defer r.Close()
	return r.Next(doc)
}

func (r *ResultSet) Next(doc interface{}) (bool, error) {
	if r.mgoIter == nil {
		r.mgoIter = r.mgoQuery.Iter()
	}

	returned := r.mgoIter.Next(doc)
	if base, ok := doc.(DocumentBase); ok && returned {
		base.SetIsNew(false)
	}

	return returned, r.mgoIter.Err()
}

func (r *ResultSet) Close() error {
	if r.IsClosed {
		return ResultSetClosed
	}

	r.IsClosed = true
	if r.mgoIter == nil {
		return nil
	}

	return r.mgoIter.Close()
}
