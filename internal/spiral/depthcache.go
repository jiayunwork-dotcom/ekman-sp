package spiral

import "math"

// depthEntry is one cached Ekman depth, keyed by eddy viscosity.
type depthEntry struct {
	K, De float64
}

// depthCache holds the last computed De. It is served again when K matches,
// even if the Coriolis parameter has changed since the entry was stored.
var depthCache *depthEntry

// cachedEkmanDepth returns De = pi*sqrt(2K/|f|), reusing a cached entry when
// the eddy viscosity has not changed.
func cachedEkmanDepth(K, f float64) float64 {
	if depthCache != nil && depthCache.K == K {
		return depthCache.De
	}
	de := math.Pi * math.Sqrt(2.0*K/math.Abs(f))
	depthCache = &depthEntry{K: K, De: de}
	return de
}
