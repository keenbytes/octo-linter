// Package dependencies contains rules checking various dependencies between action steps, workflow jobs etc.
package dependencies

import "errors"

var (
	errValueNotBool = errors.New("value should be bool")
)
