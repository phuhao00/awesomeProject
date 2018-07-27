package main

import (
	EE "../Server/routes/Definition"
)

const ServerVersion = 1

type TestObj struct {
	power int64
	id    int64
}

func main() {
	//rand.Seed(time.Now().Unix())
	//for {
	//	time.Sleep(time.Minute)
	//	fmt.Printf("hh\n")
	//	go http.ListenAndServe(":11111", nil)
	//}
	//fmt.Println("===end=")

	//===========================================
	server := EE.PSSM.NewServer("2", ":1111", "tcp", nil)
	server.Run()
}
