// Package usedactions contains rules checking paths of actions used in steps.
package usedactions

import (
	"errors"
	"fmt"
)

const (
	// ValueLocalOnly defines a configuration value for the referenced action (in 'uses' field) to be local only.
	ValueLocalOnly = "local-only"
	// ValueExternalOnly defines a configuration value for the referenced action (in 'uses' field) to be external only.
	ValueExternalOnly = "external-only"
	// ValueLocalOrExternal defines a configuration value for the referenced action (in 'uses' field) to be local or
	// external.
	ValueLocalOrExternal = "local-or-external"
)

var (
	errValueNotBool = errors.New("value should be bool")
	errValueNotString = errors.New("value should be string")
	errValueNotStringArray = errors.New("value should be []string")
	errValueNotLocalAndOrExternal = errors.New("value can contain only 'local' and/or 'external'")
	errValueNotEmptyOrLocalOrExternalOrBoth = fmt.Errorf(
		"value can be '%s', '%s', '%s' or empty string",
		ValueLocalOnly,
		ValueLocalOrExternal,
		ValueExternalOnly,
	)
)
