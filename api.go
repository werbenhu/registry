// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package registry

// 路由服务器提供查询服务的接口
type Api interface {
	Start(addr string) error
	Stop()
}
