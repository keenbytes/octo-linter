// Package usedactions contains rules checking paths of actions used in steps.
package usedactions

const (
	// ValueLocalOnly defines a configuration value for the referenced action (in 'uses' field) to be local only.
	ValueLocalOnly = "local-only"
	// ValueLocalOnly defines a configuration value for the referenced action (in 'uses' field) to be external only.
	ValueExternalOnly = "external-only"
	// ValueLocalOnly defines a configuration value for the referenced action (in 'uses' field) to be local or external.
	ValueLocalOrExternal = "local-or-external"
)
