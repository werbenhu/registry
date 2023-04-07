// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package register

import "github.com/werbenhu/srouter"

// 注册器：对注册的封装，服务提供者直接使用注册器就可以很容易的注册服务到路由服务器
type Register struct {
	serf    srouter.Discovery
	handler srouter.Handler
	member  *srouter.Member
}

// 新建一个注册
// id： 是需要注册的服务的ID
// addr： 当前服务跟路由服务器通信的地址
// advertise: 对外公布的服务发现通信的地址
// routers: 路由服务器的地址，多个用逗号隔开
// group: 服务所属的组
// service: 服务提供服务的地址
func New(id string, addr string, advertise string, routers string, group string, service string) *Register {
	member := srouter.NewMember(id, addr, advertise, routers, group, service)
	return &Register{member: member}
}

// hander将允许服务监听注册服务器收到的新注册、更新以及删除注册事件
func (r *Register) SetHandler(h srouter.Handler) {
	r.handler = h
}

// Run将注册本服务到路由服务器，并保持双方的通信
func (r *Register) Run() error {
	r.serf = srouter.NewSerf(r.member)
	r.serf.SetHandler(r.handler)
	return r.serf.Start()
}

// 关闭注册器
func (r *Register) Stop() {
	r.serf.Stop()
}
