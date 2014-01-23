package main

import (
	"flag"
	"fmt"
)

var langFlag string

func init() {
	flag.StringVar(&langFlag, "lang", "en", "ISO lang code")
}

func main() {
	flag.Parse()
	fmt.Println("hello " + langFlag)
}
