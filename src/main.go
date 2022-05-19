package main

import (
	"fmt"
	"www/src/router"
)

func main() {
	fmt.Println("Hello")
	rt := router.Init()

	rt.Logger.Fatal(rt.Start(":8080"))
}
