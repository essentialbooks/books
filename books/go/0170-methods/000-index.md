---
Title: Methods
Id: 3890
---

A method is a function tied to a type, most commonly a struct.

This is similar to classes in languages like Java or C++.

Basics of methods:

@file index.go output sha1:1016a0f355a0ac8b453d890bca63087c3783311f goplayground:c9MI-O2pvLt

In the above example method `PrintFullName` takes a reciver named `p` of type `*Person`.

People coming from other languages are tempted to name the receiver `this` (mimicking C++) or `self` (mimicking Python).

In Go, the rule for naming reciver is:
* be short
* be consistent across methods

## Value vs. pointer receiver

Method receiver can be either a value and a pointer.

@file index2.go output sha1:18e4acb218010b37cc87e48b475ede8d9df4b6b7 goplayground:2KGrg2M4pj2

As you can see, when `p` is of type `Person`, we can call both methods defined for `Person` and `*Person`. Go will automatically convert receiver `Person` to `*Person`.

It doesn't work the other way i.e. you can't call method with receiver `Person` on value with type `*Person`.

<!-- TODO: write more -->
