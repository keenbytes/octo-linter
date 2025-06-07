# octo-linter

[![Go Reference](https://pkg.go.dev/badge/github.com/keenbytes/octo-linter.svg)](https://pkg.go.dev/github.com/keenbytes/octo-linter) [![Go Report Card](https://goreportcard.com/badge/github.com/keenbytes/octo-linter)](https://goreportcard.com/report/github.com/keenbytes/octo-linter)

![octo-linter](octo-linter2.png "octo-linter")

A tool that validates GitHub Actions workflow and action YAML files. It checks for syntax errors, such as
invalid inputs and outputs, and lints for missing descriptions, invalid rules, and other best practice
violations, ensuring your workflows are error-free and adhere to GitHub Actions standards.

This application is a refactored and enhanced
[github-actions-validator](https://github.com/keenbytes/github-actions-validator), another software that I
have created few years ago.  In the new version, rules can be configured, checks are executed in parallel,
and log level command line argument has been introduced.  Also, each rule's source code is now extracted
to a separate file for better maintenance.

## Building
Run `go build -o octo-linter` to compile the binary.

### Building docker image
To build the docker image, use the following command.

    docker build -t octo-linter .


## Running
Check below help message for `validate` command:

    Usage:  octo-linter lint [FLAGS]
    
    Runs the linter on files from a specific directory
    
    Required flags: 
      -p,	 --path DIR       Path to .github directory
    
    Optional flags: 
      -c,	 --config FILE    Linter config with rules in YAML format
      -l,	 --loglevel       One of NONE,ERR,WARN,DEBUG
      -s,	 --secrets-file   Check if secret names exist in this file (one per line)
      -z,	 --vars-file      Check if variable names exist in this file (one per line)

Use `-p` argument to point to `.github` directories.  The tool will search for any actions in the `actions`
directory, where each action is in its own sub-directory and its filename is either `action.yaml` or
`action.yml`.  And, it will search for workflows' `*.yml` and `*.yaml` files in `workflows` directory.

Additionally, all the variable names (meaning `${{ var.NAME }}`) as well as secrets (`${{ secret.NAME }}`)
in the workflow can be checked against a list of possible names.  Use `-z` and `-s` arguments with paths
to files containing a list of possible variable or secret names, with names being separated by new line or
space.

### Configuration file
Octo-linter can be told what rules should be executed and which of them should be classified as errors.  The
rest will be shown as warnings.

If config is not passed, then the default one is used.  It can be found in 
[`pkg/linter/dotgithub.yml`](pkg/linter/dotgithub.yml).

### Rules
Best way to check available rules is to navigate to `pkg/linter/rule/` directory.  However, there's a list below as well.

<details>
<summary>
List of rules.
</summary>

|Rule|Description|
|----|-----------|
|rule_action_called_input_exists|Scans the action code for all input references and verifies that each has been previously defined. During action execution, if a reference to an undefined input is found, it is replaced with an empty string.|
|rule_action_called_step_output_exists|Checks whether references to step outputs correspond to outputs defined in preceding steps. During execution, referencing a non-existent step output results in an empty string. |
|rule_action_called_variable|Verifies that referenced variables such as 'env', 'var', and 'secret' follow the defined casing rule. Currently, only 'uppercase-underscores' is supported, meaning variables must be fully uppercase and may include underscores.|
|rule_action_called_variable_not_in_double_quote|Scans for all variable references enclosed in double quotes. It is safer to use single quotes, as double quotes expand certain characters and may allow the execution of sub-commands.|
|rule_action_called_variable_not_one_word|Checks for variable references that are single-word or single-level, e.g. '${{ something }}' instead of '${{ inputs.something }}'. Only the values 'true' and 'false' are permitted in this form; all other variables are considered invalid.|
|rule_action_directory_name|Checks whether the action directory name adheres to the selected naming convention. Currently, only 'lowercase-hyphens' is supported, meaning the name must be entirely lowercase and use hyphens only.|
|rule_action_file_extensions|Checks if action file extension is one of the specific values, eg. 'yml' or 'yaml'.|
|rule_action_input_required|Checks whether specific input attributes are defined (e.g. 'description'). Currently, only the 'description' attribute is supported.|
|rule_action_input_value|Verifies whether the action input field follows the specified naming convention — for example, ensuring the 'name' field uses 'lowercase-hyphens' (lowercase letters, digits, and hyphens only).|
|rule_action_output_required|Checks whether specific output attributes are defined (e.g. 'description'). Currently, only the 'description' attribute is supported.|
|rule_action_output_value|Verifies whether the action output field follows the specified naming convention — for example, ensuring the 'name' field uses 'lowercase-hyphens' (lowercase letters, digits, and hyphens only).|
|rule_action_required|Checks whether the specified action fields are present, e.g. 'name'.|
|rule_step_action|Checks whether the referenced actions have valid paths. This rule can be configured to allow local actions, external actions, or both.|
|rule_step_action_exists|Verifies that the action referenced in a step actually exists. It can be configured to allow only local actions (within the same repository), external actions, or both.|
|rule_step_action_input_valid|Verifies that all required inputs are provided when referencing an action in a step, and that no undefined inputs are used.|
|rule_step_env|Checks whether step environment variable names follow the specified naming convention. Currently, only 'uppercase-underscores' is supported, meaning variable names may contain uppercase letters, numbers, and underscores only.|
|rule_workflow_call_input_required|Checks whether specific workflow_call input attributes are defined (e.g. 'description'). Currently, only the 'description' attribute is supported.|
|rule_workflow_call_input_value|Verifies whether the workflow_call input field follows the specified naming convention — for example, ensuring the 'name' field uses 'lowercase-hyphens' (lowercase letters, digits, and hyphens only).|
|rule_workflow_called_input_exists|Scans the code for all input references and verifies that each has been previously defined. During execution, if a reference to an undefined input is found, it is replaced with an empty string.|
|rule_workflow_called_variable|Verifies that referenced variables such as 'env', 'var', and 'secret' follow the defined casing rule. Currently, only 'uppercase-underscores' is supported, meaning variables must be fully uppercase and may include underscores.|
|rule_workflow_called_variable_exists_in_file|Checks if called variables and secrets exist. This rule requires a list of variables and secrets to be checked against.|
|rule_workflow_called_variable_not_in_double_quote.go|Scans for all variable references enclosed in double quotes. It is safer to use single quotes, as double quotes expand certain characters and may allow the execution of sub-commands.|
|rule_workflow_called_variable_not_one_word|Checks for variable references that are single-word or single-level, e.g. '${{ something }}' instead of '${{ inputs.something }}'. Only the values 'true' and 'false' are permitted in this form; all other variables are considered invalid.|
|rule_workflow_dispatch_input_required|Checks whether specific workflow_dispatch input attributes are defined (e.g. 'description'). Currently, only the 'description' attribute is supported.|
|rule_workflow_dispatch_input_value|Verifies whether the workflow_dispatch input field follows the specified naming convention — for example, ensuring the 'name' field uses 'lowercase-hyphens' (lowercase letters, digits, and hyphens only).|
|rule_workflow_env|Checks whether workflow environment variable names follow the specified naming convention. Currently, only 'uppercase-underscores' is supported, meaning variable names may contain uppercase letters, numbers, and underscores only.|
|rule_workflow_file_extensions|Checks if workflow file extension is one of the specific values, eg. 'yml' or 'yaml'.|
|rule_workflow_job_env|Checks whether workflow job environment variable names follow the specified naming convention. Currently, only 'uppercase-underscores' is supported, meaning variable names may contain uppercase letters, numbers, and underscores only.|
|rule_workflow_job_needs_exist|Checks if 'needs' references existing jobs.|
|rule_workflow_job_value|Checks if workflow job fields follow specified naming convention, for example if 'name' is 'lowercase-hyphens'.|
|rule_workflow_required|Checks whether the specified workflow fields are present, e.g. 'name'.|
|rule_workflow_required_uses_or_runs_on|Checks if workflow has 'runs-on' or 'uses' field. At least of them must be defined.|
|rule_workflow_runs_on_not_latest|Checks whether 'runs-on' does not contain the 'latest' string. In some case, runner version (image) should be frozen, instead of using the latest.|
|rule_workflow_single_job_main|Checks if workflow's only job is called 'main' - just for naming consistency.|
</details>

### Example of checking secrets

    % cat ~/secrets-list.txt 
    MY_SECRET_1
    MY_SECRET_2
    % ./octo-linter validate -p /path/to/.github -s ~/secrets-list.txt -l WARN 2>&1| grep 'action_step_env:'
    wrn: action_step_env: action 'store-artifacts' step 5 env 'aws_access_key_id' must be alphanumeric uppercase and underscore only
    wrn: action_step_env: action 'store-artifacts' step 5 env 'aws_secret_access_key' must be alphanumeric uppercase and underscore only

### Using docker image
Note that the image has to be present, either built or pulled from the registry.
Replace path to the .github directory.

    docker run --rm --name tmp-octo-linter \
      -v /Users/me/my-repo/.github:/dotgithub \
      octo-linter \
	  validate -p /dotgithub


## Exit code
Tool exits with exit code `0` when everything is fine.  `1` when there are errors, `2` when there are only
warnings.  Additionally it may exit with a different code, eg. `22`.  These numbers indicate another error
whilst reading files.

