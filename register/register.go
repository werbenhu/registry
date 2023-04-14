// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package register

import registry "github.com/werbenhu/registry"

// Register can be easily userd to register a service
type Register struct {
	serf    registry.Discovery
	handler registry.Handler
	member  *registry.Member
}

// New() create a Register object
// id: service id
// bind:
//
//	The address used to register the service to registry server. If there is a firewall, please remember that the port needs to open both tcp and udp.
//
// advertise:
//
//	The address that the service will advertise to registry server. Can be used for basic NAT traversal where both the internal ip:port and external ip:port are known.
//
// registries:
//
//	The addresses of the registry servers, if there are more than one, separate them with commas, such as "192.168.1.101:7370,192.168.1.102:7370"
//
// group:
//
//	Group name of the current service belongs to.
//
// addr:
//
//	The address currently provided by this service to the client, for example, the current service is an http server, that is the address 172.16.3.3:80 that http listens to.
func New(id string, bind string, advertise string, registries string, group string, addr string) *Register {
	member := registry.NewMember(id, bind, advertise, registries, group, addr)
	return &Register{member: member}
}

// SetHandler() set event processing handler when new services are discovered
func (r *Register) SetHandler(h registry.Handler) {
	r.handler = h
}

// Start the register
func (r *Register) Start() error {
	r.serf = registry.NewSerf(r.member)
	r.serf.SetHandler(r.handler)
	return r.serf.Start()
}

// Close the register
func (r *Register) Stop() {
	r.serf.Stop()
}
