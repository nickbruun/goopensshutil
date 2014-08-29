package opensshutil

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func isPrintable(r rune) bool {
	return r >= 0x20 && r <= 0x7e
}

func isWhiteSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r'
}
