// The MIT License (MIT)
// Copyright (c) 2014 Claes Mogren
// http://opensource.org/licenses/MIT

// Code to generate a weighted random identity based on weighted input files
// http://www.keithschwarz.com/darts-dice-coins/
// http://web.eecs.utk.edu/~vose/Publications/random.pdf
package main

import (
	"./vose/"
	"bufio"
	"encoding/json"
	"encoding/xml"
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
var nrLoopsFlag int
var jsonFlag bool
var xmlFlag bool

// Rig holds a randomly generated identity
type Rig struct {
	Firstname    string `json:"firstname" xml:"firstname"`
	Lastname     string `json:"lastname" xml:"lastname"`
	Street       string `json:"street" xml:"street"`
	Streetnumber int    `json:"nr" xml:"nr"`
	Zipcode      int    `json:"zip" xml:"zip"`
	City         string `json:"city" xml:"city"`
}

// AsText prints the Rig as text
func (r Rig) AsText() string {
	str := fmt.Sprintln(r.Firstname, r.Lastname)
	if langFlag == "en_us" {
		str += fmt.Sprintln(r.Streetnumber, r.Street)
	} else {
		str += fmt.Sprintln(r.Street, r.Streetnumber)
	}
	str += fmt.Sprintln(r.Zipcode, r.City)
	return str
}

// AsJSON will output the Rig as JSON
func (r Rig) AsJSON() string {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}

// AsXML will output the Rig as XML
func (r Rig) AsXML() string {
	b, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}

// RigFile contains all the randomly generated identities
type RigFile struct {
	tot   float64
	texts [][]string
	vose  *vose.Vose
}

// RigDict is the dictionary for the given language
type RigDict struct {
	fnames, mnames, lnames, streets, zipcodes RigFile
}

func init() {
	flag.StringVar(&langFlag, "lang", "en_us", "Select ISO 639-1 language code, defaults to USA")
	flag.BoolVar(&listLangFlag, "l", false, "List available ISO language codes")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&jsonFlag, "j", false, "Print as JSON")
	flag.BoolVar(&xmlFlag, "x", false, "Print as XML")
	flag.IntVar(&nrLoopsFlag, "n", 1, "Number of identities to output")
}

func main() {
	flag.Parse()
	if listLangFlag {
		listLangs()
		os.Exit(0)
	}
	rand.Seed(time.Now().UnixNano())
	dict := loadData(langFlag)
	for i := 0; i < nrLoopsFlag; i++ {
		rig := getNext(dict)
		if jsonFlag {
			fmt.Println(rig.AsJSON())
		} else if xmlFlag {
			fmt.Println(rig.AsXML())
		} else {
			fmt.Println(rig.AsText())
		}
	}
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
	// Check for fnames, mnames, lnames, zipcodes and Streets
	srcFileNames := []string{"fnames.grig", "lnames.grig", "mnames.grig", "Streets.grig", "zipcodes.grig"}
	valid := true
	filename := ""
	for _, srcFile := range srcFileNames {
		filename = "./data/" + iso + "/" + srcFile
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
	dict.streets = loadFile(iso, "Streets.grig")
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
	rigFile.texts = make([][]string, 0)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	scanner := bufio.NewScanner(file)
	sum := 0.0
	var weights []float64
	var dataStr, str []string
	for scanner.Scan() {
		scanText := scanner.Text()
		if strings.HasPrefix(scanText, "#") {
			continue
		}
		dataStr = strings.Split(scanner.Text(), "\t")
		// string to float
		f, err := strconv.ParseFloat(dataStr[0], 64)
		if err != nil {
			// Ignore error
			fmt.Println(err)
		}
		sum += f
		str = dataStr[1:]
		weights = append(weights, f)
		rigFile.texts = append(rigFile.texts, str)
		if verbose {
			fmt.Println("P: ", f, str)
		}
	}
	rigFile.tot = sum
	rigFile.vose, err = vose.NewVose(weights, rand.New(rand.NewSource(time.Now().UnixNano())))
	if err != nil {
		fmt.Println("Vose:", err)
	}
	return rigFile
}

func getNext(dict RigDict) Rig {
	rig := Rig{}
	if rand.Intn(2) == 0 {
		rig.Firstname = dict.fnames.texts[dict.fnames.vose.Next()][0]
	} else {
		rig.Firstname = dict.mnames.texts[dict.mnames.vose.Next()][0]
	}
	rig.Lastname = dict.lnames.texts[dict.lnames.vose.Next()][0]
	rig.Street = dict.streets.texts[dict.streets.vose.Next()][0]
	rig.Streetnumber = rand.Intn(59) + 1
	zip := dict.zipcodes.texts[dict.zipcodes.vose.Next()]
	rig.City = zip[1]
	rig.Zipcode, _ = strconv.Atoi(zip[0])
	return rig
}
