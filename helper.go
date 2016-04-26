package main

import "strings"

// stolen from:
// http://stackoverflow.com/questions/8689425/removed-last-character-of-a-string
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
