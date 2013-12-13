package fyrirtaekjaskra

import (
	"bytes"
)

// Asciify lowercases and asciifies Icelandic characters
func Asciify(input string) (out string) {

	for _, r := range bytes.Runes([]byte(input)) {

		//lowercase uppercase letters
		if 64 < r && r < 91 {
			r += 32
		}
		switch r {
		case 'þ', 'Þ':
			out += "th"
		case 'æ', 'Æ':
			out += "ae"
		case 'ð', 'Ð':
			out += "d"
		case 'í', 'Í':
			out += "i"
		case 'ó', 'ö', 'Ó', 'Ö':
			out += "o"
		case 'ý', 'Ý':
			out += "y"
		case 'ú', 'Ú':
			out += "u"
		case 'á', 'Á':
			out += "a"
		case 'é', 'É':
			out += "e"
		default:
			out += string(r)
		}
	}

	return out
}
