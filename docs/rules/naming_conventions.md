# naming_conventions

Group of rules related to naming convention of action and workflow fields which check if they are in specific format.

## Rules

```yaml
version: '3'
rules:
  naming_conventions:
    action_input_name_format: dash-case
    action_output_name_format: dash-case
    action_referenced_variable_format: ALL_CAPS
    action_step_env_format: ALL_CAPS
    workflow_env_format: ALL_CAPS
    workflow_job_env_format: ALL_CAPS
    workflow_job_step_env_format: ALL_CAPS
    workflow_referenced_variable_format: ALL_CAPS
    workflow_dispatch_input_name_format: dash-case
    workflow_call_input_name_format: dash-case
    workflow_job_name_format: dash-case
    workflow_single_job_only_name: main
```

|Rule|Description|Value|
|----|-----------|-----|
|action_input_name_format|Action input name.|One of [Available Formats](#available-formats)|
|action_output_name_format|Action output name.|One of [Available Formats](#available-formats)|
|action_referenced_variable_format|Referenced variables such as `env`, `var`, and `secret`.|One of [Available Formats](#available-formats)|
|action_step_env_format|Step environment variable names.|One of [Available Formats](#available-formats)|
|workflow_env_format|Workflow environment variable names.|One of [Available Formats](#available-formats)|
|workflow_job_env_format|Workflow job environment variable names.|One of [Available Formats](#available-formats)|
|workflow_job_step_env_format|Workflow job step environment variable names.|One of [Available Formats](#available-formats)|
|workflow_referenced_variable_format|Referenced variables in a workflow such as 'env', 'var', and 'secret'.|One of [Available Formats](#available-formats)|
|workflow_dispatch_input_name_format|`workflow_dispatch` block input name.|One of [Available Formats](#available-formats)|
|workflow_call_input_name_format|`workflow_call` block input name.|One of [Available Formats](#available-formats)|
|workflow_job_name_format|Checks job name.|One of [Available Formats](#available-formats)|
|workflow_single_job_only_name|If workflow has only one job, this should be its name.|`string`|

### Available Formats

Below naming convention formats are available:

* `dash-case`
* `camelCase`
* `PascalCase`
* `snake_case`
* `ALL_CAPS`
