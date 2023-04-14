// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package registry

type Err struct {
	Msg  string
	Code int
}

func (e Err) String() string {
	return e.Msg
}

func (e Err) Error() string {
	return e.Msg
}

var (
	ErrMemberIdEmpty       = Err{Code: 10000, Msg: "id can't be empty"}
	ErrReplicasParam       = Err{Code: 10000, Msg: "member replicas param error"}
	ErrGroupNameEmpty      = Err{Code: 10001, Msg: "member group name empty"}
	ErrParseAddrToHostPort = Err{Code: 10002, Msg: "parse addr to host and port error"}
	ErrParsePort           = Err{Code: 10003, Msg: "parse port error"}
)
