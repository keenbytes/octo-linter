// Package refvars contains rules checking variables referenced in action or workflow steps, eg. ${{ var }}.
package refvars

import "errors"

var (
	errValueNotBool = errors.New("value should be bool")
)
