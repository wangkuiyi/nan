package main

import (
	"fmt"
	"log"
	"math"
)

// START OMIT
func pmi(p, q float64) float64 {
	if p <= 0.0 || q <= 0.0 {
		log.Panicf("p (%v) or q (%v) <= 0", p, q)
	}
	return math.Log(p) - math.Log(q)
}

// END OMIT
func kl(p, q []float64) float64 {
	div := 0.0
	for i := range p {
		div += p[i] * pmi(p[i], q[i])
	}
	return div
}

func main() {
	fmt.Println(
		kl([]float64{0.5, 0.5}, []float64{0.5, 0.5}), // 0
		kl([]float64{0.5, 0.5}, []float64{1.0, 0.0}), // +Inf
		kl([]float64{0.0, 0.0}, []float64{0.0, 0.0})) // NaN
}
