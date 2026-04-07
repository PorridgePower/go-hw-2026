package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	res := reverse.String("Hello, OTUS!")
	fmt.Println(res)
}
