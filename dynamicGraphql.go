package dynamicGraphql

import (
	// "fmt"
	"github.com/graphql-go/graphql"
	"reflect"
	"strings"
)

var objects map[string]*graphql.Object

func DynamicGraphql(name string, t reflect.Type) *graphql.Object {
	if objects == nil {
		objects = make(map[string]*graphql.Object)
	}
	if objects[name] != nil {
		return objects[name]
	}
	n := t.NumField()
	fields := make(map[string]*graphql.Field)
	//--
	objectConfig := graphql.ObjectConfig{
		Fields: graphql.Fields(fields),
		Name:   name,
	}
	obj := graphql.NewObject(objectConfig)
	objects[name] = obj
	//--
	for i := 0; i < n; i++ {
		f := t.Field(i)
		//--base refine
		if f.Name == "Base" {
			tBase := f.Type
			nBase := tBase.NumField()
			for j := 0; j < nBase; j++ {
				fBase := tBase.Field(j)
				fNameBase, isJsonTagBase := getFieldName(fBase)
				if !isJsonTagBase {
					continue
				}
				fieldBase := genField(fBase.Type)
				if fieldBase != nil {
					obj.AddFieldConfig(fNameBase, fieldBase)
				}
			}
			continue
		}

		fName, isJsonTag := getFieldName(f)
		//--
		if !isJsonTag {
			continue
		}
		field := genField(f.Type)
		if field != nil {
			obj.AddFieldConfig(fName, field)

		}

	}
	return obj
}

func genObject(t reflect.Type) *graphql.Object {
	objectName := t.Name()
	if objects[objectName] != nil {
		return objects[objectName]
	}
	n := t.NumField()
	fields := make(map[string]*graphql.Field)
	//--
	objectConfig := graphql.ObjectConfig{
		Fields: graphql.Fields(fields),
		Name:   objectName,
	}
	obj := graphql.NewObject(objectConfig)
	objects[objectName] = obj
	//--
	for i := 0; i < n; i++ {
		f := t.Field(i)
		//--base refine
		if f.Name == "Base" {
			tBase := f.Type
			nBase := tBase.NumField()
			for j := 0; j < nBase; j++ {
				fBase := tBase.Field(j)
				fNameBase, isJsonTagBase := getFieldName(fBase)
				if !isJsonTagBase {
					continue
				}
				fieldBase := genField(fBase.Type)
				if fieldBase != nil {
					obj.AddFieldConfig(fNameBase, fieldBase)
				}
			}
			continue
		}

		fName, isJsonTag := getFieldName(f)
		//--
		if !isJsonTag {
			continue
		}
		field := genField(f.Type)
		if field != nil {
			obj.AddFieldConfig(fName, field)
		}
	}

	return obj

}
func genField(t reflect.Type) *graphql.Field {
	fKind := t.Kind()
	switch fKind {
	case reflect.Struct:
		if t.Name() != "Time" {
			obj := genObject(t)
			return &graphql.Field{Type: obj}
		}
	case reflect.Ptr:
		s := t.Elem()
		kind := s.Kind()
		if kind != reflect.Struct {
			return &graphql.Field{Type: getType(s, kind)}
		}
		if s.Name() != "Location" && s.Name() != "Time" {
			obj := genObject(s)
			return &graphql.Field{Type: obj}
		}
	case reflect.Slice:
		e := t.Elem()
		return genListField(e)

	default:
		return &graphql.Field{Type: getType(t, fKind)}
	}
	return nil
}
func genListField(t reflect.Type) *graphql.Field {
	fKind := t.Kind()
	switch fKind {
	case reflect.Struct:
		if t.Name() != "Time" {
			obj := genObject(t)
			return &graphql.Field{Type: graphql.NewList(obj)}
		}
	case reflect.Ptr:
		s := t.Elem()
		kind := s.Kind()
		if kind != reflect.Struct {
			tempType := getType(s, kind)
			return &graphql.Field{Type: graphql.NewList(tempType)}
		}
		if s.Name() != "Location" && s.Name() != "Time" {
			obj := genObject(s)
			return &graphql.Field{Type: graphql.NewList(obj)}
		}
	case reflect.Slice:
		e := t.Elem()
		return genListField(e)

	default:
		tempType := getType(t, fKind)
		return &graphql.Field{Type: graphql.NewList(tempType)}
	}
	return nil
}
func getFieldName(f reflect.StructField) (string, bool) {
	isJsonTag := true
	fName := f.Tag.Get("json")
	fName = strings.Split(fName, ",")[0]
	//--
	if fName == "" || fName == "-" {
		isJsonTag = false
		fName = f.Name
		if fName == "ID" {
			fName = "id"
		} else {
			fName = ToLowerCamel(fName)
		}
	}
	return fName, isJsonTag
}
func getType(t reflect.Type, kind reflect.Kind) graphql.Output {
	var graphqlType graphql.Output
	switch kind {
	case reflect.String:
		graphqlType = graphql.String
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint:
		graphqlType = graphql.Int
	case reflect.Float32, reflect.Float64:
		graphqlType = graphql.Float
	case reflect.Bool:
		graphqlType = graphql.Boolean
	case reflect.Slice:
		e := t.Elem()
		return getType(e, e.Kind())
	default:
		return nil
	}
	return graphqlType
}
