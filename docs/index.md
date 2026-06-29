# go-ruby-digest documentation

**Ruby's Digest message-digest suite in pure Go — MRI-compatible, no cgo.**

`go-ruby-digest/digest` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's Digest
message-digest suite, matching reference Ruby (MRI) byte-for-byte. The module path is
`github.com/go-ruby-digest/digest`.

It was **extracted from rbgo's prelude/internals into a reusable standalone
library**: the module is standalone and importable by any Go program, and it is
the backend bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)
by `rbgo` as a native module — just like
[go-ruby-regexp](https://github.com/go-ruby-regexp) and
[go-ruby-erb](https://github.com/go-ruby-erb). The dependency runs the other
way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: digest suite complete — MRI byte-exact"
    MD5, SHA1, the **SHA2** family (SHA256/384/512) and **RIPEMD-160**, plus **BubbleBabble**, through the **`Digest::Instance`** streaming protocol (`update` / `<<` / `reset` / `finish`) and the one-shot class methods, in **`hexdigest`** / **`base64digest`** / raw **`digest`** forms. Validated by a **differential oracle** against the system `ruby` / `digest` — every digest compared byte-for-byte — at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
h := digest.SHA256.HexDigest([]byte("abc"))   // "ba7816bf..."

d := digest.New(digest.MD5)
d.Update([]byte("a")); d.Update([]byte("bc")) // streaming == one-shot

bb := digest.BubbleBabble([]byte("Pineapple")) // "xigak-nyryk-..."
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`digest`](https://github.com/go-ruby-digest/digest) | the library — Ruby's Digest suite in pure Go |
| [`docs`](https://github.com/go-ruby-digest/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-digest.github.io`](https://github.com/go-ruby-digest/go-ruby-digest.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-digest/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** Extracted from rbgo's internals; no dependency on
  the Ruby runtime — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-digest/digest](https://github.com/go-ruby-digest/digest).
