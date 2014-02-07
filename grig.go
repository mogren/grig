// The MIT License (MIT)
// Copyright (c) 2014 Claes Mogren
// http://opensource.org/licenses/MIT

// Code to generate a weighted random identity based on weighted input files
// http://www.keithschwarz.com/darts-dice-coins/
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
var verbose bool

type RigLine struct {
	weight float64
	text   []string
}

type RigFile struct {
	tot      float64
	probList []float64
	texts    [][]string
	// Vose
}

type RigDict struct {
	fnames, mnames, lnames, streets, zipcodes []string
}

func init() {
	flag.StringVar(&langFlag, "lang", "en", "Select ISO 639-1 language code")
	flag.BoolVar(&listLangFlag, "l", false, "List available ISO language codes")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
}

func main() {
	flag.Parse()
	if listLangFlag {
		listLangs()
		os.Exit(0)
	}
	loadData(langFlag)
}

func listLangs() {
	// for all dirs in data
	files, _ := ioutil.ReadDir("./data/")
	for _, f := range files {
		if f.IsDir() && !strings.HasPrefix(f.Name(), ",") {
			if validateDir(f.Name()) {
				fmt.Println(f.Name())
			}
		}
	}
}

func validateDir(iso string) bool {
	// Check for fnames, mnames, lnames, zipcodes and streets
	srcFileNames := []string{"fnames.grig", "lnames.grig", "mnames.grig", "streets.grig", "zipcodes.grig"}
	valid := true
	for _, srcFile := range srcFileNames {
		filename := "./data/" + iso + "/" + srcFile
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			if verbose {
				fmt.Println("Data file", srcFile, "missing for", iso)
			}
			valid = false
		}
	}
	return valid
}

func loadData(iso string) RigDict {
	//{"fnames.grig", "lnames.grig", "mnames.grig", "streets.grig", "zipcodes.grig"}
	dict := RigDict{}
	dict.fnames = loadFile("fname.grig")
	return dict
}

func loadFile(srcFile string) map[string]int {
	file, err := os.Open("data/" + iso + "/" + srcFile)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dataStr := scanner.Text()
	}
}
