package srouter

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
	TagGroup    = "group"
	TagService  = "service"
	TagReplicas = "replicas"
)

type Serf struct {
	events  chan serf.Event
	member  *Member
	serf    *serf.Serf
	handler Handler
	members sync.Map
}

func NewSerf(local *Member) *Serf {
	s := &Serf{
		events: make(chan serf.Event, 3),
		member: local,
	}
	return s
}

func (s *Serf) LocalMember() *Member {
	node, ok := s.members.Load(s.member.Id)
	if !ok {
		return nil
	}
	return node.(*Member)
}

func (s *Serf) Members() []*Member {
	nodes := make([]*Member, 0)
	s.members.Range(func(key any, val any) bool {
		nodes = append(nodes, val.(*Member))
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
	cfg.MemberlistConfig.AdvertiseAddr, cfg.MemberlistConfig.AdvertisePort = s.splitHostPort(s.member.Advertise)
	cfg.MemberlistConfig.BindAddr, cfg.MemberlistConfig.BindPort = s.splitHostPort(s.member.Addr)
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

	go s.Loop()
	log.Printf("[INFO] serf discovery started, current member addr:%s, advertise addr:%s\n", s.member.Addr, s.member.Advertise)
	if len(s.member.Routers) > 0 {
		members := strings.Split(s.member.Routers, ",")
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
