package main

import (
	"flag"
	"fmt"
  "os"
  "io/ioutil"
)

var langFlag string
var listLangFlag bool

func init() {
	flag.StringVar(&langFlag, "lang", "en", "Select ISO 639-1 language code")
	flag.BoolVar(&listLangFlag, "l", false, "List available ISO language codes")
}

func main() {
	flag.Parse()
  if (listLangFlag) {
    fmt.Println("list: ")
    return
  }
	fmt.Println("hello " + langFlag)
  loadData(langFlag)
}

func loadData(iso string) {
  file, err := ioutil.ReadFile("data/" + langFlag + "/fnames.grig")
  if err != nil {
    fmt.Println(err)
    os.Exit(0)
  }
  str := string(file)
  fmt.Println(str)
}

