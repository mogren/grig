package vose

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNewVose(t *testing.T) {
	type args struct {
		prob      []float64
		generator *rand.Rand
	}
	arg1 := args{prob: []float64{0.1, 0.2, 0.3, 0.4}, generator: rand.New(rand.NewSource(1))}
	arg2 := args{prob: []float64{}, generator: rand.New(rand.NewSource(1))}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Successful init", arg1, false},
		{"Missing probability", arg2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewVose(tt.args.prob, tt.args.generator)
			fmt.Printf("Vose: %v\n", v)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVose() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestVose_Next(t *testing.T) {
	type fields struct {
		limit     int
		prob      []float64
		alias     []int
		generator rand.Rand
	}
	test1 := fields{
		limit:     4,
		prob:      []float64{0.4, 0.8, 1.0, 0.8000000000000003},
		alias:     []int{3, 3, 0, 2},
		generator: *rand.New(rand.NewSource(1)),
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"Basic test 1", test1, 2},
		{"Basic test 2", test1, 3},
		{"Basic test 3", test1, 2},
		{"Basic test 4", test1, 1},
		{"Basic test 5", test1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vose{
				limit:     tt.fields.limit,
				prob:      tt.fields.prob,
				alias:     tt.fields.alias,
				generator: tt.fields.generator,
			}
			if got := v.Next(); got != tt.want {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}
