# Performance

`go-ruby-digest/digest` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `digest`. This
page records a **comparative benchmark** of that module against the reference
Ruby runtimes, part of the ecosystem-wide per-module parity suite.

## What is measured

The **same** Ruby script — `Digest::MD5`/`SHA1`/`SHA256` `hexdigest` loops over a message — is run under every runtime. `rbgo`'s
number reflects **this pure-Go library doing the work**; every other column is
that interpreter's own `digest` stdlib. So the comparison is the **Ruby-visible
operation**, apples-to-apples across interpreters. The script prints a
deterministic checksum and its output is checked **byte-identical to MRI**
before timing.

- **Host:** Apple M4 Max, macOS (darwin/arm64). **Method:** best-of-5 wall time
  (best, not mean, to suppress scheduler noise); single-shot processes, no
  warm-up beyond the script's own loop.
- **Runtimes:** `ruby 4.0.5 +PRISM` (MRI, the oracle) and `ruby --yjit`;
  `jruby 10.1.0.0` (OpenJDK 25); `truffleruby 34.0.1` (GraalVM CE Native).
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules)
  (`digest.rb` + `run.sh`). Reproduce:
  `RBGO=./rbgo TRUFFLE=truffleruby bash bench/modules/run.sh 5`.

## Result (best of 5, ms)

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-digest) | 370 | 1.09× |
| MRI (ruby 4.0.5) | 340 | 1.00× |
| MRI + YJIT | 330 | 0.97× |
| JRuby 10.1.0.0 | 1330 | 3.91× |
| TruffleRuby 34.0.1 | 480 | 1.41× |

rbgo runs on **go-ruby-digest** at ~1.1x MRI — effective parity, as both ultimately drive optimized native hash implementations.

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. These are real measured numbers from the
    2026-06-29 run — nothing is cherry-picked.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03 (post-optimization)

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path recorded above. It isolates the library primitive
from Ruby-interpreter dispatch, answering the parity question head-on: *is the
pure-Go implementation as fast as the reference runtime's own `digest`?* The
**same workload, same inputs, same iteration counts** run through the Go library
and through each reference runtime's stdlib; outputs were checked identical to
MRI before any timing. The Go driver calls the library's class one-shot
`digest.HexSum(...)` — the exact analogue of Ruby's `Digest::SHA256.hexdigest`.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` / `vs YJIT` < 1.00× means *faster than
  that runtime*. Interpreter start-up is outside the timed region, so these are
  operation costs, not `ruby file.rb` process costs.

#### md5-4KiB

| Runtime | ns/op | vs MRI | vs YJIT |
| --- | ---: | ---: | ---: |
| **go-ruby (pure Go)** | 5024.3 | 0.89× | **0.90×** |
| MRI | 5618.0 | 1.00× | — |
| MRI + YJIT | 5572.0 | 0.99× | 1.00× |
| JRuby | 4917.5 | 0.88× | — |
| TruffleRuby | 4978.4 | 0.89× | — |

#### sha256-4KiB

| Runtime | ns/op | vs MRI | vs YJIT |
| --- | ---: | ---: | ---: |
| **go-ruby (pure Go)** | 1452.0 | 0.87× | **0.89×** |
| MRI | 1676.0 | 1.00× | — |
| MRI + YJIT | 1640.0 | 0.98× | 1.00× |
| JRuby | 1578.3 | 0.94× | — |
| TruffleRuby | 15502.2 | 9.25× | — |

#### sha512-4KiB

| Runtime | ns/op | vs MRI | vs YJIT |
| --- | ---: | ---: | ---: |
| **go-ruby (pure Go)** | 2646.5 | 0.97× | 1.02× |
| MRI | 2723.0 | 1.00× | — |
| MRI + YJIT | 2584.0 | 0.95× | 1.00× |
| JRuby | 2681.3 | 0.98× | — |
| TruffleRuby | 9448.3 | 3.47× | — |

**SHA-256 now beats YJIT (0.89×) and MRI (0.87×)** — and MD5 with it (0.90× /
0.89×); SHA-512 sits at the shared floor (0.97× MRI, a run-to-run tie with
YJIT). This is the payoff of an allocation rework, not a faster hash: the
one-shots were rebuilding a `hash.Hash`, summing into a fresh slice, and
double-allocating in `EncodeToString` — **five heap allocations** per
`hexdigest`. The new path (`fastSum` + stack-buffered hex/Base64 encode) sums
into a stack array and encodes into a stack buffer, leaving the returned string
as the **sole allocation**:

| library op | before | after |
| --- | --- | --- |
| `HexSum` SHA-256 | 1511 ns · 320 B · 5 allocs | **1463 ns · 64 B · 1 alloc** |
| `HexSum` SHA-512 | 2724 ns · 576 B · 5 allocs | **2642 ns · 128 B · 1 alloc** |
| `HexSum` MD5 | 5183 ns · 176 B · 4 allocs | **5150 ns · 32 B · 1 alloc** |

*(Go `-benchmem`, best of 8, same host.)*

!!! note "The honest floor: it's a shared hardware core"
    The SHA-256 hash **core** is the same silicon on both sides — Go's
    `crypto/sha256` and MRI's OpenSSL both issue the ARMv8 `SHA256H`/`SHA256H2`
    instructions, so the pure compute (~1400 ns for 4 KiB) is **irreducible and
    equal**; neither runtime can beat the other on the hash itself. The margin
    that remains is entirely **wrapper overhead** around that core: MRI/YJIT pay
    a Ruby method dispatch, an OpenSSL `EVP` context, and a result-string
    allocation on every `hexdigest`; go-ruby now pays one string allocation and a
    concrete function call. Killing four of five allocations shrank go-ruby's
    wrapper to ~50 ns over the core versus the interpreters' ~200 ns, which is
    why go-ruby lands below YJIT rather than the 1.09× it showed before this
    change. SHA-256 was measured across six back-to-back harness runs: go-ruby
    held **1438–1470 ns** (tight) while YJIT ranged **1427–1640 ns** (noisy, one
    low outlier at 1427); go-ruby's central tendency is at or below YJIT's every
    run but the two share a hardware floor — read this as *go-ruby is now as fast
    as YJIT, usually a touch faster*, not as a core-speed victory that isn't
    physically possible.

    The TruffleRuby SHA-256/512 columns are cold-JIT outliers (Graal had not
    compiled the digest loop within the warm-up budget), not steady-state numbers.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-digest/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via
    `go.mod`), the equivalent `ruby/digest.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — most visibly TruffleRuby on the shortest loops (a few cold-JIT
    outliers are noted in the text). Sub-microsecond rows carry the most relative
    noise; treat those ratios as order-of-magnitude. Every number here is a
    **real measured value** from the dated run above — nothing is fabricated,
    estimated, or cherry-picked. The go-ruby column is the pure-Go library; every
    other column is that interpreter's own stdlib doing the equivalent work.
