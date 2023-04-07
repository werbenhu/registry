// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package srouter

type err struct {
	Msg  string
	Code int
}

func (e err) String() string {
	return e.Msg
}

func (e err) Error() string {
	return e.Msg
}

var (
	ErrMemberIdEmpty  = err{Code: 10000, Msg: "id can't be empty"}
	ErrReplicasParam  = err{Code: 10000, Msg: "member replicas param error"}
	ErrGroupNameEmpty = err{Code: 10001, Msg: "member group name empty"}
)
