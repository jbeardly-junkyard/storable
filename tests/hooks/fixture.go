package hooks

import (
	"errors"
	"fmt"

	"github.com/tyba/storable"
)

//go:generate storable gen

var Log []string

type Recur struct {
	storable.Document `bson:",inline" collection:"recur"`
	Foo               string
	R                 *Other `bson:"-"`
	MoreThings        []Thing
	MyFailer          *Failer
	MyAfterFailer     *AfterFailer
	Things            map[string][]*Thing
}

func (r Recur) BeforeInsert() error {
	Log = append(Log, "Called BeforeInsert on Recur with Foo "+r.Foo)
	return nil
}

func (r *Recur) BeforeUpdate(s *RecurStore) error {
	Log = append(Log, "Called BeforeUpdate(s) on *Recur with Foo "+r.Foo)
	return nil
}

func (r *Recur) BeforeSave() error {
	Log = append(Log, "Called BeforeSave on *Recur with Foo "+r.Foo)
	return nil
}

func (s *RecurStore) BeforeSave() error {
	panic("Shouldn't have been called!")
	return nil
}

type Other struct {
	Name string
	R2   *Recur
}

func (r Other) AfterInsert() error {
	Log = append(Log, "Called AfterInsert on Other with Name "+r.Name)
	return nil
}

func (r Other) AfterUpdate() error {
	Log = append(Log, "Called AfterUpdate on Other with Name "+r.Name)
	return nil
}

func (r *Other) AfterSave() error {
	Log = append(Log, "Called AfterSave on *Other with Name "+r.Name)
	return nil
}

func (r *Other) BeforeSave() { // Bad signature.
	panic("Shouldn't have been called!")
}

type Thing struct {
	I int
}

func (t Thing) BeforeSave() error {
	Log = append(Log, fmt.Sprintf("Called BeforeSave on Thing %v", t.I))
	return nil
}

type Failer struct{}

func (f *Failer) BeforeInsert() error {
	Log = append(Log, "Called BeforeInsert on *Failer")
	return errors.New("I failed, sorry!")
}

func (f *Failer) BeforeSave() error {
	Log = append(Log, "Called BeforeSave on *Failer")
	return errors.New("I failed, sorry!")
}

type AfterFailer struct{}

func (f *AfterFailer) AfterSave() error {
	Log = append(Log, "Called AfterSave on *AfterFailer")
	return errors.New("I failed too late, sorry!")
}
