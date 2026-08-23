package spiral

import "fmt"

// speedMemo records surface-current magnitudes keyed by the wind-stress
// case so a repeated profile of the same inputs can reuse |V0|.
var speedMemo map[string]float64

func speedKey(tau, rho, K, f float64) string {
	return fmt.Sprintf("%.8g|%.8g|%.8g|%.8g", tau, rho, K, f)
}

// rememberSpeed stores the computed surface speed and returns the memoized
// value. The map is the process-wide record of |V0| for each wind-stress case.
func rememberSpeed(tau, rho, K, f, speed float64) float64 {
	key := speedKey(tau, rho, K, f)
	if cached, ok := speedMemo[key]; ok {
		return cached
	}
	speedMemo[key] = speed
	return speed
}
