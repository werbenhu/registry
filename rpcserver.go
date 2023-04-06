package srouter

import (
	"context"
	"log"
	"net"

	"github.com/werbenhu/chash"
	"github.com/werbenhu/srouter/discovery"
	"google.golang.org/grpc"
)

type RpcServer struct {
	port string
	rpc  *grpc.Server
}

func NewRpcServer() *RpcServer {
	return &RpcServer{}
}

func (s *RpcServer) Match(ctx context.Context, req *MatchRequest) (*MatchResponse, error) {

	group, err := chash.GetGroup(req.Group)
	if err != nil {
		return nil, err
	}

	_, payload, err := group.Match(req.Key)
	if err != nil {
		return nil, err
	}

	agent := &discovery.Agent{}
	if err := agent.Unmarshal(payload); err != nil {
		return nil, err
	}

	return &MatchResponse{
		Id:    agent.Service.Id,
		Group: agent.Service.Group,
		Addr:  agent.Service.Addr,
	}, nil
}

func (s *RpcServer) Members(ctx context.Context, req *MembersRequest) (*MembersResponse, error) {

	group, err := chash.GetGroup(req.Group)
	if err != nil {
		return nil, err
	}

	elements := group.GetElements()
	services := make([]*MatchResponse, 0)
	for _, element := range elements {
		agent := &discovery.Agent{}
		if err := agent.Unmarshal(element.Payload); err == nil {
			service := &MatchResponse{
				Id:    agent.Service.Id,
				Group: agent.Service.Group,
				Addr:  agent.Service.Addr,
			}
			services = append(services, service)
		}
	}

	return &MembersResponse{
		Services: services,
	}, nil
}

func (s *RpcServer) Start(port string) error {
	var err error

	s.port = port
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Fatalf("[ERROR] rpc listen to port:%s failed, err:%s", s.port, err.Error())
		return err
	}

	s.rpc = grpc.NewServer()
	RegisterRouterServer(s.rpc, s)
	return s.rpc.Serve(listener)
}

func (s *RpcServer) Stop() {
	if s.rpc != nil {
		s.rpc.Stop()
	}
}
