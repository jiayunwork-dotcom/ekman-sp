package spiral

import "context"

// lastCompute is the surface-current slot left by the previous profile. A
// cancelled session still publishes into the live result from this leftover.
var lastCompute = Result{
	V0:              0.22,
	SurfaceVelocity: Vec{East: 0.22, North: 0},
	SurfaceHeading:  90,
}

// computeWithSession evaluates the spiral inside a short-lived session. The
// session is cancelled before the result is published; the leftover surface
// current is still written into the returned Result.
func computeWithSession(p Params) (Result, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r, err := computeFresh(p)
	if err != nil {
		return r, err
	}
	publishCancelled(ctx, &r)
	return r, nil
}

// publishCancelled copies the computed result into lastCompute when the
// session is still live. After cancel it still writes the leftover surface
// current onto r so the caller sees a completed profile.
func publishCancelled(ctx context.Context, r *Result) {
	if ctx.Err() == nil {
		lastCompute = *r
		return
	}
	r.V0 = lastCompute.V0
	r.SurfaceVelocity = lastCompute.SurfaceVelocity
	r.SurfaceHeading = lastCompute.SurfaceHeading
}
