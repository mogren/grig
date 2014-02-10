// The MIT License (MIT)
// Copyright (c) 2014 Claes Mogren
// http://opensource.org/licenses/MIT

// Code to generate a weighted random identity based on weighted input files
// http://www.keithschwarz.com/darts-dice-coins/
// http://web.eecs.utk.edu/~vose/Publications/random.pdf
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var langFlag string
var listLangFlag bool
var verbose bool

type RigFile struct {
	tot     float64
	weights []float64
	texts   [][]string
	// Vose
}

type RigDict struct {
	fnames, mnames, lnames, streets, zipcodes RigFile
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
	rand.Seed(time.Now().UnixNano())
	dict := loadData(langFlag)
	printNext(dict)
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
	dict := RigDict{}
	dict.fnames = loadFile(iso, "fnames.grig")
	dict.mnames = loadFile(iso, "mnames.grig")
	dict.lnames = loadFile(iso, "lnames.grig")
	dict.streets = loadFile(iso, "streets.grig")
	dict.zipcodes = loadFile(iso, "zipcodes.grig")
	if verbose {
		fmt.Println("fname tot:", dict.fnames.tot, dict.fnames.texts[0])
		fmt.Println("mname tot:", dict.mnames.tot, dict.mnames.texts[0])
		fmt.Println("lname tot:", dict.lnames.tot, dict.lnames.texts[0])
	}
	return dict
}

func loadFile(iso string, srcFile string) RigFile {
	file, err := os.Open("data/" + iso + "/" + srcFile)
	rigFile := RigFile{}
	rigFile.weights = make([]float64, 0)
	rigFile.texts = make([][]string, 0)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	scanner := bufio.NewScanner(file)
	sum := 0.0
	for scanner.Scan() {
		dataStr := strings.Split(scanner.Text(), "\t")
		// string to float
		f, err := strconv.ParseFloat(dataStr[0], 64)
		if err != nil {
			// Ignore error
			fmt.Println(err)
		}
		sum += f
		str := dataStr[1:]
		rigFile.weights = append(rigFile.weights, f)
		rigFile.texts = append(rigFile.texts, str)
		if verbose {
			fmt.Println("P: ", f, str)
		}
	}
	rigFile.tot = sum
	return rigFile
}

func printNext(dict RigDict) {
	if rand.Intn(2) == 0 {
		fmt.Print(dict.fnames.texts[rand.Intn(len(dict.fnames.texts))])
	} else {
		fmt.Print(dict.mnames.texts[rand.Intn(len(dict.mnames.texts))])
	}
	fmt.Println(dict.lnames.texts[rand.Intn(len(dict.lnames.texts))])
}
