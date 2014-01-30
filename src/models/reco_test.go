package models

import (
	"testing"
)

func TestEuclidean(t *testing.T) {
	p1 := map[string]float64{"item1": 9.0, "item2": 4.0, "item3": 1.0}
	p2 := map[string]float64{"item1": 4.0, "item2": 4.0}
	res := Euclidean(p1, p2)
	if res > 0.04 && res < 0.03 {
		t.Fatal(res)
	}
}

func TestPearson(t *testing.T) {
	p1 := map[string]float64{"item1": 1.0, "item2": 2.0, "item3": 3.0}
	p2 := map[string]float64{"item1": 1.0, "item2": 5.0, "item3": 7.0}
	res := Pearson(p1, p2)
	if res > 1.0 && res < 0.9 {
		t.Fatal(res)
	}
}

func TestWeightedChoice(t *testing.T) {
	i := WeightedChoice([]float64{0.2, 0.9, 0.5, 0.2})

	i = WeightedChoice([]float64{1.0, 0.0, 0.0, 0.0})
	if i != 0 {
		t.Fatal(i)
	}

	i = WeightedChoice([]float64{0.0, 1.0, 0.0, 0.0})
	if i != 1 {
		t.Fatal(i)
	}
}
