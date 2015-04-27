package storable

import (
	"strings"
)

type Field struct {
	bson string
	typ  string
}

func NewField(name, typ string) Field {
	return Field{name, typ}
}

func (f Field) Type() string {
	return f.typ
}

func (f Field) String() string {
	return f.bson
}

var (
	IdField = NewField("_id", "bson.ObjectId")
)

type Dir int

const (
	Asc  Dir = 1
	Desc Dir = -1
)

type Sort []FieldSort

type FieldSort struct {
	F Field
	D Dir
}

func (s Sort) String() string {
	var fields []string
	for _, fs := range s {
		f := ""
		if fs.D == Desc {
			f += "-"
		}

		f += fs.F.String()

		fields = append(fields, f)
	}

	return strings.Join(fields, ",")
}

func (s Sort) IsEmpty() bool {
	return len(s) == 0
}
