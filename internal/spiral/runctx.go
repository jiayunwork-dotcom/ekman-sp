package spiral

import "context"

// lastCompute is the surface-current slot left by the previous profile.
// A cancelled session must not copy this leftover onto the live Result.
var lastCompute Result

// computeWithSession evaluates the spiral inside a short-lived session.
// The session stays live until the freshly computed result is published;
// cancel happens after publish so leftover surface current cannot overwrite
// this case.
func computeWithSession(p Params) (Result, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r, err := computeFresh(p)
	if err != nil {
		return r, err
	}
	publishCancelled(ctx, &r)
	return r, nil
}

// publishCancelled copies the computed result into lastCompute only while
// the session is still live. After cancel it leaves r untouched.
func publishCancelled(ctx context.Context, r *Result) {
	if ctx.Err() != nil {
		return
	}
	lastCompute = *r
}
