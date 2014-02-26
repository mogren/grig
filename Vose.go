package grig

import (
	"errors"
	"fmt"
	"math/rand"
)

type Vose struct {
	limit     int
	prob      []float64
	alias     []int64
	generator rand.Rand
}

func NewVose(prob []float64, generator rand.Rand) (v *Vose, err error) {
	v = new(Vose)
	v.generator = generator
	v.prob = prob
	v.limit = len(prob)
	if v.limit == 0 {
		return v, errors.New("Empty prob slice!")
	}
	sum := 0.0
	// Check for negative numbers
	for i, d := range prob {
		if d < 0.0 {
			return v, errors.New(fmt.Sprintf("%.6f", prob[i]) + " is below zero!")
		}
		sum += d
	}
	scale := float64(v.limit) / sum
	scaledProb := make([]float64, v.limit)
	for i, d := range prob {
		scaledProb[i] = d
	}
	fmt.Println(v.limit, generator.Intn(v.limit), scale)
	fmt.Println(scaledProb)
	return v, nil
}

func initVose(v Vose) {
	fmt.Println(v.GetLimit())
}

func (v Vose) GetLimit() int {
	return v.limit
}
