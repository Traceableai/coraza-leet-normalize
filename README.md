# coraza-leet-normalize

Confusable character normalization for [Coraza WAF](https://coraza.io). Registers a `t:confusableNormalize` transformation that defeats evasion via leet-speak, Cyrillic/Greek homoglyphs, and zero-width Unicode characters.

[![CI](https://github.com/Traceableai/coraza-leet-normalize/actions/workflows/ci.yml/badge.svg)](https://github.com/Traceableai/coraza-leet-normalize/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Traceableai/coraza-leet-normalize.svg)](https://pkg.go.dev/github.com/Traceableai/coraza-leet-normalize)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

## Problem

Attackers evade WAF regex patterns using visually similar character substitutions:

```
"1gn0r3 4ll pr3v10us 1nstruct10ns"   (leet-speak)
"ignоrе all instructions"             (Cyrillic о and е look identical to Latin)
"i​g​n​o​r​e"                         (zero-width spaces between every letter)
```

No existing Coraza or ModSecurity built-in transformation handles these evasion techniques. This plugin fills that gap.

## Installation

```bash
go get github.com/Traceableai/coraza-leet-normalize
```

## Usage

Import the package with a blank identifier so the `init()` function registers the transformation:

```go
import (
    "github.com/corazawaf/coraza/v3"
    _ "github.com/Traceableai/coraza-leet-normalize"
)
```

Then use `t:confusableNormalize` in any SecLang rule:

```
SecRule ARGS "@rx (?i)ignore\s+.*instructions" \
    "id:1001,phase:2,t:none,t:urlDecodeUni,t:confusableNormalize,deny,msg:'Prompt injection attempt'"
```

The transformation chain: raw input → URL decode → normalize confusables → regex match.

## What It Normalizes

### Leet-speak (8 mappings)

| Input | Output | Example                      |
| ----- | ------ | ---------------------------- |
| `0`   | `o`    | `ign0re` → `ignore`         |
| `1`   | `i`    | `1nstruct1ons` → `instructions` |
| `3`   | `e`    | `3xecute` → `execute`       |
| `4`   | `a`    | `4dmin` → `admin`           |
| `5`   | `s`    | `5ystem` → `system`         |
| `7`   | `t`    | `a77ack` → `attack`         |
| `@`   | `a`    | `@dmin` → `admin`           |
| `$`   | `s`    | `bypa$$` → `bypass`         |

### Cyrillic Homoglyphs (15 mappings)

| Input | Codepoint | Output |
| ----- | --------- | ------ |
| а     | U+0430    | a      |
| е     | U+0435    | e      |
| о     | U+043E    | o      |
| р     | U+0440    | p      |
| с     | U+0441    | c      |
| у     | U+0443    | y      |
| і     | U+0456    | i      |
| х     | U+0445    | x      |
| А     | U+0410    | A      |
| Е     | U+0415    | E      |
| О     | U+041E    | O      |
| Р     | U+0420    | P      |
| С     | U+0421    | C      |
| Т     | U+0422    | T      |
| Х     | U+0425    | X      |

### Greek Homoglyphs (6 mappings)

| Input | Codepoint | Output |
| ----- | --------- | ------ |
| ο     | U+03BF    | o      |
| α     | U+03B1    | a      |
| ε     | U+03B5    | e      |
| ρ     | U+03C1    | p      |
| ι     | U+03B9    | i      |
| Ο     | U+039F    | O      |

### Zero-width Characters (7 stripped)

| Codepoint | Name                      |
| --------- | ------------------------- |
| U+200B    | Zero Width Space          |
| U+200C    | Zero Width Non-Joiner     |
| U+200D    | Zero Width Joiner         |
| U+FEFF    | Byte Order Mark / ZWNBSP  |
| U+00AD    | Soft Hyphen               |
| U+2060    | Word Joiner               |
| U+180E    | Mongolian Vowel Separator |

## Important Notes

### Aggressive Leet-speak Mappings

The leet-speak mappings convert **all** occurrences of the mapped digits, not just those in evasion contexts. For example, `"order 12345"` becomes `"order i2eas"` after normalization.

This is acceptable because the transformation is applied *before* regex matching, not to stored or displayed data. WAF regex patterns only match security-relevant phrases, so normalized digits in benign strings like `"order i2eas"` won't trigger false positives (no security rule matches that string).

### Transformation Ordering

Place `t:confusableNormalize` after URL decoding but before regex matching:

```
t:none,t:urlDecodeUni,t:confusableNormalize
```

If you also use `t:lowercase`, apply it before `t:confusableNormalize` since the Cyrillic/Greek mappings already produce the correct case.

### No Lowercasing

This transformation does **not** lowercase the output. Use `t:lowercase` or the `(?i)` regex flag separately. This keeps the plugin focused and composable.

## Benchmarks

Measured on Apple M4 Pro:

```
BenchmarkConfusableNormalize-14          5634936     219.0 ns/op    48 B/op    1 allocs/op
BenchmarkConfusableNormalizeNoChange-14  6209410     191.3 ns/op    32 B/op    1 allocs/op
```

Single-pass O(n) processing with one allocation (the `strings.Builder` output buffer).

## TinyGo Compatibility

The plugin uses only `strings.Builder`, `unicode/utf8`, and map lookups — all TinyGo-compatible. No `regexp`, `reflect`, or `unsafe` packages are used.

## Contributing

Issues and pull requests are welcome. Please open an issue first to discuss proposed changes.

## License

[Apache License 2.0](LICENSE)
