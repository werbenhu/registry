// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package registry

// Auto-discover event notification interface
type Handler interface {

	// Triggered when a new service is registered
	OnMemberJoin(*Member) error

	// Triggered when a new service leaves
	OnMemberLeave(*Member) error

	// Triggered when a service is updated
	OnMemberUpdate(*Member) error
}

// Auto-discover interface
type Discovery interface {

	// Set event processing handler when new services are discovered
	SetHandler(Handler)

	// Get members of all registry services
	Members() []*Member

	// Get current registry service
	LocalMember() *Member

	// Start the discovery service
	Start() error

	// Stop the discovery service
	Stop()
}
