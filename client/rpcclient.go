// Package client provides a gRPC client for service discovery.

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

// RpcClient is a gRPC client for service discovery.
type RpcClient struct {
	// Addr is the registry server address.
	Addr string

	// conn is the gRPC connection.
	conn *grpc.ClientConn

	// reg is the gRPC client object.
	reg registry.RClient
}

// NewRpcClient creates a new RpcClient object and connects to the registry server at `addr`.
func NewRpcClient(addr string) (*RpcClient, error) {
	client := &RpcClient{Addr: addr}
	// Connect to the registry server.
	conn, err := grpc.Dial(client.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client.conn = conn
	client.reg = registry.NewRClient(conn)
	return client, nil
}

// Close closes the gRPC client connection.
func (c *RpcClient) Close() {
	c.conn.Close()
}

// Match assigns a service to a key using the consistent hashing algorithm.
//
// Parameters:
// - group: The group name of the services.
// - key: The key, such as user ID, device ID, etc.
//
// Returns:
// - The service that matches the key.
// - An error if the service cannot be found.
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
	// The service contains three attributes: service ID, group name, and service address.
	return registry.NewService(service.Id, service.Group, service.Addr), nil
}

// Members returns the list of services in a group.
//
// Parameters:
// - group: The group name of the services.
//
// Returns:
// - The list of services in the group.
// - An error if the group does not exist or cannot be accessed.
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
