package srouter

import (
	"context"
	"net"

	"github.com/werbenhu/chash"
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

	m := &Member{}
	if err := m.Unmarshal(payload); err != nil {
		return nil, err
	}

	return &MatchResponse{
		Id:    m.Service.Id,
		Group: m.Service.Group,
		Addr:  m.Service.Addr,
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
		m := &Member{}
		if err := m.Unmarshal(element.Payload); err == nil {
			service := &MatchResponse{
				Id:    m.Service.Id,
				Group: m.Service.Group,
				Addr:  m.Service.Addr,
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
