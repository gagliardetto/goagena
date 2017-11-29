package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"unicode"

	. "github.com/dave/jennifer/jen"
	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/goagena/scanner"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

func main() {
	var pkg string
	flag.StringVar(&pkg, "pkg", "", "package you want to scan and convert to goa types")
	flag.Parse()

	sc, err := scanner.New(pkg)
	if err != nil {
		panic(err)
	}

	pks, err := sc.Scan()
	if err != nil {
		panic(err)
	}

	//spew.Dump(pks)

	f := NewDesignFile("design")

	for _, pk := range pks {
		//fmt.Println("------")
		//fmt.Println(pk.Aliases)
		//fmt.Println(pk.Path)
		//fmt.Println(pk.Name)
		//spew.Dump(pk.Structs)
		for _, str := range pk.Structs {

			spew.Dump(str.Doc)
			f.AddSimpleTypeFromStruct(str)

			// Check whether the struct has a
			// Validate() error
			// method.
			if ok, err := scanner.IsValidatable(str.Type); err == nil && ok {
				//fmt.Println(str.Name, "is validatable!")
			}

		}
	}
	fmt.Println(f)
}

func (d *DesignFile) String() string {
	result := fmt.Sprintf("%#v", d.ff)

	goaImports := `import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)`

	result = strings.Trim(result, "package "+d.packageName)

	result = fmt.Sprintf("package %s\n\n%s\n%s", d.packageName, goaImports, result)

	return result
}

func (d *DesignFile) AddSimpleTypeFromStruct(str *scanner.Struct) {
	var typeAttributes []Code
	if IsMediaType(str.Doc) {
		typeAttributes = newGoaMediaType(str.Name)
	} else {
		typeAttributes = newSimpleGoaType(str.Name)
	}

	attributes := make([]Code, 0)

	requiredAttributes := make([]string, 0)
	for _, field := range str.Fields {

		if !field.Type.IsPtr() {
			requiredAttributes = append(requiredAttributes, ExtractFieldName(field))
		}

		attribute := FieldToAttribute(field)
		attributes = append(attributes, attribute)

	}

	// Required attributes
	if len(requiredAttributes) > 0 {
		requiredParams := make([]Code, 0)
		for _, requiredFieldName := range requiredAttributes {
			requiredParams = append(requiredParams, Lit(requiredFieldName))
			// TODO: add newline that does not add a comma
		}
		required := Id("Required").Call(
			requiredParams...,
		)

		attributes = append(attributes, Line(), required)
	}

	fff := NewAttributesBlock(attributes)

	typeAttributes = append(typeAttributes, fff)

	if IsMediaType(str.Doc) {
		d.ff.Var().Id(str.Name).Op("=").Id("MediaType").Call(
			typeAttributes...,
		)
	} else {
		d.ff.Var().Id(str.Name).Op("=").Id("Type").Call(
			typeAttributes...,
		)
	}

}

func NewAttributesBlock(attributes []Code) *Statement {
	return Func().Params().Block(
		attributes...,
	)
}

func ExtractFieldName(field *scanner.Field) string {
	return toLowerSnakeCase(field.Name)
}

func FieldToAttribute(field *scanner.Field) *Statement {
	attributeParams := make([]Code, 0)

	// attribute name
	attributeNameString := ExtractFieldName(field)
	{
		attributeName := Lit(attributeNameString)
		attributeParams = append(attributeParams, attributeName)
	}

	// attribute type
	{
		var attributeType Code

		// if the type is repeated, it mean that the field is array or slice
		if field.Type.IsRepeated() {
			attributeType = Id("ArrayOf").Params(Id(typeToGoaType(field.Type.UnqualifiedName())))
		} else {
			attributeType = Id(typeToGoaType(field.Type.UnqualifiedName()))
		}

		attributeParams = append(attributeParams, attributeType)
	}

	// DEBUG:
	//fmt.Println(field.Name,
	//	"ptr:", field.Type.IsPtr(),
	//	"arr:", field.Type.IsRepeated(),
	//	"str:", field.Type.IsStruct(),
	//)

	if len(field.Doc) > 0 {
		attributeDescription := Lit(strings.Join(field.Doc, ""))
		attributeParams = append(attributeParams, attributeDescription)
	}

	attribute := Id("Attribute").Call(
		attributeParams...,
	)
	if field.Type.IsPtr() {
		attribute.Comment("Is optional")
	}

	return attribute
}

type DesignFile struct {
	packageName string
	ff          *File
}

func NewDesignFile(packageName string) *DesignFile {
	return &DesignFile{
		ff:          NewFile(packageName),
		packageName: packageName,
	}
}

func newSimpleGoaType(name string) []Code {
	typeAttributes := make([]Code, 0)
	typeAttributes = append(typeAttributes, Lit(name))
	return typeAttributes
}

func newGoaMediaType(name string) []Code {
	typeAttributes := make([]Code, 0)
	name = fmt.Sprintf("application/vnd.%s+json", toLowerSnakeCase(name))
	typeAttributes = append(typeAttributes, Lit(name))
	return typeAttributes
}

func typeToGoaType(t string) string {
	switch t {
	case "string":
		return "String"

	case "int", "int32", "int64":
		return "Integer"

	case "UUID":
		return "UUID"

	case "bool":
		return "Boolean"

	case "float", "float16", "float32", "float64":
		return "Number"

	case "Time":
		return "DateTime"

	default:
		return generator.CamelCase(toLowerSnakeCase(t))
	}
}

func toLowerSnakeCase(s string) string {
	var buf bytes.Buffer
	var lastWasUpper bool
	for i, r := range s {
		if unicode.IsUpper(r) && i != 0 && !lastWasUpper {
			buf.WriteRune('_')
		}
		lastWasUpper = unicode.IsUpper(r)
		buf.WriteRune(unicode.ToLower(r))
	}
	return buf.String()
}

func toUpperSnakeCase(s string) string {
	return strings.ToUpper(toLowerSnakeCase(s))
}

// IsMediaType returns true in case the user specified a comment on the
// struct to convert it to a MediaType.
func IsMediaType(doc []string) bool {
	for _, v := range doc {
		if strings.HasPrefix(v, "goagena:mediatype") {
			return true
		}
	}
	return false
}
