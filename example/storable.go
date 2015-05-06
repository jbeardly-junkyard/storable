package example

import (
	"github.com/tyba/storable"
	"github.com/tyba/storable/operators"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ProductStore struct {
	storable.Store
}

func NewProductStore(db *mgo.Database) *ProductStore {
	return &ProductStore{*storable.NewStore(db, "products")}
}

func (s *ProductStore) New() *Product {
	doc := &Product{}
	doc.SetIsNew(true)

	return doc
}

func (s *ProductStore) Query() *ProductQuery {
	return &ProductQuery{*storable.NewBaseQuery()}
}

func (s *ProductStore) Find(query *ProductQuery) (*ProductResultSet, error) {
	resultSet, err := s.Store.Find(query)
	if err != nil {
		return nil, err
	}

	return &ProductResultSet{*resultSet}, nil
}

func (s *ProductStore) FindOne(query *ProductQuery) (*Product, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

type ProductQuery struct {
	storable.BaseQuery
}

func (q *ProductQuery) FindById(id bson.ObjectId) {
	q.AddCriteria(operators.Eq(storable.IdField, id))
}

type ProductResultSet struct {
	storable.ResultSet
}

func (r *ProductResultSet) All() ([]*Product, error) {
	var result []*Product
	err := r.ResultSet.All(&result)

	return result, err
}

func (r *ProductResultSet) One() (*Product, error) {
	var result *Product
	_, err := r.ResultSet.One(&result)

	return result, err
}

func (r *ProductResultSet) Next() (*Product, error) {
	var result *Product
	_, err := r.ResultSet.Next(&result)

	return result, err
}

type schema struct {
	Product struct {
		Status    storable.Field
		CreatedAt storable.Field
		UpdatedAt storable.Field
		Name      storable.Field
		Price     struct {
			Amount   storable.Field
			Discount storable.Field
		}
		Discount storable.Field
		Url      storable.Field
		Tags     storable.Field
	}
}

var Schema = schema{
	Product: struct {
		Status    storable.Field
		CreatedAt storable.Field
		UpdatedAt storable.Field
		Name      storable.Field
		Price     struct {
			Amount   storable.Field
			Discount storable.Field
		}
		Discount storable.Field
		Url      storable.Field
		Tags     storable.Field
	}{
		Status:    storable.NewField("status", "int"),
		CreatedAt: storable.NewField("createdat", "time.Time"),
		UpdatedAt: storable.NewField("updatedat", "time.Time"),
		Name:      storable.NewField("name", "string"),
		Price: struct {
			Amount   storable.Field
			Discount storable.Field
		}{
			Amount:   storable.NewField("price.amount", "float64"),
			Discount: storable.NewField("price.discount", "float64"),
		},
		Discount: storable.NewField("discount", "float64"),
		Url:      storable.NewField("url", "string"),
		Tags:     storable.NewField("tags", "string"),
	},
}
