package Test

import (
	"testing"
	"fmt"
)

func TestReturn(t *testing.T) {
	a,b:=test()

	fmt.Printf("%d--%d\n",a,b)
	c,d:=test2()
	fmt.Printf("%d--%d",c,d)

}
func test() (int32,int32) {
	return 1,2
}

func test2() (a int32,b int32) {
	a=3
	b=4
	return
}