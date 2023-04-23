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
	// TagGroup is the tag key of the group name.
	TagGroup = "group"

	// TagAddr is the tag key of the service address.
	TagAddr = "addr"

	// TagReplicas is the tag key of replicas.
	TagReplicas = "replicas"
)

// Serf represents a discovery instance of hashicorp/serf.
type Serf struct {
	events  chan serf.Event // A channel to expose Serf events.
	member  *Member         // The local member of the current registry server.
	serf    *serf.Serf      // A single node that is part of a single cluster that gets events about joins/leaves/failures/etc.
	handler Handler         // An auto-discover event notification interface.
	members sync.Map        // The members of all services.
}

// NewSerf creates a new instance of Serf.
func NewSerf(local *Member) *Serf {
	s := &Serf{
		member: local,
	}
	return s
}

// LocalMember returns the current registry service.
func (s *Serf) LocalMember() *Member {
	node, ok := s.members.Load(s.member.Id)
	if !ok {
		return nil
	}
	return node.(*Member)
}

// Members returns the members of all services.
func (s *Serf) Members() []*Member {
	nodes := make([]*Member, 0)
	s.members.Range(func(key any, val any) bool {
		nodes = append(nodes, val.(*Member))
		return true
	})
	return nodes
}

// SetHandler sets the event processing handler when new services are discovered.
func (s *Serf) SetHandler(h Handler) {
	s.handler = h
}

// Stop stops the Serf server.
func (s *Serf) Stop() {
	if s.serf != nil {
		s.serf.Shutdown()
	}
	if s.events != nil {
		close(s.events)
	}
	log.Printf("[DEBUG] serf server stopped.\n")
}

// Start starts the HashiCorp Serf agent with the configuration provided in s.
func (s *Serf) Start() error {
	// Initialize variables.
	var err error
	var host string
	var port int
	cfg := serf.DefaultConfig()
	s.events = make(chan serf.Event)

	// Extract host and port from Advertise address and set them in the configuration.
	host, port, err = s.splitHostPort(s.member.Advertise)
	if err != nil {
		return err
	}
	cfg.MemberlistConfig.AdvertiseAddr = host
	cfg.MemberlistConfig.AdvertisePort = port

	// Extract host and port from Bind address and set them in the configuration.
	host, port, err = s.splitHostPort(s.member.Bind)
	if err != nil {
		return err
	}
	cfg.MemberlistConfig.BindAddr = host
	cfg.MemberlistConfig.BindPort = port
	cfg.EventCh = s.events

	// Set up the logger for Serf and the memberlist package.
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

	// Set the node name and tags in the configuration.
	cfg.NodeName = s.member.Id
	cfg.Tags = s.member.GetTags()

	// Create the Serf agent with the configuration.
	s.serf, err = serf.Create(cfg)
	if err != nil {
		return err
	}

	// Store the member in the members map and start the loop.
	s.members.Store(s.member.Id, s.member)
	go s.loop()

	// Print the bind and advertise addresses to the log.
	log.Printf("[INFO] Serf discovery started, current service bind:%s, advertise addr:%s\n", s.member.Bind, s.member.Advertise)

	// Join any registries that were specified in the member's configuration.
	if len(s.member.Registries) > 0 {
		members := strings.Split(s.member.Registries, ",")
		s.Join(members)
	}
	return nil
}

// Join joins the Serf agent to an existing Serf cluster with the specified members.
func (s *Serf) Join(members []string) error {
	_, err := s.serf.Join(members, true)
	return err
}

// splitHostPort splits an address of the form "host:port" into separate host and port strings.
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

// loop reads exposing Serf events and passes events to the handler
func (s *Serf) loop() {
	for e := range s.events {
		switch e.EventType() {
		// handle member join event
		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				addr := fmt.Sprintf("%s:%d", member.Addr, member.Port)
				latest := NewSimpleMember(member.Name, addr, addr)
				latest.SetTags(member.Tags)

				// call handler's OnMemberJoin method and store member
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

		// handle member update event
		case serf.EventMemberUpdate:
			for _, member := range e.(serf.MemberEvent).Members {
				addr := fmt.Sprintf("%s:%d", member.Addr, member.Port)
				latest := NewSimpleMember(member.Name, addr, addr)
				latest.SetTags(member.Tags)

				// call handler's OnMemberUpdate method and store member
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

		// handle member leave or failed event
		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, member := range e.(serf.MemberEvent).Members {
				addr := fmt.Sprintf("%s:%d", member.Addr, member.Port)
				latest := NewSimpleMember(member.Name, addr, addr)
				latest.SetTags(member.Tags)

				// delete member and call handler's OnMemberLeave method if it exists
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
