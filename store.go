package storable

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	NonNewDocumentErr  = errors.New("Cannot insert a non new document.")
	NewDocumentErr     = errors.New("Cannot updated a new document.")
	EmptyQueryInRawErr = errors.New("Empty queries are not allowed on raw ops.")
	EmptyIdErr         = errors.New("A document without id is not allowed.")
)

type Store struct {
	Collection *mgo.Collection
}

// NewStore returns a new Store instance
func NewStore(db *mgo.Database, collection string) *Store {
	return &Store{
		Collection: db.C(collection),
	}
}

// Insert insert the given document in the collection, returns error if no-new
// document is given. The document id is setted if is empty.
func (s *Store) Insert(doc DocumentBase) error {
	if !doc.IsNew() {
		return NonNewDocumentErr
	}

	if len(doc.GetId()) == 0 {
		doc.SetId(bson.NewObjectId())
	}

	err := s.Collection.Insert(doc)
	if err == nil {
		doc.SetIsNew(false)
	}

	return err
}

// Update update the given document in the collection, returns error if a new
// document is given.
func (s *Store) Update(doc DocumentBase) error {
	if doc.IsNew() {
		return NewDocumentErr
	}

	return s.Collection.UpdateId(doc.GetId(), doc)
}

// Save insert or update the given document in the collection, a document with
// id should be provided. Upsert is used (http://godoc.org/gopkg.in/mgo.v2#Collection.Upsert)
func (s *Store) Save(doc DocumentBase) error {
	id := doc.GetId()
	if len(id) == 0 {
		return EmptyIdErr
	}

	_, err := s.Collection.UpsertId(id, doc)
	if err == nil {
		doc.SetIsNew(false)
	}

	return err
}

// Delete remove the document from the collection
func (s *Store) Delete(doc DocumentBase) error {
	return s.Collection.RemoveId(doc.GetId())
}

// Find executes the given query in the collection
func (s *Store) Find(q Query) (*ResultSet, error) {
	mq := s.Collection.Find(q.GetCriteria())

	if !q.GetSort().IsEmpty() {
		mq.Sort(q.GetSort().String())
	}

	if q.GetSkip() != 0 {
		mq.Skip(q.GetSkip())
	}

	if q.GetLimit() != 0 {
		mq.Limit(q.GetLimit())
	}

	return &ResultSet{mgoQuery: mq}, nil
}

// RawUpdate performes a direct update in the collection, update is wrapped on
// a $set operator. If a query without criteria is given EmptyQueryInRawErr is
// returned
func (s *Store) RawUpdate(query Query, update interface{}, multi bool) error {
	criteria := query.GetCriteria()
	if len(criteria) == 0 {
		return EmptyQueryInRawErr
	}

	var err error
	if multi {
		_, err = s.Collection.UpdateAll(criteria, bson.M{"$set": update})
	} else {
		err = s.Collection.Update(criteria, bson.M{"$set": update})
	}

	return err
}

// RawDelete performes a direct remove in the collection. If a query without
// criteria is given EmptyQueryInRawErr is returned
func (s *Store) RawDelete(query Query, multi bool) error {
	criteria := query.GetCriteria()
	if len(criteria) == 0 {
		return EmptyQueryInRawErr
	}

	var err error
	if multi {
		_, err = s.Collection.RemoveAll(criteria)
	} else {
		err = s.Collection.Remove(criteria)
	}

	return err
}
