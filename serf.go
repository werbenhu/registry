package registry

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/logutils"
	"github.com/hashicorp/serf/serf"
	"github.com/natefinch/lumberjack"
)

const (
	// tag key of group name
	TagGroup = "group"

	// tag key of service address
	TagAddr = "addr"

	// tag key of replicas
	TagReplicas = "replicas"
)

type Serf struct {

	// Event is a generic interface for exposing Serf events
	// Clients will usually need to use a type switches to get
	// to a more useful type
	events chan serf.Event

	// local member of current registry server
	member *Member

	// Serf is a single node that is part of a single cluster that gets
	// events about joins/leaves/failures/etc. It is created with the Create
	// method.
	serf *serf.Serf

	// Auto-discover event notification interface
	handler Handler

	// Members of all services
	members sync.Map
}

// NewSerf create a discovery instance of hashicorp/serf
func NewSerf(local *Member) *Serf {
	s := &Serf{
		member: local,
		events: make(chan serf.Event),
	}
	return s
}

// LocalMember get current registry service
func (s *Serf) LocalMember() *Member {
	node, ok := s.members.Load(s.member.Id)
	if !ok {
		return nil
	}
	return node.(*Member)
}

// Members get members of all services
func (s *Serf) Members() []*Member {
	nodes := make([]*Member, 0)
	s.members.Range(func(key any, val any) bool {
		nodes = append(nodes, val.(*Member))
		return true
	})
	return nodes
}

// Set event processing handler when new services are discovered
func (s *Serf) SetHandler(h Handler) {
	s.handler = h
}

func (s *Serf) Stop() {
	if s.serf != nil {
		s.serf.Shutdown()
	}
	if s.events != nil {
		close(s.events)
		s.events = nil
	}
}

// Start hashicorp/serf agent
func (s *Serf) Start() error {
	var err error
	var host string
	var port int
	cfg := serf.DefaultConfig()

	host, port, err = s.splitHostPort(s.member.Advertise)
	if err != nil {
		return err
	}
	cfg.MemberlistConfig.AdvertiseAddr = host
	cfg.MemberlistConfig.AdvertisePort = port

	host, port, err = s.splitHostPort(s.member.Bind)
	if err != nil {
		return err
	}
	cfg.MemberlistConfig.BindAddr = host
	cfg.MemberlistConfig.BindPort = port
	cfg.EventCh = s.events

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("ERROR"),
		Writer: io.MultiWriter(&lumberjack.Logger{
			Filename:   "./log/serf.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     28,
		}, os.Stderr),
	}

	cfg.Logger = log.New(os.Stderr, "", log.LstdFlags)
	cfg.Logger.SetOutput(filter)
	cfg.MemberlistConfig.Logger = cfg.Logger
	cfg.NodeName = s.member.Id
	cfg.Tags = s.member.GetTags()

	s.serf, err = serf.Create(cfg)
	if err != nil {
		return err
	}

	s.members.Store(s.member.Id, s.member)
	go s.loop()
	log.Printf("[INFO] serf discovery started, current service bind:%s, advertise addr:%s\n", s.member.Bind, s.member.Advertise)
	if len(s.member.Registries) > 0 {
		members := strings.Split(s.member.Registries, ",")
		s.Join(members)
	}
	return nil
}

// Join joins an existing Serf cluster.
func (s *Serf) Join(members []string) error {
	_, err := s.serf.Join(members, true)
	return err
}

// Split address into host names and ports
func (s *Serf) splitHostPort(addr string) (string, int, error) {
	h, p, err := net.SplitHostPort(addr)
	if err != nil {
		return "", -1, ErrParseAddrToHostPort
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		return "", -1, ErrParsePort
	}
	return h, port, nil
}

// Loop read the exposing Serf events and pass events to the handler
func (s *Serf) loop() {
	for e := range s.events {
		switch e.EventType() {
		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				addr := fmt.Sprintf("%s:%d", member.Addr, member.Port)
				latest := NewSimpleMember(member.Name, addr, addr)
				latest.SetTags(member.Tags)

				if s.handler != nil {
					if err := s.handler.OnMemberJoin(latest); err == nil {
						s.members.Store(latest.Id, latest)
						continue
					} else {
						log.Printf("[ERROR] serf handle member join err:%s\n", err.Error())
					}
				}
				s.members.Store(latest.Id, latest)
			}

		case serf.EventMemberUpdate:
			for _, member := range e.(serf.MemberEvent).Members {
				addr := fmt.Sprintf("%s:%d", member.Addr, member.Port)
				latest := NewSimpleMember(member.Name, addr, addr)
				latest.SetTags(member.Tags)

				if s.handler != nil {
					if err := s.handler.OnMemberUpdate(latest); err == nil {
						s.members.Store(latest.Id, latest)
						continue
					} else {
						log.Printf("[ERROR] serf handle member update err:%s\n", err.Error())
					}
				}
				s.members.Store(latest.Id, latest)
			}

		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, member := range e.(serf.MemberEvent).Members {
				addr := fmt.Sprintf("%s:%d", member.Addr, member.Port)
				latest := NewSimpleMember(member.Name, addr, addr)
				latest.SetTags(member.Tags)

				s.members.Delete(latest.Id)
				if s.handler != nil {
					if err := s.handler.OnMemberLeave(latest); err != nil {
						log.Printf("[ERROR] serf handle member leave err:%s\n", err.Error())
					}
				}
			}
		}
	}
}
