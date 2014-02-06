package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var langFlag string
var listLangFlag bool

type RigDict struct {
	fnames, mnames, lnames, streets, zipcodes map[string]int
}

func init() {
	flag.StringVar(&langFlag, "lang", "en", "Select ISO 639-1 language code")
	flag.BoolVar(&listLangFlag, "l", false, "List available ISO language codes")
}

func main() {
	flag.Parse()
	if listLangFlag {
		listLangs()
		return
	}
	loadData(langFlag)
}

func listLangs() {
	// for all dirs in data
	files, _ := ioutil.ReadDir("./data/")
	for _, f := range files {
		if f.IsDir() && !strings.HasPrefix(f.Name(), ",") {
			fmt.Println(f.Name())
		}
	}
}

func validateDir(iso string) {
	// Check for fnames, mnames, lnames, zipcodes and streets
}

func loadData(iso string) {
	file, err := os.Open("data/" + langFlag + "/fnames.grig")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
