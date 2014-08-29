package opensshutil

// OpenSSH configuration option.
type ConfigOption []string

// OpenSSH configuration file host match entry.
type HostMatchConfig struct {
	// Condition type.
	ConditionType string

	// Condition arguments.
	ConditionArgs []string

	// Match-local options.
	Options map[string]ConfigOption
}

// OpenSSH configuration file host entry.
type HostConfig struct {
	// Host patterns.
	HostPatterns []string

	// Host-global options.
	Options map[string]ConfigOption

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
