package main
import (
"fmt"
_ "github.com/mattn/go-sqlite3"
"math/rand"
"net/http"
_ "net/http/pprof"
"time"
)

const ServerVersion = 1
type TestObj struct {
	power int64
	id    int64
}
func main() {
	rand.Seed(time.Now().Unix())
	//http.HandleFunc("/", httphandler)
	for {
		time.Sleep(time.Minute)
		fmt.Printf("hh\n")
		go http.ListenAndServe(":11111", nil)
	}
	fmt.Println("===end=")
}
