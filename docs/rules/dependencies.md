# dependencies

Group of rules checking various dependencies between action steps, workflow jobs etc.

## Rules

```yaml
version: '3'
rules:
  dependencies:
    workflow_needs_field_must_contain_already_existing_jobs: true
    action_referenced_input_must_exists: true
    action_referenced_step_output_must_exist: true
    workflow_referenced_input_must_exists: true
    workflow_referenced_variable_must_exists_in_attached_file: true
```

|Rule|Description|Value|
|----|-----------|-----|
|workflow_needs_field_must_contain_already_existing_jobs|Checks if `needs` field references existing jobs.|`bool`|
|action_referenced_input_must_exists|Scans the action code for all input references and verifies that each has been previously defined. During action execution, if a reference to an undefined input is found, it is replaced with an empty string.|`bool`|
|action_referenced_step_output_must_exist|Checks whether references to step outputs correspond to outputs defined in preceding steps. During execution, referencing a non-existent step output results in an empty string. |`bool`|
|workflow_referenced_input_must_exists|Scans the code for all input references and verifies that each has been previously defined. During execution, if a reference to an undefined input is found, it is replaced with an empty string.|`bool`|
|workflow_referenced_variable_must_exists_in_attached_file|Checks if called variables and secrets exist. This rule requires a list of variables and secrets to be checked against.|`bool`|
