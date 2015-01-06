package klash

import (
	"bytes"
	"unicode"
)

func DecomposeName(name string, lower bool) string {
	length := len(name)
	var buf bytes.Buffer
	var caseChange func(rune) rune
	var joinRune rune

	if lower {
		caseChange = unicode.ToLower
		joinRune = '-'
	} else {
		caseChange = unicode.ToUpper
		joinRune = '_'
	}

	for idx, c := range name {
		if idx > 0 && unicode.IsUpper(c) && idx != length-1 {
			buf.WriteRune(joinRune)
		}
		buf.WriteRune(caseChange(c))
	}
	return buf.String()
}

func Dashed(name string) string {
	if len(name) > 1 {
		return "--" + name
	}
	return "-" + name
}
