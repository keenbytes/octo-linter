# used_actions_*

Groups of rules checking paths of actions used in steps.

## Rules

```yaml
version: '3'
rules:
  used_actions_in_action_steps:
    source: local-or-external
    must_exist: ['local', 'external']
    must_have_valid_inputs: true
  
  used_actions_in_workflow_job_steps:
    source: local-or-external
    must_exist: ['local', 'external']
    must_have_valid_inputs: true
```

|Rule|Description|Value|
|----|-----------|-----|
|source|Referenced action (in `uses`) in steps must have valid path. This rule can be configured to allow local actions, external actions, or both.|One of [Allowed Scopes](#allowed-sources)|
|must_exist|Verifies that the action referenced in a step actually exists. It can be configured to allow only local actions (within the same repository), external actions, or both.|`[]string` that contains `local` and/or `external`|
|must_have_valid_inputs|Verifies that all required inputs are provided when referencing an action in a step, and that no undefined inputs are used.|`bool`|

### Allowed Sources

Below is the list of possible values for the allowed action source:

* `local-or-external`
* `local`
* `external`
