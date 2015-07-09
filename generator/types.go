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

	Collection  string
	Type        string
	Fields      []*Field
	CheckedNode *types.Named
	NewFunc     *types.Func
	Package     *types.Package

	Hooks      []Hook
	StoreHooks []Hook
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

func (m *Model) NewArgs() string {
	if m.NewFunc == nil {
		return ""
	}

	var ret []string
	sig := m.NewFunc.Type().(*types.Signature)

	for i := 0; i < sig.Params().Len(); i++ {
		param := sig.Params().At(i)
		typeName := types.TypeString(param.Type(), types.RelativeTo(m.Package))
		ret = append(ret, fmt.Sprintf("%v %v", param.Name(), typeName))
	}

	return strings.Join(ret, ", ")
}

func (m *Model) NewArgVars() string {
	if m.NewFunc == nil {
		return ""
	}

	var ret []string
	sig := m.NewFunc.Type().(*types.Signature)

	for i := 0; i < sig.Params().Len(); i++ {
		ret = append(ret, sig.Params().At(i).Name())
	}

	return strings.Join(ret, ", ")
}

func (m *Model) NewReturns() string {
	if m.NewFunc == nil {
		return "(doc *" + m.Name + ")"
	}

	var ret []string
	hasError := false
	sig := m.NewFunc.Type().(*types.Signature)

	for i := 0; i < sig.Results().Len(); i++ {
		res := sig.Results().At(i)
		typeName := types.TypeString(res.Type(), types.RelativeTo(m.Package))
		if isTypeOrPtrTo(res.Type(), m.CheckedNode) {
			ret = append(ret, "doc "+typeName)
		} else if isBuiltinError(res.Type()) && !hasError {
			ret = append(ret, "err "+typeName)
			hasError = true
		} else if res.Name() != "" {
			ret = append(ret, fmt.Sprintf("r%d %v", i, res.Name()))
		} else {
			ret = append(ret, fmt.Sprintf("r%d %v", i, typeName))
		}
	}

	return "(" + strings.Join(ret, ", ") + ")"
}

func (m *Model) NewRetVars() string {
	if m.NewFunc == nil {
		return "doc"
	}

	var ret []string
	hasError := false
	sig := m.NewFunc.Type().(*types.Signature)

	for i := 0; i < sig.Results().Len(); i++ {
		res := sig.Results().At(i)
		if isTypeOrPtrTo(res.Type(), m.CheckedNode) {
			ret = append(ret, "doc")
		} else if isBuiltinError(res.Type()) && !hasError {
			ret = append(ret, "err")
			hasError = true
		} else {
			ret = append(ret, fmt.Sprintf("r%d", i))
		}
	}

	return strings.Join(ret, ", ")
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
	Hooks       []Hook
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
	done := map[*Field]bool{}
	for recursive != nil {
		if recursive.isMap {
			path = append(path, "[map]")
		}

		path = append(path, recursive.DbName())
		recursive = recursive.Parent
		if done[recursive] {
			break
		}
		done[recursive] = true
	}

	return strings.Join(reverseSliceStrings(path), ".")
}

func (f *Field) ContainsMap() bool {
	return f.containsMap(map[*Field]bool{})
}

func (f *Field) containsMap(checked map[*Field]bool) bool {
	if checked[f] {
		return false
	}
	checked[f] = true

	if !f.isMap && f.Parent != nil {
		return f.Parent.containsMap(checked)
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

type Hook struct {
	Before bool
	Action HookAction
}

func (h Hook) MethodName() string {
	var ret string

	if h.Before {
		ret += "Before"
	} else {
		ret += "After"
	}

	ret += string(h.Action)

	return ret
}

type HookAction string

const (
	InsertHook HookAction = "Insert"
	UpdateHook HookAction = "Update"
	SaveHook   HookAction = "Save"
)

func isTypeOrPtrTo(ptr types.Type, named *types.Named) bool {
	switch ty := ptr.(type) {
	case *types.Pointer:
		if elem, ok := ty.Elem().(*types.Named); ok && elem == named {
			return true
		}
	case *types.Named:
		if ty == named {
			return true
		}
	}
	return false
}
