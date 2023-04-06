package discovery

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

type Serf struct {
	events  chan serf.Event
	agent   *Agent
	serf    *serf.Serf
	handler Handler
	agents  sync.Map
}

func NewSerf(agent *Agent) *Serf {
	s := &Serf{
		events: make(chan serf.Event, 3),
		agent:  agent,
	}
	return s
}

func (s *Serf) LocalAgent() *Agent {
	node, ok := s.agents.Load(s.agent.Id)
	if !ok {
		return nil
	}
	return node.(*Agent)
}

func (s *Serf) Agents() []*Agent {
	nodes := make([]*Agent, 0)
	s.agents.Range(func(key any, val any) bool {
		nodes = append(nodes, val.(*Agent))
		return true
	})
	return nodes
}

func (s *Serf) SetHandler(h Handler) {
	s.handler = h
}

func (s *Serf) Stop() {
	s.serf.Shutdown()
	close(s.events)
}

func (s *Serf) Start() error {
	var err error
	cfg := serf.DefaultConfig()
	cfg.MemberlistConfig.AdvertiseAddr, cfg.MemberlistConfig.AdvertisePort = s.splitHostPort(s.agent.Advertise)
	cfg.MemberlistConfig.BindAddr, cfg.MemberlistConfig.BindPort = s.splitHostPort(s.agent.Addr)
	cfg.EventCh = s.events

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("ERROR"),
		Writer: io.MultiWriter(&lumberjack.Logger{
			Filename:   "./log/serf.log",
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		}, os.Stderr),
	}

	cfg.Logger = log.New(os.Stderr, "", log.LstdFlags)
	cfg.Logger.SetOutput(filter)
	cfg.MemberlistConfig.Logger = cfg.Logger
	cfg.NodeName = s.agent.Id
	cfg.Tags = s.agent.GetTags()

	s.serf, err = serf.Create(cfg)
	if err != nil {
		return err
	}

	go s.Loop()
	log.Printf("[INFO] serf discovery started, current agent addr:%s, advertise addr:%s\n", s.agent.Addr, s.agent.Advertise)
	if len(s.agent.Members) > 0 {
		members := strings.Split(s.agent.Members, ",")
		s.Join(members)
	}
	return nil
}

func (s *Serf) Join(members []string) error {
	_, err := s.serf.Join(members, true)
	return err
}

func (s *Serf) splitHostPort(addr string) (string, int) {
	h, p, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatalf("[ERROR] serf discovery parse addr:%s err:%s", addr, err.Error())
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatalf("[ERROR] serf discovery parse port:%s err:%s", p, err.Error())
	}
	return h, port
}

func (s *Serf) Loop() {
	for e := range s.events {
		switch e.EventType() {
		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				addr := fmt.Sprintf("%s:%d", member.Addr, member.Port)
				latest := newSimpleAgent(member.Name, addr, addr)
				latest.SetTags(member.Tags)

				if s.handler != nil {
					if err := s.handler.OnAgentJoin(latest); err == nil {
						s.agents.Store(latest.Id, latest)
						continue
					} else {
						log.Printf("[ERROR] serf handle agent join err:%s\n", err.Error())
					}
				}
				s.agents.Store(latest.Id, latest)
			}

		case serf.EventMemberUpdate:
			for _, member := range e.(serf.MemberEvent).Members {
				addr := fmt.Sprintf("%s:%d", member.Addr, member.Port)
				latest := newSimpleAgent(member.Name, addr, addr)
				latest.SetTags(member.Tags)

				if s.handler != nil {
					if err := s.handler.OnAgentUpdate(latest); err == nil {
						s.agents.Store(latest.Id, latest)
						continue
					} else {
						log.Printf("[ERROR] serf handle agent update err:%s\n", err.Error())
					}
				}
				s.agents.Store(latest.Id, latest)
			}

		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, member := range e.(serf.MemberEvent).Members {
				addr := fmt.Sprintf("%s:%d", member.Addr, member.Port)
				latest := newSimpleAgent(member.Name, addr, addr)
				latest.SetTags(member.Tags)

				s.agents.Delete(latest.Id)
				if s.handler != nil {
					if err := s.handler.OnAgentLeave(latest); err != nil {
						log.Printf("[ERROR] serf handle agent leave err:%s\n", err.Error())
					}
				}

			}
		}
	}
}
