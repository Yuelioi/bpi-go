// Package client implements the shared HTTP, session, signing, response, and
// request-policy module used by every bpi-go domain. Most applications should
// construct the root bpi.Client; direct use is intended for composing a single
// domain module or building advanced custom requests.
package client
