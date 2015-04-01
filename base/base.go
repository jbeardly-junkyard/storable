package base

import (
	"errors"

	"gopkg.in/maxwellhealth/bongo.v0"
	"gopkg.in/mgo.v2/bson"
)

var (
	NonNewDocumentErr = errors.New("Cannot insert a non new document.")
	NewDocumentErr    = errors.New("Cannot updated a new document.")
)

type Store struct {
	Collection string

	conn *bongo.Connection
	coll *bongo.Collection
}

func NewStore(conn *bongo.Connection, collection string) *Store {
	return &Store{
		Collection: collection,
		conn:       conn,
		coll:       conn.Collection(collection),
	}
}

func (s *Store) Insert(doc bongo.Document) error {
	if !s.isNew(doc) {
		return NonNewDocumentErr
	}

	return s.coll.Save(doc)
}

func (s *Store) Update(doc bongo.Document) error {
	if s.isNew(doc) {
		return NewDocumentErr
	}

	return s.coll.Save(doc)
}

func (s *Store) isNew(doc bongo.Document) bool {
	isNew := true
	if newt, ok := doc.(bongo.NewTracker); ok {
		isNew = newt.IsNew()
	}

	return isNew
}

func (s *Store) Delete(doc bongo.Document) error {
	return s.coll.DeleteDocument(doc)
}

func (s *Store) Find(query Query) (*ResultSet, error) {
	resultSet := s.coll.Find(query.GetCriteria())
	if resultSet.Error != nil {
		return nil, resultSet.Error
	}

	return &ResultSet{resultSet}, nil
}

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

type ResultSet struct {
	rs *bongo.ResultSet
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
	return r.rs.Free()
}
