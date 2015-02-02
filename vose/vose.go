// Package vose is a weighed random generator

// The MIT License (MIT)
// Copyright (c) 2014 Claes Mogren
// http://opensource.org/licenses/MIT

// Code to generate a weighted random identity based on weighted input files
// http://www.keithschwarz.com/darts-dice-coins/
// http://web.eecs.utk.edu/~vose/Publications/random.pdf
package vose

import (
	"errors"
	"fmt"
	"math/rand"
)

// Struct holding the weighted generator
type Vose struct {
	limit     int
	prob      []float64
	alias     []int
	generator rand.Rand
}

// Set up a new Vose struct
func NewVose(prob []float64, generator *rand.Rand) (v *Vose, err error) {
	if len(prob) == 0 {
		return nil, errors.New("Empty prob slice!")
	}
	v = new(Vose)
	v.generator = *generator
	v.limit = len(prob)
	v.alias = make([]int, v.limit)
	v.prob = make([]float64, v.limit)
	sum := 0.0
	// Check for negative numbers and get sum
	for i, d := range prob {
		if d < 0.0 {
			return nil, errors.New(fmt.Sprintf("%.6f", prob[i]) + " is below zero!")
		}
		sum += d
	}
	// Normalize weights
	scale := float64(v.limit) / sum
	scaledProb := make([]float64, v.limit)
	for i, d := range prob {
		scaledProb[i] = d * scale
	}
	initVose(v, scaledProb)
	return v, nil
}

func initVose(v *Vose, scaledProb []float64) {
	small := make([]int, v.limit)
	large := make([]int, v.limit)
	ns := 0
	nl := 0
	// Partition
	for i, sd := range scaledProb {
		if sd > 1.0 {
			large[nl] = i
			nl++
		} else {
			small[ns] = i
			ns++
		}
	}
	// Build alias
	for ns != 0 && nl != 0 {
		ns--
		nl--
		j := small[ns]
		k := large[nl]
		v.prob[j] = scaledProb[j]
		v.alias[j] = k
		scaledProb[k] = (scaledProb[k] + scaledProb[j]) - 1.0
		if scaledProb[k] < 1.0 {
			small[ns] = k
			ns++
		} else {
			large[nl] = k
			nl++
		}
	}
	for ns != 0 {
		ns--
		v.prob[small[ns]] = 1
	}
	for nl != 0 {
		nl--
		v.prob[large[nl]] = 1
	}
}

// Get the next weighted radom numver from the Vose
func (v Vose) Next() int {
	u := float64(v.limit) * v.generator.Float64()
	j := int(u)
	p := u - float64(j)
	if p <= v.prob[j] {
		return j
	}
	return v.alias[j]
}
