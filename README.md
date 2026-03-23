# coraza-leet-normalize

Confusable character normalization for [Coraza WAF](https://coraza.io). Registers a `t:confusableNormalize` transformation that defeats evasion via leet-speak, Cyrillic/Greek homoglyphs, and zero-width Unicode characters.

[![CI](https://github.com/Traceableai/coraza-leet-normalize/actions/workflows/ci.yml/badge.svg)](https://github.com/Traceableai/coraza-leet-normalize/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Traceableai/coraza-leet-normalize.svg)](https://pkg.go.dev/github.com/Traceableai/coraza-leet-normalize)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

## Problem

Attackers evade WAF regex patterns using visually similar character substitutions:

```
"4cc3ss d3n13d"              (leet-speak)
"аccess denied"              (Cyrillic а looks identical to Latin a)
"a​c​c​e​s​s"               (zero-width spaces between every letter)
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
SecRule ARGS "@rx (?i)access\s+denied" \
    "id:1001,phase:2,t:none,t:urlDecodeUni,t:confusableNormalize,deny,msg:'Evasion attempt detected'"
```

The transformation chain: raw input → URL decode → normalize confusables → regex match.

## What It Catches

Given a WAF rule matching `(?i)access\s+denied`, here's what happens with and without the plugin:

```
Evasion technique   | Attack payload                              | Without plugin | With plugin
--------------------|---------------------------------------------|----------------|------------
leet-speak          | 4cc3ss d3n13d                               | BYPASSED       | BLOCKED
cyrillic homoglyphs | [U+0430]ccess denied                        | BYPASSED       | BLOCKED
greek homoglyphs    | [U+03B1]ccess denied                        | BYPASSED       | BLOCKED
zero-width chars    | a[ZW]c[ZW]c[ZW]e[ZW]s[ZW]s denied          | BYPASSED       | BLOCKED
mixed evasion       | [U+0430]cc3[ZW]ss d3n13d                    | BYPASSED       | BLOCKED
benign input        | order 12345 confirmed                       | allowed        | allowed
```

`[U+0430]` = Cyrillic а (looks identical to Latin a), `[U+03B1]` = Greek α, `[ZW]` = zero-width space. These characters are visually indistinguishable from their Latin equivalents, so they bypass regex rules undetected.

Every evasion technique slips past the raw regex but gets caught after normalization. Benign input remains unaffected — no false positives.

## What It Normalizes

### Leet-speak (8 mappings)

| Input | Output | Example                      |
| ----- | ------ | ---------------------------- |
| `0`   | `o`    | `p4ssw0rd` → `password`     |
| `1`   | `i`    | `adm1n` → `admin`           |
| `3`   | `e`    | `s3cur3` → `secure`         |
| `4`   | `a`    | `4cc3ss` → `access`         |
| `5`   | `s`    | `5erver` → `server`         |
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
