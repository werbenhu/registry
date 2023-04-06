package client

import (
	"context"
	"time"

	"github.com/werbenhu/srouter"
	"google.golang.org/grpc"
)

type Client struct {
	Addr string
	conn *grpc.ClientConn
	rpc  srouter.RouterClient
}

func New(router string) (*Client, error) {
	client := &Client{Addr: router}
	conn, err := grpc.Dial(client.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client.conn = conn
	client.rpc = srouter.NewRouterClient(conn)
	return client, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Match(group string, key string) (*srouter.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service, err := c.rpc.Match(ctx, &srouter.MatchRequest{
		Group: group,
		Key:   key,
	})
	if err != nil {
		return nil, err
	}
	return srouter.NewService(service.Id, service.Group, service.Addr), nil
}

func (c *Client) Members(group string) ([]*srouter.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	services := make([]*srouter.Service, 0)

	members, err := c.rpc.Members(ctx, &srouter.MembersRequest{
		Group: group,
	})
	if err != nil {
		return services, err
	}
	for _, member := range members.Services {
		services = append(services, srouter.NewService(member.Id, member.Group, member.Addr))
	}

	return services, nil
}
