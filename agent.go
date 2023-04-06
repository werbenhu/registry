package srouter

import (
	"encoding/json"
	"sync"
)

type Service struct {
	Id    string `json:"id"`
	Group string `json:"group"`
	Addr  string `json:"addr"`
}

func NewService(id string, group string, addr string) *Service {
	return &Service{
		Id:    id,
		Group: group,
		Addr:  addr,
	}
}

type Agent struct {
	sync.Mutex
	Id        string  `json:"id"`
	Addr      string  `json:"addr"`
	Advertise string  `json:"advertise"`
	Routers   string  `json:"-"`
	Replicas  string  `json:"replicas"`
	Service   Service `json:"service"`
	tags      map[string]string
}

func newSimpleAgent(id string, addr string, advertise string) *Agent {
	return &Agent{
		Id:        id,
		Addr:      addr,
		Advertise: advertise,
		Routers:   "",
		Replicas:  DefaultReplicas,
		Service: Service{
			Id: id,
		},
	}
}

func NewAgent(id string, addr string, advertise string, routers string, group string, serviceAddr string) *Agent {
	return &Agent{
		Id:        id,
		Addr:      addr,
		Advertise: advertise,
		Routers:   routers,
		Replicas:  DefaultReplicas,
		Service: Service{
			Id:    id,
			Group: group,
			Addr:  serviceAddr,
		},
	}
}

func (a *Agent) IsSelf(b *Agent) bool {
	return a.Id == b.Id
}

func (a *Agent) SetTag(key string, val string) {
	a.Lock()
	defer a.Unlock()
	if a.tags == nil {
		a.tags = make(map[string]string)
	}
	a.tags[key] = val
}

func (a *Agent) GetTag(key string) (string, bool) {
	a.Lock()
	defer a.Unlock()
	if a.tags == nil {
		return "", false
	}
	val, ok := a.tags[key]
	return val, ok
}

func (a *Agent) SetTags(tags map[string]string) {
	if a.tags == nil {
		a.tags = make(map[string]string)
	}
	for k, v := range tags {
		a.SetTag(k, v)
	}

	a.Service.Group, _ = a.GetTag(TagGroup)
	a.Service.Addr, _ = a.GetTag(TagService)
	a.Replicas, _ = a.GetTag(TagReplicas)
}

func (a *Agent) GetTags() map[string]string {
	a.SetTag(TagService, a.Service.Addr)
	a.SetTag(TagGroup, a.Service.Group)
	a.SetTag(TagReplicas, a.Replicas)

	a.Lock()
	defer a.Unlock()
	clone := make(map[string]string)
	for k, v := range a.tags {
		clone[k] = v
	}
	return clone
}

func (a *Agent) Marshal() ([]byte, error) {
	a.Lock()
	defer a.Unlock()
	return json.Marshal(a)
}

func (a *Agent) Unmarshal(paylaod []byte) error {
	a.Lock()
	defer a.Unlock()
	return json.Unmarshal(paylaod, a)
}
