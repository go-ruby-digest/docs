# Roadmap

`go-ruby-digest/digest` is grown **test-first**, each capability differential-tested against MRI
rather than built in isolation. Ruby's Digest message-digest suite — the
deterministic, interpreter-independent slice extracted from rbgo's internals — is
**complete**.

| Stage | What | Status |
| --- | --- | --- |
| Digest suite | MD5, SHA1, the SHA2 family (SHA256/384/512) and RIPEMD-160, each producing the same `hexdigest` / `digest` / `base64digest` output as reference Ruby. | **Done** |
| Digest::Instance protocol | The streaming `update` / `<<` / `reset` / `finish` interface and the one-shot `Digest::MD5.hexdigest(s)` class methods, mirroring MRI's protocol. | **Done** |
| BubbleBabble | `Digest.bubblebabble` / `hexencode` producing the human-readable BubbleBabble encoding byte-for-byte as Ruby's `digest/bubblebabble` does. | **Done** |
| Hex & base64 forms | `hexdigest`, `base64digest` and raw `digest`, with MRI's exact byte and character output for each. | **Done** |
| Streaming equivalence | Feeding input incrementally via `update` yields the same digest as a single one-shot call, the property reference Ruby guarantees. | **Done** |
| Differential oracle & coverage | A wide input corpus hashed both here and by the system `ruby`/`digest`, compared byte-for-byte; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic algorithm; it
  never runs arbitrary Ruby. Anything that needs a live binding or evaluation is
  the consumer's job — that is why `rbgo` binds this module rather than the
  reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's
  behaviour; differences across MRI releases are matched to the reference used by
  the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime;
  the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
