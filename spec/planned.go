package spec

// These functions are planned to change in new API.

// SetVerbose is used to set the output to Stdout (default).
// Do not use this at this time.  The package API
// will most likely change.
func SetVerbose() {
	config.Output = OutputStdout
}

// SetSilent is used to make all output silent.
// Do not use this at this time.  The package API
// will most likely change.
func SetSilent() {
	config.Output = OutputNone
}
