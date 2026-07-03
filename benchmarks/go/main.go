// SPDX-License-Identifier: BSD-3-Clause
package main

import "github.com/go-ruby-digest/digest"

func hexOf(name string, data []byte) string {
	d := must(name)
	d.Update(data)
	return d.HexFinish()
}
func must(name string) digest.Digest { d, _ := digest.New(name); return d }

func main() {
	data := detBytes(4096)
	bench("sha256-4KiB", 1000, func() { sink = hexOf("SHA256", data) })
	bench("md5-4KiB", 1000, func() { sink = hexOf("MD5", data) })
	bench("sha512-4KiB", 1000, func() { sink = hexOf("SHA512", data) })
}
