// The MIT License (MIT)
// Copyright (c) 2014 Claes Mogren
// http://opensource.org/licenses/MIT

// Code to generate a weighted random identity based on weighted input files
// http://www.keithschwarz.com/darts-dice-coins/
package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/mogren/grig/vose"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	langFlag     string
	listLangFlag bool
	verbose      bool
	nrLoopsFlag  int
	jsonFlag     bool
	xmlFlag      bool
	// Global RNG with mutex protection
	globalRNG = struct {
		sync.Mutex
		rng *rand.Rand
	}{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
)

// Rig holds a randomly generated identity
type Rig struct {
	Firstname    string `json:"firstname" xml:"firstname"`
	Lastname     string `json:"lastname" xml:"lastname"`
	Street       string `json:"street" xml:"street"`
	StreetNumber int    `json:"nr" xml:"nr"`
	Zipcode      int    `json:"zip" xml:"zip"`
	City         string `json:"city" xml:"city"`
}

// AsText prints the Rig as text
func (r Rig) AsText() string {
	var sb strings.Builder
	sb.WriteString(r.Firstname + " " + r.Lastname + "\n")
	if langFlag == "en_us" {
		fmt.Fprintf(&sb, "%d %s\n", r.StreetNumber, r.Street)
	} else {
		fmt.Fprintf(&sb, "%s %d\n", r.Street, r.StreetNumber)
	}
	fmt.Fprintf(&sb, "%d %s\n", r.Zipcode, r.City)
	return sb.String()
}

// AsJSON will output the Rig as JSON
func (r Rig) AsJSON(isLast bool) string {
	var sb strings.Builder
	b, err := json.Marshal(r) // Use Marshal instead of MarshalIndent for performance
	if err != nil {
		fmt.Println("error:", err)
	}
	sb.Write(b)
	if !isLast {
		sb.WriteString(",")
	}
	return sb.String()
}

// AsXML will output the Rig as XML
func (r Rig) AsXML() string {
	var sb strings.Builder
	b, err := xml.Marshal(r) // Use Marshal instead of MarshalIndent for performance
	if err != nil {
		fmt.Println("error:", err)
	}
	sb.Write(b)
	return sb.String()
}

// RigFile contains all the randomly generated identities
type RigFile struct {
	total float64
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
	
	dict := loadData(langFlag)
	
	// Pre-allocate identities for better performance when generating multiple
	identities := make([]Rig, 0, nrLoopsFlag)
	
	// Generate all identities
	for i := 0; i < nrLoopsFlag; i++ {
		identities = append(identities, getNext(dict))
	}
	
	// Output in requested format
	if jsonFlag {
		if nrLoopsFlag > 1 {
			fmt.Println("\"list\": [")
		}
		
		for i, rig := range identities {
			fmt.Println(rig.AsJSON(i == (nrLoopsFlag - 1)))
		}
		
		if nrLoopsFlag > 1 {
			fmt.Println("]")
		}
	} else if xmlFlag {
		for _, rig := range identities {
			fmt.Println(rig.AsXML())
		}
	} else {
		for _, rig := range identities {
			fmt.Println(rig.AsText())
		}
	}
}

func listLangs() {
	// for all dirs in data
	files, _ := os.ReadDir("./data/")
	for _, f := range files {
		if f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			if validateDir(f.Name()) {
				fmt.Println(f.Name())
			}
		}
	}
}

func validateDir(iso string) bool {
	// Check for fnames, mnames, lnames, zipcodes and Streets
	srcFileNames := []string{"fnames.grig", "lnames.grig", "mnames.grig", "streets.grig", "zipcodes.grig"}
	valid := true
	
	for _, srcFile := range srcFileNames {
		filePath := "./data/" + iso + "/" + srcFile
		_, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				if verbose {
					fmt.Println("Data file", srcFile, "missing for", iso)
				}
				valid = false
			} else {
				// Handle other errors (permissions, etc.)
				if verbose {
					fmt.Printf("Error accessing %s: %v\n", filePath, err)
				}
				valid = false
			}
		}
	}
	return valid
}

func loadData(iso string) RigDict {
	dict := RigDict{}
	var wg sync.WaitGroup
	wg.Add(5)
	
	// Load files in parallel
	go func() { dict.fnames = loadFile(iso, "fnames.grig"); wg.Done() }()
	go func() { dict.mnames = loadFile(iso, "mnames.grig"); wg.Done() }()
	go func() { dict.lnames = loadFile(iso, "lnames.grig"); wg.Done() }()
	go func() { dict.streets = loadFile(iso, "streets.grig"); wg.Done() }()
	go func() { dict.zipcodes = loadFile(iso, "zipcodes.grig"); wg.Done() }()
	
	wg.Wait()
	
	if verbose {
		fmt.Printf("fname total: %.f\n", dict.fnames.total)
		fmt.Printf("mname total: %.f\n", dict.mnames.total)
		fmt.Printf("lname total: %.f\n", dict.lnames.total)
		fmt.Printf("streets total: %.f\n", dict.streets.total)
		fmt.Printf("zipcodes total: %.f\n", dict.zipcodes.total)
	}
	return dict
}

func loadFile(iso string, srcFile string) RigFile {
	// Initial capacity estimates based on typical file sizes
	const initialCapacity = 1000
	
	file, err := os.Open("data/" + iso + "/" + srcFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
			os.Exit(1)
		}
	}(file)
	
	rigFile := RigFile{}
	rigFile.texts = make([][]string, 0, initialCapacity)
	scanner := bufio.NewScanner(file)
	sum := 0.0
	weights := make([]float64, 0, initialCapacity)
	
	for scanner.Scan() {
		scanText := scanner.Text()
		if strings.HasPrefix(scanText, "#") {
			continue
		}
		dataStr := strings.Split(scanText, "\t")
		// string to float
		f, err := strconv.ParseFloat(dataStr[0], 64)
		if err != nil {
			fmt.Printf("Error parsing float in %s: %v\n", srcFile, err)
			continue
		}
		sum += f
		weights = append(weights, f)
		rigFile.texts = append(rigFile.texts, dataStr[1:])
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", srcFile, err)
	}
	
	rigFile.total = sum
	
	// Use the global RNG with mutex protection
	globalRNG.Lock()
	rigFile.vose, err = vose.NewVose(weights, globalRNG.rng)
	globalRNG.Unlock()
	
	if err != nil {
		fmt.Printf("Vose error with %s: %v\n", srcFile, err)
		panic(1)
	}
	return rigFile
}

func getNext(dict RigDict) Rig {
	rig := Rig{}
	
	// Use bit mask for coin flip (faster than rand.Intn(2))
	globalRNG.Lock()
	coinFlip := globalRNG.rng.Uint32() & 1
	globalRNG.Unlock()
	
	if coinFlip == 0 {
		rig.Firstname = dict.fnames.texts[dict.fnames.vose.Next()][0]
	} else {
		rig.Firstname = dict.mnames.texts[dict.mnames.vose.Next()][0]
	}
	
	rig.Lastname = dict.lnames.texts[dict.lnames.vose.Next()][0]
	rig.Street = dict.streets.texts[dict.streets.vose.Next()][0]
	
	// Use global RNG with mutex protection
	globalRNG.Lock()
	rig.StreetNumber = globalRNG.rng.Intn(150) + 1
	globalRNG.Unlock()
	
	zip := dict.zipcodes.texts[dict.zipcodes.vose.Next()]
	rig.City = zip[1]
	zipcode, err := strconv.Atoi(zip[0])
	if err != nil {
		// Fallback to default value if conversion fails
		zipcode = 10000
	}
	rig.Zipcode = zipcode
	return rig
}
