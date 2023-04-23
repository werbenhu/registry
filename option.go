package registry

import (
	"os"

	"github.com/rs/xid"
)

// Option represents the options for registry server.
type Option struct {

	// Id is the service ID.
	Id string

	// Bind is the address used to register the service.
	// If there is a firewall, ensure that the port is open for both TCP and UDP.
	Bind string

	// BindAdvertise is the address that the service will advertise to other services for registering.
	// Can be used for basic NAT traversal where both the internal IP:port and external IP:port are known.
	BindAdvertise string

	// Registries are the addresses of other registry servers.
	// If there are more than one, separate them with commas, such as "192.168.1.101:7370,192.168.1.102:7370".
	Registries string

	// Addr is the address used for service discovery.
	Addr string

	// Advertise is the address that will be advertised to clients for service discovery.
	Advertise string
}

// IOption represents a function that modifies the Option.
type IOption func(o *Option)

// OptId sets the service ID option.
func OptId(id string) IOption {
	return func(o *Option) {
		if id != "" {
			o.Id = id
		}
	}
}

// OptAddr sets the service discovery address option.
func OptAddr(addr string) IOption {
	return func(o *Option) {
		o.Addr = addr
	}
}

// OptAdvertise sets the advertised address for service discovery option.
func OptAdvertise(addr string) IOption {
	return func(o *Option) {
		o.Advertise = addr
	}
}

// OptBindAdvertise sets the advertised address for service registration option.
func OptBindAdvertise(addr string) IOption {
	return func(o *Option) {
		if addr != "" {
			o.BindAdvertise = addr
		}
	}
}

// OptBind sets the address used for service registration option.
func OptBind(addr string) IOption {
	return func(o *Option) {
		if addr != "" {
			o.Bind = addr
		}
	}
}

// OptRegistries sets the addresses of other registry servers option.
func OptRegistries(registries string) IOption {
	return func(o *Option) {
		if registries != "" {
			o.Registries = registries
		}
	}
}

// DefaultOption returns the default options for registering a server.
func DefaultOption() *Option {
	hostname, _ := os.Hostname()
	return &Option{
		Id:            hostname + "-" + xid.New().String(),
		Bind:          ":7370",
		BindAdvertise: ":7370",
	}
}
