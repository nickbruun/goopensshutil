package opensshutil

// OpenSSH configuration file token.
type configToken int

const (
	CT_ILLEGAL configToken = iota

	// End-of-file.
	CT_EOF

	// End-of-line.
	CT_EOL

	// Equal (=).
	CT_EQUAL

	// Keyword (alphanumeric starting with a letter.)
	CT_KEYWORD

	// String.
	CT_STRING
)

func (t configToken) String() string {
	if t > CT_STRING {
		return "CT_UNKNOWN"
	}

	return map[configToken]string{
		CT_ILLEGAL: "CT_ILLEGAL",
		CT_EOF:     "CT_EOF",
		CT_EOL:     "CT_EOL",
		CT_EQUAL:   "CT_EQUAL",
		CT_KEYWORD: "CT_KEYWORD",
		CT_STRING:  "CT_STRING",
	}[t]
}
