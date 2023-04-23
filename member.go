package registry

import (
	"encoding/json"
	"sync"
)

// Service represents a service object.
type Service struct {
	// The ID of the service.
	Id string `json:"id"`

	// The group name of this service.
	Group string `json:"group"`

	// The service address provided to the client.
	Addr string `json:"addr"`
}

// NewService creates a new service object.
func NewService(id string, group string, addr string) *Service {
	return &Service{
		Id:    id,
		Group: group,
		Addr:  addr,
	}
}

// Member is used for auto-discovery. When a service is discovered, a Member object is created.
type Member struct {
	sync.Mutex

	// The ID of the service.
	Id string `json:"id"`

	// The address used to register the service to the registry server.
	Bind string `json:"bind"`

	// The address that the service will advertise to the registry server.
	Advertise string `json:"advertise"`

	// The addresses of the registry servers. If there are more than one, separate them with commas, such as "192.168.1.101:7370,192.168.1.102:7370".
	Registries string `json:"-"`

	// The number of replicated elements of a service that need to be virtualized.
	Replicas string `json:"replicas"`

	// Service information.
	Service Service `json:"service"`

	// Tags for extra information.
	tags map[string]string
}

// NewSimpleMember creates a simple Member object. It does not contain the address of the service.
func NewSimpleMember(id string, bind string, advertise string) *Member {
	return &Member{
		Id:         id,
		Bind:       bind,
		Advertise:  advertise,
		Registries: "",
		Replicas:   DefaultReplicas,
		Service: Service{
			Id: id,
		},
	}
}

// NewMember creates a new Member object with the given attributes.
func NewMember(id string, bind string, advertise string, registries string, group string, addr string) *Member {
	return &Member{
		Id:         id,
		Bind:       bind,
		Advertise:  advertise,
		Registries: registries,
		Replicas:   DefaultReplicas,
		Service: Service{
			Id:    id,
			Group: group,
			Addr:  addr,
		},
	}
}

// IsSelf returns true if the given Member object has the same ID as this Member object.
func (m *Member) IsSelf(b *Member) bool {
	return m.Id == b.Id
}

// SetTag sets the extra information associated with the given tag for this Member object.
func (m *Member) SetTag(key string, val string) {
	m.Lock()
	defer m.Unlock()

	// If tags are not initialized, initialize them.
	if m.tags == nil {
		m.tags = make(map[string]string)
	}

	// Set the tag's value.
	m.tags[key] = val

	// Update service attributes based on specific tags.
	if key == TagAddr {
		m.Service.Addr = val
	} else if key == TagGroup {
		m.Service.Group = val
	} else if key == TagReplicas {
		m.Replicas = val
	}
}

// GetTag retrieves the value associated with the given tag for this Member object.
func (m *Member) GetTag(key string) (string, bool) {
	m.Lock()
	defer m.Unlock()

	// If tags are not initialized, return false.
	if m.tags == nil {
		return "", false
	}

	// Retrieve the tag's value.
	val, ok := m.tags[key]
	return val, ok
}

// SetTags set tags for this service
func (m *Member) SetTags(tags map[string]string) {
	if m.tags == nil {
		m.tags = make(map[string]string)
	}

	// Set each tag using SetTag method.
	for k, v := range tags {
		m.SetTag(k, v)
	}

	// Update service attributes based on specific tags.
	m.Service.Group, _ = m.GetTag(TagGroup)
	m.Service.Addr, _ = m.GetTag(TagAddr)
	m.Replicas, _ = m.GetTag(TagReplicas)
}

// GetTags retrieves all tags and their values for this Member object.
func (m *Member) GetTags() map[string]string {

	// Update specific tags before retrieving all tags.
	m.SetTag(TagAddr, m.Service.Addr)
	m.SetTag(TagGroup, m.Service.Group)
	m.SetTag(TagReplicas, m.Replicas)

	m.Lock()
	defer m.Unlock()

	// Clone the tags to avoid data races.
	clone := make(map[string]string)
	for k, v := range m.tags {
		clone[k] = v
	}
	return clone
}

// Marshal returns the JSON encoding of this Member object.
func (m *Member) Marshal() ([]byte, error) {
	m.Lock()
	defer m.Unlock()
	return json.Marshal(m)
}

// Unmarshal parses the given JSON-encoded data and stores
func (m *Member) Unmarshal(paylaod []byte) error {
	m.Lock()
	defer m.Unlock()
	return json.Unmarshal(paylaod, m)
}
