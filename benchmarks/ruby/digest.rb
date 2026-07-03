# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
require "digest"
require_relative "_harness"
data = det_bytes(4096)
bench("sha256-4KiB", 1000) { Digest::SHA256.hexdigest(data) }
bench("md5-4KiB",    1000) { Digest::MD5.hexdigest(data) }
bench("sha512-4KiB", 1000) { Digest::SHA512.hexdigest(data) }
