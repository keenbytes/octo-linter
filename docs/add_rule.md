# New rule

## Requirements

To add a new rule the following things have to be created:

* rule struct with its functionality which implements the methods as in the `rule.Rule` interface, such as `Validate` and `Lint`
* rule struct is either created in a new group or existing one, see directories in `internal/linter/rule`
* configuration for the rule must be added to the default configuration file in `internal/linter/dotgithub.yml`
* rule must be linked with a specific configuration key, and that is done in the `gen.go` file
* tests should be written
* documentation must be updated

See below sections to get more information on the above topics.

:warning: **Multiple rules in the configuration file can be handled by one rule struct.  Please read the whole docs.**

### Rule group

Every rule is placed in a specific group. In configuration these can be find as the second-level keys in the `rules` section,
and for example such guys are `naming_conventions` or `required_fields`.  These groups correspond to specific directories in 
the `internal/rules` directory.  For previous examples, these are `naming` and `required`.

New rule can be put in a new group or an existing one.

### Rule struct

The easiest is to use an existing rule as a template, copy it and modify it.  Every rule needs to implement methods from the
interface which is in the `internal/linter/rule` (shown below).

```go
// Rule represents a rule.
type Rule interface {
	Validate(conf interface{}) error
	Lint(
		config interface{},
		f dotgithub.File,
		d *dotgithub.DotGithub,
		chErrors chan<- glitch.Glitch,
	) (bool, error)
	ConfigName(fileType int) string
	FileType() int
}
```

* `Validate` is a method that checks the configuration value.
* `Lint` is the main method that runs the lint on a specific file (action or workflow) with specific configuration value etc.
* `ConfigName` returns the name of the key in the configuration. This is used in many places to link the rule with configuration and it's shown in the error (warning) messages. And also, it takes as argument an integer that defines what is the type of the file that is checked (action or workflow).
* `FileType` returns integer which is a bitmask indicating what file types are linted by this rule.

#### FileType method
If rule is just for action file it would look as shown below:

```go

// FileType returns an integer that specifies the file types (action and/or workflow) the rule targets.
func (r ActionReferencedStepOutputExists) FileType() int {
	return rule.DotGithubFileTypeAction
}
```

However, if it lints both types of files then it would return `rule.DotGithubFileTypeAction | rule.DotGithubFileTypeWorkflow`.


#### ConfigName method
This method might get tricky when a rule struct is used to validate multiple rules (keys) in the configuration file.

When rule struct is used just for a single rule then the method is simple, as shown below.

```go
// ConfigName returns the name of the rule as defined in the configuration file.
func (r ActionReferencedStepOutputExists) ConfigName(int) string {
	return "dependencies__action_referenced_step_output_must_exist"
}
```

However, when the rule struct is used to lint both action and workflow files it would return configuration key name dependent on the type of the file that is
checked, like shown on the example below.

```go

// ConfigName returns the name of the rule as defined in the configuration file.
func (r ReferencedInputExists) ConfigName(t int) string {
	switch t {
	case rule.DotGithubFileTypeWorkflow:
		return "dependencies__workflow_referenced_input_must_exists"
	case rule.DotGithubFileTypeAction:
		return "dependencies__action_referenced_input_must_exists"
	default:
		return "dependencies__*_referenced_input_must_exists"
	}
}
```

There is yet another scenario, where a rule has an additional custom field and the name is dependent on it. See below.

```go
// ConfigName returns the name of the rule as defined in the configuration file.
func (r Action) ConfigName(int) string {
	switch r.Field {
	case ActionFieldAction:
		return "required_fields__action_requires"
	case ActionFieldInput:
		return "required_fields__action_input_requires"
	case ActionFieldOutput:
		return "required_fields__action_output_requires"
	default:
		return "required_fields__action_*_requires"
	}
}
```

#### Lint method
To distinguish lint error from other error, it is called a `Glitch`, and instance of `glitch.Glitch` must be send to the `chErrors` channel.

To write that method, just use any of the existing ones. Check the default configuration file in `internal/linter/dotgithub.yml` and find one that is similar (again, check the ConfigName section for the three different scenarios).

#### Validate method
Depending on the value, use any of the existing code as a template.

### Configuration file
New rule must be added to the default configuration file found in the `internal/linter/dotgithub.yml`.

### Link configuration key with rule struct
When octo-linter parses configuration file, it needs to create instantiate rule structs from it. Hence, every configuration key must correspond to a specific
rule. This is done by a loop that is generated with `gen.go` file. 

Going back to the three scenarios (though it can be more) described in the ConfigName method section, here are corresponding snippets of code for that.

```go
			"dependencies__action_referenced_step_output_must_exist": {
				N: "dependencies.ActionReferencedStepOutputExists",
			},
```

```go

			"dependencies__action_referenced_input_must_exists": {
				N: "dependencies.ReferencedInputExists",
			},
			// ...
			"dependencies__workflow_referenced_input_must_exists": {
				N: "dependencies.ReferencedInputExists",
			},
```

```go
			"required_fields__action_requires": {
				N: "required.Action",
				F: map[string]string{"Field": `required.ActionFieldAction`},
			},
			"required_fields__action_input_requires": {
				N: "required.Action",
				F: map[string]string{"Field": `required.ActionFieldInput`},
			},
			"required_fields__action_output_requires": {
				N: "required.Action",
				F: map[string]string{"Field": `required.ActionFieldOutput`},
			},
```

### Documentation
Once new rule is working, and covered with tests, it must be properly documented.
