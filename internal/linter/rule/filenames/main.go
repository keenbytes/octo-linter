// Package filenames contains rules related to action and workflow filenames.
package filenames

import (
	"errors"
	"fmt"
)

const (
	// ValueDashCase is a configuration value indicating that a field should follow the dash-case convention.
	ValueDashCase = "dash-case"
	// ValueDashCaseUnderscore is a configuration value indicating that a field should follow the dash-case with
	// underscore prefix allowed convention.
	ValueDashCaseUnderscore = "dash-case;underscore-prefix-allowed"
	// ValueCamelCase is a configuration value indicating that a field should follow the camel-case convention.
	ValueCamelCase = "camelCase"
	// ValuePascalCase is a configuration value indicating that a field should follow the pascal-case convention.
	ValuePascalCase = "PascalCase"
	// ValueAllCaps is a configuration value indicating that a field should follow the all-caps-case convention.
	ValueAllCaps = "ALL_CAPS"
)

var (
	errValueNotString      = errors.New("value should be string")
	errValueNotStringArray = errors.New("value should be []string")
	errValueNotValid       = fmt.Errorf(
		"value can be one of: %s, %s, %s, %s",
		ValueDashCase, ValueCamelCase, ValuePascalCase, ValueAllCaps,
	)
	errValueNotValidIncludingDashCaseUnderscore = fmt.Errorf(
		"value can be one of: %s, %s, %s, %s, %s",
		ValueDashCase, ValueDashCaseUnderscore, ValueCamelCase, ValuePascalCase, ValueAllCaps,
	)
	errValueNotYmlOrYaml = errors.New("value can contain only 'yml' and/or 'yaml'")
	errFileInvalidType   = errors.New("file is of invalid type")
)
