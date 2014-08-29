package opensshutil

import (
	"io/ioutil"
	"strings"
)

// Parse OpenSSH configuration file.
func ParseConfig(src []byte) (c *Config, err error) {
	// Create the empty configuration.
	c = &Config{
		Hosts: make([]HostConfig, 0),
	}

	// Create a scanner.
	s := newConfigScanner(src)

	// Iterate across the tokens.
	var host *HostConfig
	var match *HostMatchConfig
	var keyword string
	args := make([]string, 0, 1)

	for {
		tok, lit := s.Scan()
		if err = s.Err(); err != nil {
			return
		}

		// Handle illegal characters.
		if tok == CT_ILLEGAL {
			err = &ParseError{"unexpected or illegal character in configuration file"}
			return
		}

		// Handle end-of-file state.
		if tok == CT_EOF {
			break
		}

		// Handle end-of-line.
		if tok == CT_EOL {
			// Handle complete lines.
			if keyword != "" {
			}
		}

		// Handle normal states.
		switch {
		case tok == CT_KEYWORD:
			if keyword == "" {
				keyword = lit
			} else {
				args = append(args, lit)
			}

		case tok == CT_EQUAL:
			break

		case tok == CT_STRING:
			if keyword == "" {
				err = &ParseError{"expected configuration keyword at beginning of line"}
				return
			}

			args = append(args, lit)

		case tok == CT_EOL && keyword != "":
			kwLower := strings.ToLower(keyword)

			if host == nil && kwLower != "host" {
				err = &ParseError{"expected Host directive before other configuration diretives"}
				return
			}

			switch kwLower {
			case "host":
				host = &HostConfig{
					HostPatterns: args,
					Options: make(ConfigOptionMap),
					Matches: make([]HostMatchConfig, 0),
				}
				match = nil
				c.Hosts = append(c.Hosts, *host)

			case "match":
				if len(args) < 1 {
					err = &ParseError{"a Match directive must at specify contain a condition type"}
				}

				match = &HostMatchConfig{
					ConditionType: args[0],
					ConditionArgs: args[1:],
					Options: make(ConfigOptionMap),
				}
				host.Matches = append(host.Matches, *match)

			default:
				if match != nil {
					match.Options.Set(keyword, args)
				} else {
					host.Options.Set(keyword, args)
				}
			}

			keyword = ""
			args = make([]string, 0, 1)
		}
	}

	return
}

// Read OpenSSH configuration file.
func ReadConfig(path string) (c *Config, err error) {
	var src []byte
	if src, err = ioutil.ReadFile(path); err != nil {
		return
	}

	return ParseConfig(src)
}
