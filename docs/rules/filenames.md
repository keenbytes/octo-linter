# filenames

Group of rules related to action and workflow filenames.

## Rules

```yaml
version: '3'
rules:
  filenames:
    action_filename_extensions_allowed: ['yml']
    action_directory_name_format: dash-case
    workflow_filename_extensions_allowed: ['yml']
    workflow_filename_base_format: dash-case;underscore-prefix-allowed
```

|Rule|Description|Value|
|----|-----------|-----|
|action_filename_extensions_allowed|Action filename extension must be one of the specified, eg. `yml` or `yaml`.|`[]string`|
|action_directory_name_format|Action directory name adheres to the selected naming convention.|One of [Available Formats](#available-formats)|
|workflow_filename_extensions_allowed|Workflow file extension must be one of specified values, eg. `yml` or `yaml`.|`[]string`|
|workflow_filename_base_format|Workflow file basename (without extension) adheres to the selected naming convention.|One of [Available Formats](#available-formats)|

### Available Formats

Below naming convention formats are available:

* `dash-case`
* `dash-case;underscore-prefix-allowed`
* `camelCase`
* `PascalCase`
* `snake_case`
* `ALL_CAPS`

In case of `dash-case;underscore-prefix-allowed` filename is allowed to have an underscore (`_`) character in the beginning. In some places
it is used to distinguish sub-workflows.
