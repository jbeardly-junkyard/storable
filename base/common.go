package base

import (
	"strings"
)

type Field string

func (f Field) String() string {
	return string(f)
}

const (
	IdField Field = "_id"
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
