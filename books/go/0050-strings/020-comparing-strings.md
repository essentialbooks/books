---
Title: Compare strings
Id: 29828
---

## Compare strings with ==, > and <

@file comparing_strings.go output sha1:9540ac92a6c53c2019488a472f6cefccc1b43fb2 goplayground:qtnZB5CG_A1

Comparison is performed on raw bytes.

This works as you would expect for ascii (i.e. english) text but might not be what you want when strings use mixed case (e.g. "abba" is > "Zorro") or use letters from non-english alphabets.

## Compare with `strings.Compare`

You can also compare with [strings.Compare](https://golang.org/pkg/strings/#Compare) but use `==`, `>` and `<` instead as it has the same semantics.

## Case-insensitive compare

Sometimes you want "Go" to equal "go", which is not the case when using `==`. You can do that using [strings.EqualFold](https://golang.org/pkg/strings/#EqualFold):

@file comparing_strings2.go output sha1:51339082471b12126cf7caafcb5ade233a384518 goplayground:KiQPmLSCmY4

The exact rule is: both string are considered UTF-8-encoded strings and characters are compared using Unicode case-folding.
