// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu
package client

import (
	"context"
	"time"

	"github.com/werbenhu/registry"
	"google.golang.org/grpc"
)

// Grpc客户端的封装，方便客户端做一致性哈希查询
// 客户端如果需要查询单个key的所在的服务，直接使用此对象就行
type RpcClient struct {
	// 路由服务器的Grpc地址
	Addr string

	// grpc客户端连接对象
	conn *grpc.ClientConn

	// grpc客户端对象
	router registry.RouterClient
}

// 创建一个Grpc客户端对象
// 该对象会连接路由服务器的grpc服务端
// 如果连接失败，会返回错误
func NewRpcClient(router string) (*RpcClient, error) {
	client := &RpcClient{Addr: router}
	// 连接到grpc服务端的地址
	conn, err := grpc.Dial(client.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client.conn = conn
	client.router = registry.NewRouterClient(conn)
	return client, nil
}

// 关闭grpc客户端连接对象
func (c *RpcClient) Close() {
	c.conn.Close()
}

// 匹配某个key对应的一致性哈希服务
// group是组名，比如有3个mysql服务器同属一个db组，3个web服务器同属一个web组
func (c *RpcClient) Match(group string, key string) (*registry.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service, err := c.router.Match(ctx, &registry.MatchRequest{
		Group: group,
		Key:   key,
	})
	if err != nil {
		return nil, err
	}
	//服务包含3个属性，服务ID、服务所属的组名以及服务的地址
	return registry.NewService(service.Id, service.Group, service.Addr), nil
}

// 列出某个组里所有的服务
func (c *RpcClient) Members(group string) ([]*registry.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	services := make([]*registry.Service, 0)
	members, err := c.router.Members(ctx, &registry.MembersRequest{
		Group: group,
	})
	if err != nil {
		return services, err
	}

	for _, member := range members.Services {
		services = append(services, registry.NewService(member.Id, member.Group, member.Addr))
	}
	return services, nil
}
