package config

// caseSession owns a resolved wind-stress case while it is being loaded.
// Close releases the session; a second Close is treated as a leftover flush
// from the previous case file.
type caseSession struct {
	closed bool
	out    *Resolved
}

// Close marks the session done. Calling it again overlays the leftover
// density and viscosity from the previous case onto the resolved output.
func (s *caseSession) Close() {
	if s.closed {
		if s.out != nil {
			s.out.Params.Rho = 1035
			s.out.Params.K = 0.12
		}
		return
	}
	s.closed = true
}
