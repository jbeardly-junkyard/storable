package example

import (
	"time"

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

func (s *ProductStore) New(name string, price Price, createdAt time.Time) (doc *Product, err error) {
	doc, err = newProduct(name, price, createdAt)
	doc.SetIsNew(true)
	doc.SetId(bson.NewObjectId())
	return
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

func (s *ProductStore) MustFind(query *ProductQuery) *ProductResultSet {
	resultSet, err := s.Find(query)
	if err != nil {
		panic(err)
	}

	return resultSet
}

func (s *ProductStore) FindOne(query *ProductQuery) (*Product, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

func (s *ProductStore) MustFindOne(query *ProductQuery) *Product {
	doc, err := s.FindOne(query)
	if err != nil {
		panic(err)
	}

	return doc
}

func (s *ProductStore) Insert(doc *Product) error {
	if err := doc.BeforeInsert(); err != nil {
		return storable.HookError{
			Hook:  "BeforeInsert",
			Field: "",
			Cause: err,
		}
	}
	if err := doc.BeforeSave(); err != nil {
		return storable.HookError{
			Hook:  "BeforeSave",
			Field: "",
			Cause: err,
		}
	}
	if err := doc.Status.BeforeInsert(); err != nil {
		return storable.HookError{
			Hook:  "BeforeInsert",
			Field: ".Status",
			Cause: err,
		}
	}

	err := s.Store.Insert(doc)
	if err != nil {
		return err
	}
	if err := doc.Status.AfterInsert(); err != nil {
		return storable.HookError{
			Hook:  "AfterInsert",
			Field: ".Status",
			Cause: err,
		}
	}

	return nil
}

func (s *ProductStore) Update(doc *Product) error {
	if err := doc.BeforeSave(); err != nil {
		return storable.HookError{
			Hook:  "BeforeSave",
			Field: "",
			Cause: err,
		}
	}

	err := s.Store.Update(doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductStore) Save(doc *Product) (updated bool, err error) {
	if err := doc.BeforeSave(); err != nil {
		return updated, storable.HookError{
			Hook:  "BeforeSave",
			Field: "",
			Cause: err,
		}
	}

	updated, err = s.Store.Save(doc)
	if err != nil {
		return false, err
	}

	if updated {

	} else {
		if err := doc.Status.AfterInsert(); err != nil {
			return updated, storable.HookError{
				Hook:  "AfterInsert",
				Field: ".Status",
				Cause: err,
			}
		}

	}

	return
}

type ProductQuery struct {
	storable.BaseQuery
}

func (q *ProductQuery) FindById(ids ...bson.ObjectId) {
	var vs []interface{}
	for _, id := range ids {
		vs = append(vs, id)
	}
	q.AddCriteria(operators.In(storable.IdField, vs...))
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
	err := r.ResultSet.One(&result)

	return result, err
}

func (r *ProductResultSet) Next() (*Product, error) {
	var result *Product
	_, err := r.ResultSet.Next(&result)

	return result, err
}

type schema struct {
	Product *schemaProduct
}

type schemaProduct struct {
	Status    storable.Field
	CreatedAt storable.Field
	UpdatedAt storable.Field
	Name      storable.Field
	Price     *schemaProductPrice
	Discount  storable.Field
	Url       storable.Field
	Tags      storable.Field
}

type schemaProductPrice struct {
	Amount   storable.Field
	Discount storable.Field
}

var Schema = schema{
	Product: &schemaProduct{
		Status:    storable.NewField("status", "int"),
		CreatedAt: storable.NewField("createdat", "time.Time"),
		UpdatedAt: storable.NewField("updatedat", "time.Time"),
		Name:      storable.NewField("name", "string"),
		Price: &schemaProductPrice{
			Amount:   storable.NewField("price.amount", "float64"),
			Discount: storable.NewField("price.discount", "float64"),
		},
		Discount: storable.NewField("discount", "float64"),
		Url:      storable.NewField("url", "string"),
		Tags:     storable.NewField("tags", "string"),
	},
}
