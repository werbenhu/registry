// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package registry

// Api is an api interface for service discovery
type Api interface {

	// Start the discovery server
	// addr: the addr that discovery server listens to
	Start(addr string) error

	// Stop the discovery server
	Stop()
}
