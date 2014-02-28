package vose

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	prob := make([]float64, 3)
	prob[0] = 0.1
	prob[1] = 0.2
	prob[2] = 0.7
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	v, err := NewVose(prob, generator)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v.Next())
		fmt.Println(v.Next())
		fmt.Println(v.Next())
		fmt.Println(v.Next())
	}
	fmt.Println("ok")
}

func TestScale(t *testing.T) {
	prob := make([]float64, 4)
	prob[0] = 11.0
	prob[1] = 100.2
	prob[2] = 120.0
	prob[3] = 12120.0
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	v, err := NewVose(prob, generator)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v.Next())
		fmt.Println(v.Next())
		fmt.Println(v.Next())
		fmt.Println(v.Next())
		fmt.Println(v.Next())
		fmt.Println(v.Next())
	}
}
