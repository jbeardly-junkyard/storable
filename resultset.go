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

// Count returns the total number of documents in the ResultSet.
func (r *ResultSet) Count() (int, error) {
	return r.mgoQuery.Count()
}

// All returns all the documents in the ResultSet and close it. Dont use it
// with large results.
func (r *ResultSet) All(result interface{}) error {
	defer r.Close()
	return r.mgoQuery.All(result)
}

// One return a document from the ResultSet and close it, the following calls
// to One returns ResultSetClosed error.
func (r *ResultSet) One(doc interface{}) (bool, error) {
	defer r.Close()
	return r.Next(doc)
}

// Next return a document from the ResultSet, can be called multiple times.
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

// Close close the ResultSet closing the internal iter.
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
