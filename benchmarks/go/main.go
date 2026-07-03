// SPDX-License-Identifier: BSD-3-Clause
package main

import "github.com/go-ruby-digest/digest"

// hexOf is the exact Go analogue of Ruby's Digest::ALGO.hexdigest(data): the
// library's class one-shot HexSum, which the Ruby side calls verbatim. (The
// earlier New+Update+HexFinish spelling was a hand-inlined expansion of the same
// call; HexSum is the documented one-shot and the allocation-minimal hot path.)
func hexOf(name string, data []byte) string { s, _ := digest.HexSum(name, data); return s }

func main() {
	data := detBytes(4096)
	bench("sha256-4KiB", 1000, func() { sink = hexOf("SHA256", data) })
	bench("md5-4KiB", 1000, func() { sink = hexOf("MD5", data) })
	bench("sha512-4KiB", 1000, func() { sink = hexOf("SHA512", data) })
}
