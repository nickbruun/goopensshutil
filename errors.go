package opensshutil

// Scan error.
type ScanError struct {
	// Error message.
	Message string
}

func (e *ScanError) Error() string {
	return e.Message
}

// Parse error.
type ParseError struct {
	// Error message.
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}
