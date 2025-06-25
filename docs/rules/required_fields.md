# required_fields

Group of rules checking if required fields within actions and workflows are defined.

## Rules

```yaml
version: '3'
rules:
  required_fields:
    action_requires: ['name', 'description']
    action_input_requires: ['description']
    action_output_requires: ['description']
    workflow_requires: ['name']
    workflow_dispatch_input_requires: ['description']
    workflow_call_input_requires: ['description']
    workflow_requires_uses_or_runs_on: true
```

|Rule|Description|Value|
|----|-----------|-----|
|action_requires|Fields in the root of the action.|`[]string`|
|action_input_requires|Fields in action inputs.|`[]string`|
|action_output_requires|Fields in action outputs.|`[]string`|
|workflow_requires|Fields in the root of workflow.|`[]string`|
|workflow_dispatch_input_requires|`workflow_dispatch` inputs fields.|`[]string`|
|workflow_call_input_requires|`workflow_call` input fields.|`[]string`|
|workflow_requires_uses_or_runs_on|Checks if workflow has `runs-on` or `uses` field. At least of them must be defined.|`bool`|
