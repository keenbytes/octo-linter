// Package casematch provides helper functions for enforcing naming conventions.
package casematch

import "regexp"

// Match checks whether a string follows a specific case convention.
func Match(str string, cas string) bool {
	var re *regexp.Regexp

	switch cas {
	case "dash-case":
		re = regexp.MustCompile(`^[a-z0-9][a-z0-9\-]+$`)
	case "dash-case;underscore-prefix-allowed":
		re = regexp.MustCompile(`^[_]{0,1}[a-z0-9][a-z0-9\-]+$`)
	case "camelCase":
		re = regexp.MustCompile(`^[a-z][A-Za-z0-9]+$`)
	case "PascalCase":
		re = regexp.MustCompile(`^[A-Z][A-Za-z0-9]+$`)
	case "ALL_CAPS":
		re = regexp.MustCompile(`^[A-Z][A-Z0-9_]+$`)
	default:
		return true
	}

	return re.MatchString(str)
}
