package opensshutil

import (
	"errors"
)

// OpenSSH configuration file scanner.
//
// The configuration file format is pretty straight forward except for the
// caveat of equal signs (`=`) being contextual. An equal sign, when
// encountered as part of the first token or when encountered immediately after
// the first token is an assignment operator, while equal signs encountered in
// later tokens are part of the token.
type configScanner struct {
	src       []byte
	dOffset   int
	offset    int
	lTokCount int
	lHasEqual bool
	expectWs  bool
	err       error
	ch        rune
}

// New OpenSSH configuration file scanner.
func newConfigScanner(src []byte) *configScanner {
	s := &configScanner{
		src: src,
	}
	s.next()

	return s
}

// Assign error.
func (s *configScanner) error(msg string) {
	s.err = errors.New(msg)
}

// Read the next ASCII character into s.ch.
//
// s.ch is < 0 indicates end-of-file.
func (s *configScanner) next() {
	if s.dOffset < len(s.src) {
		s.offset = s.dOffset

		if s.ch == '\n' {
			s.lTokCount = 0
			s.lHasEqual = false
		}

		s.ch = rune(s.src[s.dOffset])
		s.dOffset++
	} else {
		s.offset = len(s.src)
		s.ch = -1
	}
}

// Skip white space.
func (s *configScanner) skipWhiteSpace() {
	for isWhiteSpace(s.ch) {
		s.next()
	}
}

// Skip comment.
func (s *configScanner) skipComment() {
	for s.ch != '\n' && s.ch != -1 {
		s.next()
	}
}

// Scan string.
//
// Scan until a white space, EOL, EOF or non-printable character is met.
func (s *configScanner) scanString() string {
	start := s.offset
	s.next()

	for s.ch != -1 && s.ch != '\n' && !isWhiteSpace(s.ch) && isPrintable(s.ch) {
		s.next()
	}

	return string(s.src[start:s.offset])
}

// Scan quoted string.
func (s *configScanner) scanQuotedString() string {
	qch := s.ch
	start := s.offset + 1
	s.next()

	for s.ch != -1 && s.ch != '\n' && s.ch != qch {
		s.next()
	}

	switch s.ch {
	case -1:
		s.error("unexpected end of file while reading quoted string")

	case '\n':
		s.error("unexpected end of line while reading quoted string")

	case qch:
		s.expectWs = true
		s.next()
	}

	return string(s.src[start : s.offset-1])
}

// Scan keyword or string.
func (s *configScanner) scanKeywordOrString() (configToken, string) {
	start := s.offset
	s.next()

	// Attempt to scan as keyword.
	for isAlphanumeric(s.ch) {
		s.next()
	}

	// If we reached end of line, end of file, white space or an equal sign in
	// the first line token, return as keyword.
	if s.ch == -1 || s.ch == '\n' || isWhiteSpace(s.ch) || (s.lTokCount == 1 && s.ch == '=') {
		return CT_KEYWORD, string(s.src[start:s.offset])
	}

	// Scan as string.
	for s.ch != -1 && s.ch != '\n' && !isWhiteSpace(s.ch) && isPrintable(s.ch) {
		s.next()
	}

	return CT_STRING, string(s.src[start:s.offset])
}

// Scan advances the line scanner to the next token.
//
// Returns false when the scan stops, either by reaching the end of the input
// or an error.
//
// Argument tok will contain the token type after the scan, while lit contains
// the literal data representation. For quoted string, lit contains the
// unquoted string.
func (s *configScanner) Scan() (tok configToken, lit string) {
scan:
	// Check for expected white space.
	if s.expectWs {
		if s.ch != -1 && s.ch != '\n' && !isWhiteSpace(s.ch) {
			s.error("expected white space or end-of-line")
		}
		s.expectWs = false
	}

	// Skip white space.
	s.skipWhiteSpace()

	// Determine token type from the first character.
	lTokCount := s.lTokCount
	s.lTokCount++

	switch {
	case isLetter(s.ch):
		// Scan as either a keyword or a string.
		tok, lit = s.scanKeywordOrString()

	case s.ch == '\n':
		tok = CT_EOL
		s.next()

	case s.ch == '=':
		// If the line token count is 1, we treat the equal sign as an
		// individual token.
		if lTokCount == 1 {
			tok, lit = CT_EQUAL, "="
			s.next()
		} else {
			tok = CT_STRING
			lit = s.scanString()
		}

	case s.ch == '#':
		// Skip the comment and repeat the scan.
		s.skipComment()
		goto scan

	case s.ch == '"':
		// Scan the quoted string.
		tok = CT_STRING
		lit = s.scanQuotedString()

	case isPrintable(s.ch):
		// Any other printable character is treated as the beginning of an
		// unquoted string.
		tok = CT_STRING
		lit = s.scanString()

	case s.ch == -1:
		tok = CT_EOF

	default:
		tok = CT_ILLEGAL
		lit = string(s.ch)
	}

	return
}

// Error.
func (s *configScanner) Err() error {
	return s.err
}
