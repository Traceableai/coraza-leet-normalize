package confusables

var confusableMap = map[rune]byte{
	// Leet-speak
	'0': 'o',
	'1': 'i',
	'3': 'e',
	'4': 'a',
	'5': 's',
	'7': 't',
	'@': 'a',
	'$': 's',

	// Cyrillic lowercase
	'\u0430': 'a', // а
	'\u0435': 'e', // е
	'\u043E': 'o', // о
	'\u0440': 'p', // р
	'\u0441': 'c', // с
	'\u0443': 'y', // у
	'\u0456': 'i', // і
	'\u0445': 'x', // х

	// Cyrillic uppercase
	'\u0410': 'A', // А
	'\u0415': 'E', // Е
	'\u041E': 'O', // О
	'\u0420': 'P', // Р
	'\u0421': 'C', // С
	'\u0422': 'T', // Т
	'\u0425': 'X', // Х

	// Greek lowercase
	'\u03BF': 'o', // ο
	'\u03B1': 'a', // α
	'\u03B5': 'e', // ε
	'\u03C1': 'p', // ρ
	'\u03B9': 'i', // ι

	// Greek uppercase
	'\u039F': 'O', // Ο
}

func isZeroWidth(r rune) bool {
	switch r {
	case '\u200B', // Zero Width Space
		'\u200C', // Zero Width Non-Joiner
		'\u200D', // Zero Width Joiner
		'\uFEFF', // Byte Order Mark / ZWNBSP
		'\u00AD', // Soft Hyphen
		'\u2060', // Word Joiner
		'\u180E': // Mongolian Vowel Separator
		return true
	}
	return false
}
