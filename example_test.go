package confusables

import (
	"fmt"
	"regexp"
)

func Example_evasionDetection() {
	rule := regexp.MustCompile(`(?i)access\s+denied`)

	attacks := []struct {
		label string
		input string
	}{
		{"leet-speak", "4cc3ss d3n13d"},
		{"cyrillic homoglyphs", "\u0430cc\u0435ss denied"},
		{"greek homoglyphs", "\u03B1ccess denied"},
		{"zero-width chars", "a\u200Bc\u200Bc\u200Be\u200Bs\u200Bs denied"},
		{"mixed evasion", "4cc\u200B\u0435ss d3n13d"},
		{"benign input", "order 12345 confirmed"},
	}

	for _, a := range attacks {
		normalized, _, _ := confusableNormalize(a.input)
		withoutPlugin := rule.MatchString(a.input)
		withPlugin := rule.MatchString(normalized)

		fmt.Printf("%-22s | raw match: %-5v | normalized match: %-5v | normalized: %q\n",
			a.label, withoutPlugin, withPlugin, normalized)
	}

	// Output:
	// leet-speak             | raw match: false | normalized match: true  | normalized: "access denied"
	// cyrillic homoglyphs    | raw match: false | normalized match: true  | normalized: "access denied"
	// greek homoglyphs       | raw match: false | normalized match: true  | normalized: "access denied"
	// zero-width chars       | raw match: false | normalized match: true  | normalized: "access denied"
	// mixed evasion          | raw match: false | normalized match: true  | normalized: "access denied"
	// benign input           | raw match: false | normalized match: false | normalized: "order i2eas confirmed"
}
