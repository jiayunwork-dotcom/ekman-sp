package spiral

// transportHold is the last depth-integrated volume transport handed to the
// report. The hold is always overwritten from this call's wind stress so a
// leftover slot cannot be shared by magnitude and vector assembly.
var transportHold float64

// relayTransport returns Me = tau/(rho*|f|), going through the hold slot so
// the vector assembly and the magnitude helper share one freshly computed
// value for this case.
func relayTransport(tau, rho, f float64) (float64, error) {
	if err := checkInput(rho, f, 1.0, tau); err != nil {
		return 0, err
	}
	me := tau / (rho * absF(f))
	transportHold = me
	return me, nil
}
