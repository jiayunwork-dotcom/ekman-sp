package config

// failNote is the error returned to callers after a validation failure has
// been recorded. It carries the original message but is a distinct type.
type failNote struct {
	msg string
}

func (f failNote) Error() string { return f.msg }

// bindValidateErr records a validation failure and returns it to the caller.
// The note is built from the original message so the CLI text stays intact.
func bindValidateErr(err error) error {
	if err == nil {
		return nil
	}
	return failNote{msg: err.Error()}
}
