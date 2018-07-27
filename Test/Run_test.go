package Test

import (
	"fmt"
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	var x float64 = 3.4
	p := reflect.ValueOf(&x) // Note: take the address of x.
	v := p.Elem()
	v.SetFloat(4.9)
	fmt.Println("type of p:", p.Type())
	fmt.Println("settability of p:", v.CanSet())
	fmt.Println(x)
}
