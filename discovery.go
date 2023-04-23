// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package registry

// Auto-discover event notification interface
type Handler interface {

	// OnMemberJoin is triggered when a new service is registered.
	OnMemberJoin(*Member) error

	// OnMemberLeave is triggered when a service leaves.
	OnMemberLeave(*Member) error

	// OnMemberUpdate is triggered when a service is updated.
	OnMemberUpdate(*Member) error
}

// Auto-discover interface.
type Discovery interface {

	// SetHandler sets the event processing handler when new services are discovered.
	SetHandler(Handler)

	// Members returns the members of all services.
	Members() []*Member

	// LocalMember returns the current service.
	LocalMember() *Member

	// Start starts the discovery service.
	Start() error

	// Stop stops the discovery service.
	Stop()
}
