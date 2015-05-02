package storable

import (
	"errors"

	"github.com/tyba/storable/operators"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	NonNewDocumentErr  = errors.New("Cannot insert a non new document.")
	NewDocumentErr     = errors.New("Cannot updated a new document.")
	EmptyQueryInRawErr = errors.New("Empty queries are not allowed on raw ops.")
)

type Store struct {
	Collection string

	collection *mgo.Collection
}

func NewStore(db *mgo.Database, collection string) *Store {
	return &Store{
		Collection: collection,
		collection: db.C(collection),
	}
}

func (s *Store) Insert(doc DocumentBase) error {
	if !doc.IsNew() {
		return NonNewDocumentErr
	}

	doc.SetId(bson.NewObjectId())
	err := s.collection.Insert(doc)
	if err == nil {
		doc.SetIsNew(false)
	}

	return err
}

func (s *Store) Update(doc DocumentBase) error {
	if doc.IsNew() {
		return NewDocumentErr
	}

	q := NewBaseQuery()
	q.AddCriteria(operators.Eq(IdField, doc.GetId()))

	return s.collection.Update(q.GetCriteria(), doc)
}

func (s *Store) Save(doc DocumentBase) error {
	if doc.IsNew() {
		return s.Insert(doc)
	}

	return s.Update(doc)
}

func (s *Store) Delete(doc DocumentBase) error {
	q := NewBaseQuery()
	q.AddCriteria(operators.Eq(IdField, doc.GetId()))

	return s.collection.Remove(q.GetCriteria())
}

func (s *Store) Find(q Query) (*ResultSet, error) {
	mq := s.collection.Find(q.GetCriteria())

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

func (s *Store) RawUpdate(query Query, update interface{}) error {
	criteria := query.GetCriteria()
	if len(criteria) == 0 {
		return EmptyQueryInRawErr
	}

	return s.collection.Update(criteria, update)
}
