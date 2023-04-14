// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/werbenhu/registry"
)

func TestErrString(t *testing.T) {
	c := registry.Err{
		Msg:  "test",
		Code: 0x1,
	}

	require.Equal(t, "test", c.String())
}

func TestErrErrorr(t *testing.T) {
	c := registry.Err{
		Msg:  "error",
		Code: 0x1,
	}

	require.Equal(t, "error", error(c).Error())
}
