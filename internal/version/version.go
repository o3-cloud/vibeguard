// Package version provides version information for vibeguard.
package version

// Version is the semantic version of vibeguard.
// This can be overridden at build time using:
//
//	go build -ldflags="-X github.com/vibeguard/vibeguard/internal/version.Version=v1.0.0"
var Version = "v0.1.0-dev"

// String returns the version string.
func String() string {
	return Version
}
