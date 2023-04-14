package registry

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

type Member struct {
	sync.Mutex
	Id        string  `json:"id"`
	Addr      string  `json:"addr"`
	Advertise string  `json:"advertise"`
	Routers   string  `json:"-"`
	Replicas  string  `json:"replicas"`
	Service   Service `json:"service"`
	tags      map[string]string
}

func NewSimpleMember(id string, addr string, advertise string) *Member {
	return &Member{
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

func NewMember(id string, addr string, advertise string, routers string, group string, serviceAddr string) *Member {
	return &Member{
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

func (m *Member) IsSelf(b *Member) bool {
	return m.Id == b.Id
}

func (m *Member) SetTag(key string, val string) {
	m.Lock()
	defer m.Unlock()
	if m.tags == nil {
		m.tags = make(map[string]string)
	}
	m.tags[key] = val
}

func (m *Member) GetTag(key string) (string, bool) {
	m.Lock()
	defer m.Unlock()
	if m.tags == nil {
		return "", false
	}
	val, ok := m.tags[key]
	return val, ok
}

func (m *Member) SetTags(tags map[string]string) {
	if m.tags == nil {
		m.tags = make(map[string]string)
	}
	for k, v := range tags {
		m.SetTag(k, v)
	}

	m.Service.Group, _ = m.GetTag(TagGroup)
	m.Service.Addr, _ = m.GetTag(TagService)
	m.Replicas, _ = m.GetTag(TagReplicas)
}

func (m *Member) GetTags() map[string]string {
	m.SetTag(TagService, m.Service.Addr)
	m.SetTag(TagGroup, m.Service.Group)
	m.SetTag(TagReplicas, m.Replicas)

	m.Lock()
	defer m.Unlock()
	clone := make(map[string]string)
	for k, v := range m.tags {
		clone[k] = v
	}
	return clone
}

func (m *Member) Marshal() ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	return json.Marshal(m)
}

func (m *Member) Unmarshal(paylaod []byte) error {
	m.Lock()
	defer m.Unlock()
	return json.Unmarshal(paylaod, m)
}
