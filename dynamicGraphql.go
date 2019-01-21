package dynamicGraphql

import (
	// "fmt"
	"github.com/graphql-go/graphql"
	"reflect"
)

func DynamicGraphql(name string, t reflect.Type) *graphql.Object {
	n := t.NumField()
	fields := make(map[string]*graphql.Field)
	for i := 0; i < n; i++ {
		f := t.Field(i)
		genField(f, fields)

	}
	objectConfig := graphql.ObjectConfig{
		Fields: fields,
	}
	if name != "" {
		objectConfig.Name = name
	}
	obj := graphql.NewObject(objectConfig)

	return obj

}

func genField(f reflect.StructField, fields map[string]*graphql.Field) {
	name := f.Tag.Get("json")
	if name == "" {
		return
	}
	t := f.Type
	fKind := f.Type.Kind()
	// name := f.Name
	// if name == "ID" {
	// 	name = "id"
	// } else {
	// 	name = ToLowerCamel(name)
	// }
	switch fKind {
	case reflect.Struct:
		if t.Name() != "Time" {
			obj := DynamicGraphql("", t)
			fields[name] = &graphql.Field{Type: obj}
		}
	case reflect.Ptr:
		s := t.Elem()
		if s.Name() != "Location" {
			obj := DynamicGraphql("", s)
			fields[name] = &graphql.Field{Type: obj}
		}
	default:
		fields[name] = &graphql.Field{Type: graphql.Int}
	}
}
