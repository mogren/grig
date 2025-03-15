// The MIT License (MIT)
// Copyright (c) 2014 Claes Mogren
// http://opensource.org/licenses/MIT

// Code to generate a weighted random identity based on weighted input files
// http://www.keithschwarz.com/darts-dice-coins/

// Package vose is a weighed random generator
package vose

import (
	"errors"
	"fmt"
	"math/rand"
)

// Vose is holding the weighted generator
type Vose struct {
	limit     int
	prob      []float64
	alias     []int
	generator rand.Rand
}

// NewVose sets up a new Vose struct
func NewVose(prob []float64, generator *rand.Rand) (v *Vose, err error) {
	if len(prob) == 0 {
		return nil, errors.New("empty probability slice")
	}
	
	probLen := len(prob)
	v = &Vose{
		limit: probLen,
		alias: make([]int, probLen),
		prob:  make([]float64, probLen),
		generator: *generator,
	}
	
	// Check for negative numbers and get sum (in one pass)
	sum := 0.0
	for i, d := range prob {
		if d < 0.0 {
			return nil, fmt.Errorf("%.6f is below zero!", prob[i])
		}
		sum += d
	}
	
	// Normalize weights (with preallocated slice)
	scale := float64(v.limit) / sum
	scaledProb := make([]float64, probLen)
	for i, d := range prob {
		scaledProb[i] = d * scale
	}
	
	v.init(scaledProb)
	return v, nil
}

func (v *Vose) init(scaledProb []float64) {
	// Preallocate small and large arrays
	small := make([]int, v.limit)
	large := make([]int, v.limit)
	ns := 0
	nl := 0
	
	// Partition (more cache-friendly)
	for i, sd := range scaledProb {
		if sd > 1.0 {
			large[nl] = i
			nl++
		} else {
			small[ns] = i
			ns++
		}
	}
	
	// Build alias (optimized loop)
	for ns > 0 && nl > 0 {
		ns--
		nl--
		j := small[ns]
		k := large[nl]
		v.prob[j] = scaledProb[j]
		v.alias[j] = k
		
		// Update scaledProb[k] in one operation
		scaledProb[k] = scaledProb[k] + scaledProb[j] - 1.0
		
		if scaledProb[k] < 1.0 {
			small[ns] = k
			ns++
		} else {
			large[nl] = k
			nl++
		}
	}
	
	// Finalize remaining entries
	for ns > 0 {
		ns--
		v.prob[small[ns]] = 1
	}
	for nl > 0 {
		nl--
		v.prob[large[nl]] = 1
	}
}

// Next will get the next weighted random number from the Vose struct
func (v Vose) Next() int {
	// Die roll
	j := v.generator.Intn(v.limit)
	
	// Weighted coin toss
	if v.generator.Float64() < v.prob[j] {
		return j
	}
	return v.alias[j]
}
