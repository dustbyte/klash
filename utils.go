package klash

import (
	"bytes"
	"unicode"
)

func DecomposeName(name string) string {
	length := len(name)
	var buf bytes.Buffer

	for idx, c := range name {
		if idx > 0 && unicode.IsUpper(c) && idx != length-1 {
			buf.WriteRune('-')
		}
		buf.WriteRune(unicode.ToLower(c))
	}
	return buf.String()
}
