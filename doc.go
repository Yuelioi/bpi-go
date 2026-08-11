// Package bpi provides the compatibility facade for an idiomatic Go client
// for Bilibili HTTP interfaces. Shared transport behavior lives in package
// client, while each public domain package owns its client implementation,
// parameters, models, and contract tests.
//
// Clients are independent and safe for concurrent use. Every domain network
// operation accepts a context, and cancellation propagates through the
// configured net/http transport. Construction does not read configuration
// files, install global loggers, mutate package state, or perform network I/O.
//
// Credentials must be supplied explicitly. The client scopes session Cookies
// to approved Bilibili hosts, keeps logging quiet by default, and removes
// sensitive query values from optional structured logs. Responses are bounded
// before decoding. A ResponseDecodeError retains an explicit private copy of a
// mismatched response body for local recovery, but its Error and formatting
// methods never expose that body.
package bpi
