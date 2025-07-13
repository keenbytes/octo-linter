// Package runners contains rules checking GitHub Actions' runners.
package runners

import "errors"

var (
	errValueNotBool = errors.New("value should be bool")
)
