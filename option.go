package srouter

import (
	"os"

	"github.com/rs/xid"
)

type Option struct {
	Id        string
	Addr      string
	Advertise string
	Routers   string
	ApiAddr   string
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

func OptApiAddr(addr string) IOption {
	return func(o *Option) {
		o.ApiAddr = addr
	}
}

func OptService(addr string) IOption {
	return func(o *Option) {
		o.Service = addr
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

func OptRouters(routers string) IOption {
	return func(o *Option) {
		if routers != "" {
			o.Routers = routers
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
