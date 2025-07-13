// Package required contains rules checking if required fields within actions and workflows are defined.
package required

import "errors"

var (
	errValueNotBool              = errors.New("value should be bool")
	errValueNotStringArray       = errors.New("value should be []string")
	errValueNotNameOrDescription = errors.New("value can contain only 'name' and/or 'description'")
	errValueNotName              = errors.New("value can contain only 'name'")
	errValueNotDescription       = errors.New("value can contain only 'description'")
	errFileInvalidType           = errors.New("file is of invalid type")
)
