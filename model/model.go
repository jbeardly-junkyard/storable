package model

import (
	"fmt"
	"reflect"
	"strings"
)

type Model struct {
	Name       string
	Collection string
	Fields     []*Field
}

func NewModel(name string) *Model {
	return &Model{
		Name:   name,
		Fields: make([]*Field, 0),
	}
}

func (m *Model) String() string {
	fields := make([]string, 0)
	for _, f := range m.Fields {
		fields = append(fields, f.String())
	}

	fieldsStr := strings.Join(fields, ", ")
	str := fmt.Sprintf("(Model '%s' [ %s ]", m.Name, fieldsStr)

	return str
}

type Field struct {
	Name string
	Type string
	Tag  *reflect.StructTag
}

func (f *Field) String() string {
	return fmt.Sprintf("%s %s %s", f.Name, f.Type, f.Tag)
}
func (f *Field) GetTagValue(key string) string {
	if f.Tag == nil {
		return ""
	}
	return f.Tag.Get(key)
}

func (f *Field) DbName() string {
	name := f.GetTagValue("bson")
	if name == "" {
		name = strings.ToLower(f.Name)
	}

	return name
}

func NewModel(name string) *Model {
	return &Model{
		Name:   name,
		Fields: make([]Field, 0),
	}
}
