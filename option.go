package registry

import (
	"os"

	"github.com/rs/xid"
)

// Options for register server
type Option struct {

	// Id: service id
	Id string

	// The address used to register the service.
	// If there is a firewall, please remember that the port needs to open both tcp and udp.
	Bind string

	// The address that the service will advertise to other services for registering.
	// Can be used for basic NAT traversal where both the internal ip:port and external ip:port are known.
	BindAdvertise string

	// The addresses of the other registry servers, if there are more than one, separate them with commas,
	// such as "192.168.1.101:7370,192.168.1.102:7370"
	Registries string

	// The address used for service discovery (default ":8080").
	Addr string

	// The address will advertise to client for service discover
	Advertise string
}

type IOption func(o *Option)

func OptId(id string) IOption {
	return func(o *Option) {
		if id != "" {
			o.Id = id
		}
	}
}

func OptAddr(addr string) IOption {
	return func(o *Option) {
		o.Addr = addr
	}
}

func OptAdvertise(addr string) IOption {
	return func(o *Option) {
		o.Advertise = addr
	}
}

func OptBindAdvertise(addr string) IOption {
	return func(o *Option) {
		if addr != "" {
			o.BindAdvertise = addr
		}
	}
}

func OptBind(addr string) IOption {
	return func(o *Option) {
		if addr != "" {
			o.Bind = addr
		}
	}
}

func OptRegistries(registries string) IOption {
	return func(o *Option) {
		if registries != "" {
			o.Registries = registries
		}
	}
}

func DefaultOption() *Option {
	hostname, _ := os.Hostname()
	return &Option{
		Id:            hostname + "-" + xid.New().String(),
		Bind:          ":7370",
		BindAdvertise: ":7370",
	}
}
