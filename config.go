package opensshutil

import (
	"strings"
)

// OpenSSH configuration option set.
type ConfigOptionMap map[string][]string

// Set option.
func (m ConfigOptionMap) Set(key string, args []string) {
	keyLower := strings.ToLower(key)

	for k, _ := range m {
		if strings.ToLower(k) == keyLower {
			m[k] = args
			return
		}
	}

	m[key] = args
}

// Get option.
//
// Returns nil if the option was not found.
func (m ConfigOptionMap) Get(key string) []string {
	keyLower := strings.ToLower(key)

	for k, args := range m {
		if strings.ToLower(k) == keyLower {
			return args
		}
	}

	return nil
}

// OpenSSH configuration file host match entry.
type HostMatchConfig struct {
	// Condition type.
	ConditionType string

	// Condition arguments.
	ConditionArgs []string

	// Match-local options.
	Options ConfigOptionMap
}

// OpenSSH configuration file host entry.
type HostConfig struct {
	// Host patterns.
	HostPatterns []string

	// Host-global options.
	Options ConfigOptionMap

	// Matches.
	Matches []HostMatchConfig
}

// OpenSSH configuration file.
//
// Naive representation of a OpenSSH configuration file format, in which keywords
// are not validated.
type Config struct {
	// Hosts.
	Hosts []HostConfig
}
