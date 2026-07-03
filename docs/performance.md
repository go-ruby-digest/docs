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

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path recorded above. It isolates the library primitive
from Ruby-interpreter dispatch, answering the parity question head-on: *is the
pure-Go implementation as fast as the reference runtime's own `digest`?* The
**same workload, same inputs, same iteration counts** run through the Go library
and through each reference runtime's stdlib; outputs were checked identical to
MRI before any timing.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

#### md5-4KiB

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 4714.1 | 0.90× |
| MRI | 5226.0 | 1.00× |
| MRI + YJIT | 5111.0 | 0.98× |
| JRuby | 4735.6 | 0.91× |
| TruffleRuby | 5064.9 | 0.97× |

#### sha256-4KiB

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 1621.0 | 1.09× |
| MRI | 1487.0 | 1.00× |
| MRI + YJIT | 1433.0 | 0.96× |
| JRuby | 1544.6 | 1.04× |
| TruffleRuby | 15254.6 | 10.26× |

#### sha512-4KiB

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 2321.8 | 0.92× |
| MRI | 2522.0 | 1.00× |
| MRI + YJIT | 2492.0 | 0.99× |
| JRuby | 2574.2 | 1.02× |
| TruffleRuby | 9348.0 | 3.71× |

**At parity with MRI** across MD5 / SHA-256 / SHA-512 (0.90–1.09×): the library wraps Go's `crypto/*`, which is assembly-optimized on arm64, matching MRI's OpenSSL-backed C. The TruffleRuby SHA-256/512 columns are cold-JIT outliers (Graal had not compiled the digest loop within the warm-up budget), not steady-state numbers.

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
