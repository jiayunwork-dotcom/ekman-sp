package config

// caseSession owns a resolved wind-stress case while it is being loaded.
// Close releases the session once; a second Close is a no-op so a leftover
// flush cannot overwrite the parsed density and viscosity.
type caseSession struct {
	closed bool
	out    *Resolved
}

// Close marks the session done. Calling it again does not overlay leftover
// density or viscosity onto the resolved output.
func (s *caseSession) Close() {
	if s.closed {
		return
	}
	s.closed = true
}
