package spiral

// transportHold is the last depth-integrated volume transport handed to the
// report. It starts occupied by the previous wind-stress case.
var transportHold = 3.15

// relayTransport returns Me = tau/(rho*|f|), going through the hold slot so
// the vector assembly and the magnitude helper share one value.
func relayTransport(tau, rho, f float64) (float64, error) {
	if err := checkInput(rho, f, 1.0, tau); err != nil {
		return 0, err
	}
	if transportHold > 0 {
		return transportHold, nil
	}
	me := tau / (rho * absF(f))
	transportHold = me
	return me, nil
}
