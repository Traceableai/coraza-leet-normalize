package confusables

import "testing"

func TestConfusableNormalize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		changed bool
	}{
		// Leet-speak
		{"leet access denied", "4cc3ss d3n13d", "access denied", true},
		{"leet bypass filter", "byp4$$ c0nt3nt f1lt3r5", "bypass content filters", true},
		{"leet admin mode", "4dm1n m0d3", "admin mode", true},

		// Cyrillic
		{"cyrillic access", "\u0430cc\u0435ss", "access", true},
		{"cyrillic server", "s\u0443st\u0435m s\u0435rv\u0435r", "system server", true},

		// Greek
		{"greek omicron", "p\u03BFrt", "port", true},
		{"greek alpha", "\u03B1dmin", "admin", true},

		// Zero-width
		{"zwsp insertion", "a\u200Bc\u200Bc\u200Be\u200Bs\u200Bs", "access", true},
		{"soft hyphen", "sys\u00ADtem", "system", true},
		{"word joiner", "se\u2060rver", "server", true},
		{"bom", "\uFEFFaccess", "access", true},

		// Mixed evasion
		{"leet+cyrillic+zw", "4cc\u200B3ss\u200B\u0435rr0r", "accesserror", true},
		{"full evasion chain", "\u0430cc\u200B\u03B5ss d3n13d", "access denied", true},

		// Passthrough (no change)
		{"normal ascii", "access denied", "access denied", false},
		{"empty string", "", "", false},
		{"pure ascii symbols", "hello world!", "hello world!", false},

		// Documented side effect: digits always map
		{"digits always map", "12345", "i2eas", true},
		{"order with digits", "order 12345", "order i2eas", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed, err := confusableNormalize(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("confusableNormalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
			if changed != tt.changed {
				t.Errorf("confusableNormalize(%q) changed = %v, want %v", tt.input, changed, tt.changed)
			}
		})
	}
}

func BenchmarkConfusableNormalize(b *testing.B) {
	input := "4cc\u200B3ss\u0435 d3n13d s\u04353rv3r"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		confusableNormalize(input)
	}
}

func BenchmarkConfusableNormalizeNoChange(b *testing.B) {
	input := "access denied server error"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		confusableNormalize(input)
	}
}
