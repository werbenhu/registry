package registry

import (
	"context"
	"log"
	"net"

	"github.com/werbenhu/chash"
	"google.golang.org/grpc"
)

// RpcServer is a grpc server for client for service discovery
type RpcServer struct {
	addr string
	rpc  *grpc.Server
}

// NewRpcServer create a new RpcServer object
func NewRpcServer() *RpcServer {
	return &RpcServer{}
}

// Match assign a service to a key with consistent hashing algorithm
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

// Members get services list of a group
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

// Start running the grpc server
func (s *RpcServer) Start(addr string) error {
	var err error

	s.addr = addr
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.rpc = grpc.NewServer()
	RegisterRServer(s.rpc, s)
	return s.rpc.Serve(listener)
}

// Stop stop the grpc server
func (s *RpcServer) Stop() {
	if s.rpc != nil {
		s.rpc.Stop()
		log.Printf("[DEBUG] rpc server is stoped.\n")
	}
}
