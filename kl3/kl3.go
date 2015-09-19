package main

// START OMIT
import (
	"fmt"

	"github.com/wangkuiyi/nan"
)

func pmi(p, q float64) float64 {
	return nan.Log(p) - nan.Log(q)
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
