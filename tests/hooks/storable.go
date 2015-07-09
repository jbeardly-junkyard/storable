package hooks

import (
	"github.com/tyba/storable"
	"github.com/tyba/storable/operators"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type RecurStore struct {
	storable.Store
}

func NewRecurStore(db *mgo.Database) *RecurStore {
	return &RecurStore{*storable.NewStore(db, "recur")}
}

func (s *RecurStore) New() (doc *Recur) {
	doc = &Recur{}
	doc.SetIsNew(true)
	return
}

func (s *RecurStore) Query() *RecurQuery {
	return &RecurQuery{*storable.NewBaseQuery()}
}

func (s *RecurStore) Find(query *RecurQuery) (*RecurResultSet, error) {
	resultSet, err := s.Store.Find(query)
	if err != nil {
		return nil, err
	}

	return &RecurResultSet{*resultSet}, nil
}

func (s *RecurStore) MustFind(query *RecurQuery) *RecurResultSet {
	resultSet, err := s.Find(query)
	if err != nil {
		panic(err)
	}

	return resultSet
}

func (s *RecurStore) FindOne(query *RecurQuery) (*Recur, error) {
	resultSet, err := s.Find(query)
	if err != nil {
		return nil, err
	}

	return resultSet.One()
}

func (s *RecurStore) MustFindOne(query *RecurQuery) *Recur {
	doc, err := s.FindOne(query)
	if err != nil {
		panic(err)
	}

	return doc
}

func (s *RecurStore) Insert(doc *Recur) error {
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
	if doc.R != nil {
		// Loop: .R.R2.R2 *..Recur

	}
	for k0, _ := range doc.MoreThings {
		if err := doc.MoreThings[k0].BeforeSave(); err != nil {
			return storable.HookError{
				Hook:  "BeforeSave",
				Field: ".MoreThings[k0]",
				Cause: err,
			}
		}

	}
	if doc.MyFailer != nil {
		if err := doc.MyFailer.BeforeInsert(); err != nil {
			return storable.HookError{
				Hook:  "BeforeInsert",
				Field: ".MyFailer",
				Cause: err,
			}
		}
		if err := doc.MyFailer.BeforeSave(); err != nil {
			return storable.HookError{
				Hook:  "BeforeSave",
				Field: ".MyFailer",
				Cause: err,
			}
		}

	}
	for k0, _ := range doc.Things {
		for k1, _ := range doc.Things[k0] {
			if doc.Things[k0][k1] != nil {
				if err := doc.Things[k0][k1].BeforeSave(); err != nil {
					return storable.HookError{
						Hook:  "BeforeSave",
						Field: ".Things[k0][k1]",
						Cause: err,
					}
				}

			}

		}

	}

	err := s.Store.Insert(doc)
	if err != nil {
		return err
	}
	if doc.R != nil {
		if err := doc.R.AfterInsert(); err != nil {
			return storable.HookError{
				Hook:  "AfterInsert",
				Field: ".R",
				Cause: err,
			}
		}
		if err := doc.R.AfterSave(); err != nil {
			return storable.HookError{
				Hook:  "AfterSave",
				Field: ".R",
				Cause: err,
			}
		}
		// Loop: .R.R2.R2 *..Recur

	}
	if doc.MyAfterFailer != nil {
		if err := doc.MyAfterFailer.AfterSave(); err != nil {
			return storable.HookError{
				Hook:  "AfterSave",
				Field: ".MyAfterFailer",
				Cause: err,
			}
		}

	}

	return nil
}

func (s *RecurStore) Update(doc *Recur) error {
	if err := doc.BeforeSave(); err != nil {
		return storable.HookError{
			Hook:  "BeforeSave",
			Field: "",
			Cause: err,
		}
	}
	if doc.R != nil {
		// Loop: .R.R2.R2 *..Recur

	}
	for k0, _ := range doc.MoreThings {
		if err := doc.MoreThings[k0].BeforeSave(); err != nil {
			return storable.HookError{
				Hook:  "BeforeSave",
				Field: ".MoreThings[k0]",
				Cause: err,
			}
		}

	}
	if doc.MyFailer != nil {
		if err := doc.MyFailer.BeforeSave(); err != nil {
			return storable.HookError{
				Hook:  "BeforeSave",
				Field: ".MyFailer",
				Cause: err,
			}
		}

	}
	for k0, _ := range doc.Things {
		for k1, _ := range doc.Things[k0] {
			if doc.Things[k0][k1] != nil {
				if err := doc.Things[k0][k1].BeforeSave(); err != nil {
					return storable.HookError{
						Hook:  "BeforeSave",
						Field: ".Things[k0][k1]",
						Cause: err,
					}
				}

			}

		}

	}
	if err := doc.BeforeUpdate(s); err != nil {
		return storable.HookError{
			Hook:  "BeforeUpdate",
			Field: ".",
			Cause: err,
		}
	}

	err := s.Store.Update(doc)
	if err != nil {
		return err
	}
	if doc.R != nil {
		if err := doc.R.AfterUpdate(); err != nil {
			return storable.HookError{
				Hook:  "AfterUpdate",
				Field: ".R",
				Cause: err,
			}
		}
		if err := doc.R.AfterSave(); err != nil {
			return storable.HookError{
				Hook:  "AfterSave",
				Field: ".R",
				Cause: err,
			}
		}
		// Loop: .R.R2.R2 *..Recur

	}
	if doc.MyAfterFailer != nil {
		if err := doc.MyAfterFailer.AfterSave(); err != nil {
			return storable.HookError{
				Hook:  "AfterSave",
				Field: ".MyAfterFailer",
				Cause: err,
			}
		}

	}

	return nil
}

type RecurQuery struct {
	storable.BaseQuery
}

func (q *RecurQuery) FindById(id bson.ObjectId) {
	q.AddCriteria(operators.Eq(storable.IdField, id))
}

func (q *RecurQuery) FindByIds(ids ...bson.ObjectId) {
	q.AddCriteria(operators.In(storable.IdField, ids))
}

type RecurResultSet struct {
	storable.ResultSet
}

func (r *RecurResultSet) All() ([]*Recur, error) {
	var result []*Recur
	err := r.ResultSet.All(&result)

	return result, err
}

func (r *RecurResultSet) One() (*Recur, error) {
	var result *Recur
	_, err := r.ResultSet.One(&result)

	return result, err
}

func (r *RecurResultSet) Next() (*Recur, error) {
	var result *Recur
	_, err := r.ResultSet.Next(&result)

	return result, err
}

type schema struct {
	Recur *schemaRecur
}

type schemaRecur struct {
	Foo           storable.Field
	R             *schemaRecurR
	MoreThings    *schemaRecurMoreThings
	MyFailer      *schemaRecurMyFailer
	MyAfterFailer *schemaRecurMyAfterFailer
	Things        *schemaRecurThings
}

type schemaRecurR struct {
	Name storable.Field
	R2   *schemaRecurRR2
}

type schemaRecurMoreThings struct {
	I storable.Map
}

type schemaRecurMyFailer struct {
}

type schemaRecurMyAfterFailer struct {
}

type schemaRecurThings struct {
	I storable.Map
}

type schemaRecurRR2 struct {
	Foo           storable.Field
	R             *schemaRecurR
	MoreThings    *schemaRecurMoreThings
	MyFailer      *schemaRecurMyFailer
	MyAfterFailer *schemaRecurMyAfterFailer
	Things        *schemaRecurThings
}

var Schema = schema{
	Recur: &schemaRecur{
		Foo: storable.NewField("-.r2.foo", "string"),
		R: &schemaRecurR{
			Name: storable.NewField("r2.-.name", "string"),
			R2: &schemaRecurRR2{
				Foo: storable.NewField("-.r2.foo", "string"),
				R:   nil,
				MoreThings: &schemaRecurMoreThings{
					I: storable.NewMap("-.r2.things.[map].i", "int"),
				},
				MyFailer:      &schemaRecurMyFailer{},
				MyAfterFailer: &schemaRecurMyAfterFailer{},
				Things: &schemaRecurThings{
					I: storable.NewMap("-.r2.things.[map].i", "int"),
				},
			},
		},
		MoreThings:    nil,
		MyFailer:      nil,
		MyAfterFailer: nil,
		Things:        nil,
	},
}
