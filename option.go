package srouter

import (
	"os"

	"github.com/rs/xid"
)

type Option struct {
	Id        string
	Addr      string
	Advertise string
	Members   string
	ApiPort   string
	Service   string
}

type IOption func(o *Option)

func OptId(id string) IOption {
	return func(o *Option) {
		if id != "" {
			o.Id = id
		}
	}
}

func OptApiPort(port string) IOption {
	return func(o *Option) {
		o.ApiPort = port
	}
}

func OptAdvertise(addr string) IOption {
	return func(o *Option) {
		if addr != "" {
			o.Advertise = addr
		}
	}
}

func OptAddr(addr string) IOption {
	return func(o *Option) {
		if addr != "" {
			o.Addr = addr
		}
	}
}

func OptMembers(members string) IOption {
	return func(o *Option) {
		if members != "" {
			o.Members = members
		}
	}
}

func OptService(service string) IOption {
	return func(o *Option) {
		if service != "" {
			o.Service = service
		}
	}
}

func DefaultOption() *Option {
	hostname, _ := os.Hostname()
	return &Option{
		Id:        hostname + "-" + xid.New().String(),
		Addr:      ":7933",
		Advertise: ":7933",
	}
}
