package base

import (
	"errors"

	"gopkg.in/maxwellhealth/bongo.v0"
)

var (
	ResultSetClosed = errors.New("Cannot close a closed resultset.")
)

type ResultSet struct {
	IsClosed bool
	rs       *bongo.ResultSet
}

func (r *ResultSet) All(result interface{}) error {
	defer r.Close()
	return r.rs.Query.All(result)
}

func (r *ResultSet) One(doc interface{}) (bool, error) {
	defer r.Close()
	return r.Next(doc)
}

func (r *ResultSet) Next(doc interface{}) (bool, error) {
	returned := r.rs.Next(doc)
	if r.rs.Error != nil {
		return false, r.rs.Error
	}

	return returned, nil
}

func (r *ResultSet) Close() error {
	if r.IsClosed {
		return ResultSetClosed
	}

	r.IsClosed = true
	return r.rs.Free()
}
