package main

import (
	"fmt"
	"www/src/router"
)

func main() {
	fmt.Println("Wellcome to Golang")
	fmt.Println("Start to Programming golang")
	rt := router.Init()

	rt.Logger.Fatal(rt.Start(":8080"))
}
