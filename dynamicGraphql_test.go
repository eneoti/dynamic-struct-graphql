package dynamicGraphql

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type testStruct struct {
	name     string `json:"name"`
	Bool     bool
	Integer  int
	Float    float64
	Time     time.Time
	LargeNum testSubStruct
	Other    *otherStruct
}

type otherStruct struct {
	name        string
	description string
	child       *childStruct
}
type childStruct struct {
	name    string
	address string
}
type testSubStruct struct {
	count  int
	IsZero bool
}

func TestGenGraphql(t *testing.T) {
	data := reflect.TypeOf(testStruct{})
	obj := DynamicGraphql("test", data)
	fmt.Println(obj != nil)
	// t.Logf("target: %+v", to)
	// t.Logf("to: %+v", target)

	// if expected := from.String; target.String != expected {
	// 	t.Errorf("want %s got %s", expected, target.String)
	// }
}
