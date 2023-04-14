// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package registry

// 自动发现事件通知接口
type Handler interface {
	// 当有新服务注册过来的时候会触发
	OnMemberJoin(*Member) error

	// 当有新服务离开的时候会触发
	OnMemberLeave(*Member) error

	// 当有新服务更新的时候会触发
	OnMemberUpdate(*Member) error
}

// 自动发现接口
type Discovery interface {
	// 设置服务发现事件处理接口
	SetHandler(Handler)

	//获取所有服务列表
	Members() []*Member

	//获取当前自身服务
	LocalMember() *Member

	//启动发现服务
	Start() error

	//停止服务
	Stop()
}
