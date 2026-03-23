package confusables

import (
	"strings"
	"unicode/utf8"

	"github.com/corazawaf/coraza/v3/experimental/plugins"
)

func init() {
	plugins.RegisterTransformation("confusableNormalize", confusableNormalize)
}

func confusableNormalize(input string) (string, bool, error) {
	var b strings.Builder
	b.Grow(len(input))
	changed := false

	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])

		if isZeroWidth(r) {
			changed = true
			i += size
			continue
		}

		if mapped, ok := confusableMap[r]; ok {
			b.WriteByte(mapped)
			changed = true
		} else {
			b.WriteRune(r)
		}
		i += size
	}

	return b.String(), changed, nil
}
