// Package filenames contains rules related to action and workflow filenames.
package filenames

const (
	// ValueDashCase is a configuration value indicating that a field should follow the dash-case convention.
	ValueDashCase = "dash-case"
	// ValueDashCase is a configuration value indicating that a field should follow the dash-case with underscore prefix allowed convention.
	ValueDashCaseUnderscore = "dash-case;underscore-prefix-allowed"
	// ValueDashCase is a configuration value indicating that a field should follow the camel-case convention.
	ValueCamelCase = "camelCase"
	// ValueDashCase is is a configuration value indicating that a field should follow the pascal-case convention.
	ValuePascalCase = "PascalCase"
	// ValueDashCase is a configuration value indicating that a field should follow the all-caps-case convention.
	ValueAllCaps = "ALL_CAPS"
)
