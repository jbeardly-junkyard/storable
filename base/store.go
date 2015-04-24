package base

import (
	"errors"

	"gopkg.in/maxwellhealth/bongo.v0"
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

func (s *Store) Find(q *Query) (*ResultSet, error) {
	resultSet := s.coll.Find(q.GetCriteria())
	if resultSet.Error != nil {
		return nil, resultSet.Error
	}

	if !q.Sort.IsEmpty() {
		resultSet.Query.Sort(q.Sort.String())
	}

	if q.Skip != 0 {
		resultSet.Query.Skip(q.Skip)
	}

	if q.Limit != 0 {
		resultSet.Query.Limit(q.Limit)
	}

	return &ResultSet{rs: resultSet}, nil
}

func (s *Store) RawUpdate(query *Query, update interface{}) error {
	return s.coll.Collection().Update(query.GetCriteria(), update)
}
