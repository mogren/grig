package vose

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

type Vose struct {
	limit     int
	prob      []float64
	alias     []int
	generator rand.Rand
}

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
	scale := float64(v.limit) / sum
	scaledProb := make([]float64, v.limit)
	for i, d := range prob {
		scaledProb[i] = d * scale
	}
	fmt.Println("Limit:", v.limit)
	fmt.Println("Scale:", scale)
	fmt.Println(scaledProb)
	initVose(v, scaledProb)
	fmt.Println(v.prob)
	fmt.Println(v.alias)
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
			nl += 1
		} else {
			small[ns] = i
			ns += 1
		}
	}
	// Build alias
	for ns != 0 && nl != 0 {
		ns -= 1
		nl -= 1
		j := small[ns]
		k := large[nl]
		v.prob[j] = scaledProb[j]
		v.alias[j] = k
		scaledProb[k] = (scaledProb[k] + scaledProb[j]) - 1.0
		if scaledProb[k] > 1.0 {
			large[nl] = k
			nl += 1
		} else {
			small[ns] = k
			ns += 1
		}
	}
	for ns != 0 {
		ns -= 1
		v.prob[small[ns]] = 1
	}
	for nl != 0 {
		nl -= 1
		v.prob[large[nl]] = 1
	}
}

func (v Vose) GetLimit() int {
	return v.limit
}

func (v Vose) Next() int {
	u := float64(v.limit) * v.generator.Float64()
	fl := math.Floor(u)
	j := int(fl)
	p := u - fl
	if p <= v.prob[j] {
		return j
	} else {
		return v.alias[j]
	}
}
