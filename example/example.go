package nan

import (
	"fmt"
	"math"
)

func exp(x float64) float64 {
	return math.Exp(x)
}

func ExampleDivide() {
	fmt.Println(1.0 / math.Exp(1.0))
}
