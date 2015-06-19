package generator

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/tools/go/types"
)

var findableTypes = map[string]bool{
	"string":    true,
	"int":       true,
	"int8":      true,
	"int16":     true,
	"int32":     true,
	"int64":     true,
	"uint":      true,
	"uint8":     true,
	"uint16":    true,
	"uint32":    true,
	"uint64":    true,
	"float32":   true,
	"float64":   true,
	"struct":    true,
	"bool":      true,
	"map":       true,
	"time.Time": true,
}

type Package struct {
	Name      string
	Models    []*Model
	Structs   []string
	Functions []string
}

func (p *Package) StructIsDefined(name string) bool {
	for _, n := range p.Structs {
		if name == n {
			return true
		}
	}

	return false
}

func (p *Package) FunctionIsDefined(name string) bool {
	for _, n := range p.Functions {
		if name == n {
			return true
		}
	}

	return false
}

const (
	StoreNamePattern     = "%sStore"
	QueryNamePattern     = "%sQuery"
	ResultSetNamePattern = "%sResultSet"
)

type Model struct {
	Name          string
	StoreName     string
	QueryName     string
	ResultSetName string

	Collection string
	Type       string
	Fields     []*Field
}

func NewModel(n string) *Model {
	return &Model{
		Name:          n,
		StoreName:     fmt.Sprintf(StoreNamePattern, n),
		QueryName:     fmt.Sprintf(QueryNamePattern, n),
		ResultSetName: fmt.Sprintf(ResultSetNamePattern, n),
		Type:          "struct",
		Fields:        make([]*Field, 0),
	}
}

func (m *Model) String() string {
	fields := make([]string, 0)
	for _, f := range m.Fields {
		fields = append(fields, "\t"+f.String()+"\n")
	}

	fieldsStr := strings.Join(fields, "")
	str := fmt.Sprintf("(Model '%s' [\n %s]", m.Name, fieldsStr)

	return str
}

func (m *Model) ValidFields() []*Field {
	fields := make([]*Field, 0)
	for _, f := range m.Fields {
		if f.Findable() {
			fields = append(fields, f)
		}
	}

	return fields
}

type Function struct {
	Name string
	Args string
}

func NewFunction() {
}

type Field struct {
	Name        string
	Type        string
	CheckedNode *types.Var
	Tag         reflect.StructTag
	Fields      []*Field
	Parent      *Field
	isMap       bool
}

func NewField(n, t string, tag reflect.StructTag) *Field {
	return &Field{
		Name:   n,
		Type:   t,
		Tag:    tag,
		Fields: make([]*Field, 0),
		isMap:  strings.HasPrefix(t, "map["),
	}
}

func (f *Field) SetFields(sf []*Field) {
	for _, field := range sf {
		f.AddField(field)
	}
}

func (f *Field) AddField(field *Field) {
	field.Parent = f
	f.Fields = append(f.Fields, field)
}

func (f *Field) GetPath() string {
	recursive := f
	path := make([]string, 0)
	for recursive != nil {
		if recursive.isMap {
			path = append(path, "[map]")
		}

		path = append(path, recursive.DbName())
		recursive = recursive.Parent
	}

	return strings.Join(reverseSliceStrings(path), ".")
}

func (f *Field) ContainsMap() bool {
	if !f.isMap && f.Parent != nil {
		return f.Parent.ContainsMap()
	}

	return f.isMap
}

func (f *Field) GetTagValue(key string) string {
	if f.Tag == "" {
		return ""
	}

	return f.Tag.Get(key)
}

func (f *Field) DbName() string {
	name := f.GetTagValue("bson")
	endFieldName := strings.Index(name, ",")
	if endFieldName != -1 {
		name = name[:endFieldName]
	}

	if name == "" {
		name = strings.ToLower(f.Name)
	}

	return name
}

func (f *Field) ValidFields() []*Field {
	fields := make([]*Field, 0)
	for _, f := range f.Fields {
		if f.Findable() {
			fields = append(fields, f)
		}
	}

	return fields
}

func (f *Field) FindableType() string {
	startType := strings.Index(f.Type, "]")
	if startType != -1 {
		return f.Type[startType+1:]
	}

	return f.Type
}

func (f *Field) Findable() bool {
	return findableTypes[f.FindableType()]
}

func (f *Field) String() string {
	return f.toString(0)
}

func (f *Field) toString(l int) string {
	if l > 3 {
		return "... more depth ..."
	}
	fields := make([]string, 0)
	for _, f := range f.Fields {
		fields = append(fields, f.toString(l+1))
	}

	fieldsStr := strings.Join(fields, ", ")

	return fmt.Sprintf("%s %s %s [%s]", f.Name, f.Type, f.Tag, fieldsStr)
}

func reverseSliceStrings(input []string) []string {
	if len(input) == 0 {
		return input
	}

	return append(reverseSliceStrings(input[1:]), input[0])
}
