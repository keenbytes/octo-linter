// Package refvars contains rules checking variables referenced in action or workflow steps, eg. ${{ var }}.
package refvars

import "errors"

var (
	errFileInvalidType = errors.New("file is of invalid type")
	errValueNotBool    = errors.New("value should be bool")
)
