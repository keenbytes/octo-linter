# Configuration

octo-linter requires a configuration file that defines compliance rules. This section explains how to create one in more detail.

## Initialise configuration file
Use `init` command that will create a sample configuration file named `dotgithub.yml` in
current directory. Use `-d` to write it in another place.

## Requirements
Let’s consider a GitHub repository that contains workflows and actions within the `.github` directory. Several 
developers are contributing to it, and we want to enforce the following rules for the files in that directory:

* Action names must use only lowercase alphanumeric characters and hyphens
* Action and workflow files should have a .yml extension
* Named-value variables should not be enclosed in double quotes
* The use of the latest runner version should be avoided
* Actions, along with their inputs and outputs (where applicable), must include both `name` and `description` fields
* Only local actions should be used
* Environment variables in steps must use uppercase alphanumeric characters, optionally including underscores

Additionally, it would be useful to automatically verify that all referenced inputs, outputs, and similar entities are properly defined.

There are many more possible rules, but we’ll focus on these for the purpose of this example.

## Configuration file
Tweak the configuration file with rules that the application would use.

Based on the list in previous section, the configuration file can look as shown below.

````yaml
version: '2'
rules:
  # Action names must use only lowercase alphanumeric characters and hyphens
  action_directory_name: lowercase-hyphens

  # Action and workflow files should have a .yml extension
  action_file_extensions: ['yml']
  workflow_file_extensions: ['yml']
  action_called_variable_not_in_double_quote: true

  # Named-value variables should not be enclosed in double quotes
  workflow_called_variable_not_in_double_quote: true

  # The use of the latest runner version should be avoided
  workflow_runs_on_not_latest: true

  # Actions, along with their inputs and outputs (where applicable), must include both name and description fields
  action_required__name: true 
  action_required__description: true
  action_input_required__description: true
  action_output_required__description: true

  # Only local actions should be used
  action_step_action: local-only

  # Environment variables in steps must use uppercase alphanumeric characters, optionally including underscores
  action_step_env: uppercase-underscores 
  
  # All referenced inputs, outputs, and similar entities are properly defined
  action_called_input_exists: true 
  action_called_step_output_exists: true
  workflow_called_variable_exists_in_file: true
  workflow_called_input_exists: true
````

### Error or warning
A non-compliant rule can be treated either as an error or a warning. If a rule is intended to trigger only a warning, it should be included in the warning_only list, as shown below:

````yaml
warning_only:
  - action_directory_name
  - action_file_extensions
  - workflow_file_extensions
````

Continue to the next section to learn how to run `octo-linter` using the prepared configuration.
