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

// RpcClient is a grpc client for discovery
type RpcClient struct {
	// registry server address
	Addr string

	// grpc connection
	conn *grpc.ClientConn

	// grpc client object
	reg registry.RClient
}

// NewRpcClient create a new RpcClient object
func NewRpcClient(router string) (*RpcClient, error) {
	client := &RpcClient{Addr: router}
	// connecting to registry server
	conn, err := grpc.Dial(client.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client.conn = conn
	client.reg = registry.NewRClient(conn)
	return client, nil
}

// Close close the rpc client
func (c *RpcClient) Close() {
	c.conn.Close()
}

// Match assign a service to a key with consistent hashing algorithm
// groupName:
//
//	the group name of the services
//
// key:
//
//	the key such as user ID, device ID, etc
func (c *RpcClient) Match(group string, key string) (*registry.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service, err := c.reg.Match(ctx, &registry.MatchRequest{
		Group: group,
		Key:   key,
	})
	if err != nil {
		return nil, err
	}
	// The service contains three attributes: service ID, group name, and service address
	return registry.NewService(service.Id, service.Group, service.Addr), nil
}

// Members get services list of a group
// groupName:
//
//	the group name of the services
func (c *RpcClient) Members(group string) ([]*registry.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	services := make([]*registry.Service, 0)
	members, err := c.reg.Members(ctx, &registry.MembersRequest{
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
