// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package registry

// Err represents a custom error type with an error message and error code.
type Err struct {
	Msg  string // the error message
	Code int    // the error code
}

// String returns the error message as a string.
func (e Err) String() string {
	return e.Msg
}

// Error returns the error message as a string (implements the error interface).
func (e Err) Error() string {
	return e.Msg
}

// Pre-defined error instances with specific error codes and messages.
var (
	ErrMemberIdEmpty       = Err{Code: 10000, Msg: "id can't be empty"}
	ErrReplicasParam       = Err{Code: 10000, Msg: "member replicas param error"}
	ErrGroupNameEmpty      = Err{Code: 10001, Msg: "member group name empty"}
	ErrParseAddrToHostPort = Err{Code: 10002, Msg: "parse addr to host and port error"}
	ErrParsePort           = Err{Code: 10003, Msg: "parse port error"}
)
