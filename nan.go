package nan

import (
	"log"
	"math"
)

// START OMIT
func Log(x float64) float64 {
	if x <= 0.0 {
		log.Panicf("x (%v) <= 0", x)
	}
	return math.Log(x)
}

// END OMIT
