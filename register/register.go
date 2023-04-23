// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package register

import registry "github.com/werbenhu/registry"

// Register represents a service registration instance.
type Register struct {
	serf    registry.Discovery // The underlying discovery implementation.
	handler registry.Handler   // The handler function that will be executed when new services are discovered.
	member  *registry.Member   // The service registration metadata.
}

// New creates a new Register instance.
//
// id: The service ID.
// bind: The address used to register the service to registry server. If there is a firewall, please remember that the port needs to open both tcp and udp.
// advertise: The address that the service will advertise to registry server. Can be used for basic NAT traversal where both the internal ip:port and external ip:port are known.
// registries: The addresses of the registry servers, separated by commas, such as "192.168.1.101:7370,192.168.1.102:7370"
// group: Group name of the current service belongs to.
// addr: The address currently provided by this service to the client. For example, if the current service is an HTTP server, this would be the address that HTTP listens to, such as 172.16.3.3:80.
func New(id string, bind string, advertise string, registries string, group string, addr string) *Register {
	member := registry.NewMember(id, bind, advertise, registries, group, addr)
	return &Register{member: member}
}

// SetHandler sets the handler function that will be executed when new services are discovered.
func (r *Register) SetHandler(h registry.Handler) {
	r.handler = h
}

// Start starts the service registration process.
func (r *Register) Start() error {
	r.serf = registry.NewSerf(r.member)
	r.serf.SetHandler(r.handler)
	return r.serf.Start()
}

// Stop stops the service registration process.
func (r *Register) Stop() {
	r.serf.Stop()
}
