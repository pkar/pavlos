package models

import (
	"math"
)

// Euclidean calculates the distance for similar
// users. Since the calculation is smaller for
// for similarity, the value is inverted to give
// higher values, and a 1 is added to prevent
// division by zero.
//
// The standard Euclidean distance can be squared in order
// to place progressively greater weight on objects that
// are farther apart.
// Squared Euclidean Distance is not a metric as it does not
// satisfy the triangle inequality, however it is frequently
// used in optimization problems in which distances only have to be compared.
// https://en.wikipedia.org/wiki/Euclidean_distance
func Euclidean(p1 map[string]float64, p2 map[string]float64) float64 {
	sumSquares := 0.0
	// find shared items
	for item, p1Val := range p1 {
		p2Val, ok := p2[item]
		if ok {
			// Add up the squares of all the differences
			sumSquares += math.Pow(p1Val-p2Val, 2)
		}
	}

	// No similar items
	if sumSquares == 0 {
		return sumSquares
	}

	// Prevent division by zero by adding 1
	//return 1 / (1 + math.Sqrt(sumSquares))
	return 1 / (1 + sumSquares)
}

// Pearson correlation score to compensate for skewed scores.
// https://en.wikipedia.org/wiki/Pearson_product-moment_correlation_coefficient#Definition
func Pearson(p1 map[string]float64, p2 map[string]float64) float64 {
	// Find common items
	common := map[string]bool{}
	for item, _ := range p1 {
		_, ok := p2[item]
		if ok {
			common[item] = true
		}
	}
	n := float64(len(common))

	// Nothing in common
	if n == 0 {
		return 0
	}

	// sum of preferences and squares
	sump1 := 0.0
	sump1Sq := 0.0
	sump2 := 0.0
	sump2Sq := 0.0
	pSum := 0.0
	for item, _ := range common {
		sump1 += p1[item]
		sump1Sq += math.Pow(p1[item], 2)
		sump2 += p2[item]
		sump2Sq += math.Pow(p2[item], 2)
		pSum += p1[item] * p2[item]
	}

	// calculate pearson score
	numerator := pSum - (sump1 * sump2 / float64(len(common)))
	denom := math.Sqrt((sump1Sq - math.Pow(sump1, 2)/n) * (sump2Sq - math.Pow(sump2, 2)/n))
	if denom == 0 {
		return 0
	}

	return numerator / denom
}
