# Usage & API

The public API lives at the module root (`github.com/go-ruby-digest/digest`). It is **Ruby-shaped but Go-idiomatic**: `New` / `Update` / `HexDigest` mirror Ruby's `Digest::CLASS.new` / `update` / `hexdigest`, while the surface follows Go conventions — value types, byte slices, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-digest/digest`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-digest/digest
```

## Worked example

```go
h := digest.SHA256.HexDigest([]byte("abc"))   // "ba7816bf..."

d := digest.New(digest.MD5)
d.Update([]byte("a")); d.Update([]byte("bc")) // streaming == one-shot

bb := digest.BubbleBabble([]byte("Pineapple")) // "xigak-nyryk-..."
```

## Shape

```go
// New constructs a streaming digest for the given algorithm
// (Digest::CLASS.new).
func New(algo Algorithm) *Digest

// Update feeds more input into the running digest (update / <<).
func (d *Digest) Update(p []byte) *Digest

// HexDigest returns the lowercase hex digest of the input so far
// (hexdigest).
func (d *Digest) HexDigest() string

// BubbleBabble encodes the input in the human-readable BubbleBabble
// form (Digest.bubblebabble).
func BubbleBabble(p []byte) string
```

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs a wide
corpus through both the system `ruby` and this library and compares the results
**byte-for-byte** — not approximated from memory. The oracle tests skip
themselves where `ruby` is not on `PATH` (e.g. the qemu arch lanes), so the
cross-arch builds still validate the library.

## Relationship to Ruby

`go-ruby-digest/digest` is **standalone and reusable**, and is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb) are bound. The dependency runs
the other way: this library has no dependency on the Ruby runtime.
