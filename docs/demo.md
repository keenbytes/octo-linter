# Demo

## Clone and run
An example configuration file and `.github` directory structure can be found in the `example` directory of the 
repository. Use below commands to see octo-linter in action:

````
git clone https://github.com/keenbytes/octo-linter.git
cd octo-linter/example

docker run --platform=linux/amd64 --rm --name octo-linter \
  -v $(pwd)/dot-github:/dot-github \
  -v $(pwd):/config \
  keenbytes/octo-linter:v2.0.0 \
  lint -p /dot-github -l WARN -c /config/config.yml
````

## Output
This should generate an output similar to the following:

````
time=2025-06-08T21:18:42.537Z level=WARN msg="action_file_extensions: action 'InvalidActionExtension' file extension must be one of: yml"
time=2025-06-08T21:18:42.548Z level=WARN msg="action_directory_name: action directory name 'InvalidActionExtension' must be lower-case and hyphens only"
time=2025-06-08T21:18:42.548Z level=WARN msg="action_directory_name: action directory name 'InvalidActionName' must be lower-case and hyphens only"
time=2025-06-08T21:18:42.548Z level=ERROR msg="action_step_action: action 'InvalidActionExtension' step 1 calls action 'actions/checkout@v4' that is not a valid local path"
time=2025-06-08T21:18:42.548Z level=ERROR msg="action_step_action: action 'InvalidActionName' step 1 calls action 'actions/checkout@v4' that is not a valid local path"
time=2025-06-08T21:18:42.548Z level=ERROR msg="action_called_step_output_exists: action 'some-action' calls a step 'non-existing-step' output 'output1' but step does not exist"
time=2025-06-08T21:18:42.549Z level=ERROR msg="action_step_action: action 'some-action' step 1 calls action 'actions/checkout@v4' that is not a valid local path"
time=2025-06-08T21:18:42.550Z level=ERROR msg="action_input_required: action 'some-action' input 'output-without-description' does not have a required description"
time=2025-06-08T21:18:42.552Z level=ERROR msg="action_step_env: action 'some-action' step 2 env 'InvalidEnvName' must be alphanumeric uppercase and underscore only"
time=2025-06-08T21:18:42.552Z level=ERROR msg="action_called_input_exists: action 'some-action' calls an input 'non-existing' that does not exist"
time=2025-06-08T21:18:42.552Z level=ERROR msg="workflow_called_input_exists: workflow 'workflow1.yaml' calls an input 'non-existing' that does not exist"
time=2025-06-08T21:18:42.553Z level=ERROR msg="workflow_runs_on_not_latest: workflow 'workflow1.yaml' job 'job-2' should not use 'latest' in 'runs-on' field"
time=2025-06-08T21:18:42.554Z level=WARN msg="workflow_file_extensions: workflow 'workflow1' file extension must be one of: yml"
````

## Exit code
Tool exits with exit code `0` when everything is fine.  `1` when there are errors, `2` when there are only
warnings.  Additionally it may exit with a different code, eg. `22`.  These numbers indicate another error
whilst reading files.

## Checking secrets and vars
octo-linter can scan the code for `secrets` and `variables` and compare them with file containing list of defined one.  If there is any `secret`
or `var` that is not on the list, tool will output info about it.  See below run and its output.

````
docker run --platform=linux/amd64 --rm --name octo-linter \
  -v $(pwd)/dot-github:/dot-github \
  -v $(pwd):/config \
  keenbytes/octo-linter:v2.0.0 \
  lint -p /dot-github -l WARN -c /config/config.yml \
  -s /config/secrets_list.txt \
  -z /config/vars_list.txt \
  2>&1 | grep NON_EXISTING_ONE
time=2025-06-08T22:09:18.788Z level=ERROR msg="workflow_called_variable_exists_in_file: workflow 'workflow1.yaml' calls a variable 'NON_EXISTING_ONE' that does not exist in the vars file"
time=2025-06-08T22:09:18.789Z level=ERROR msg="workflow_called_variable_exists_in_file: workflow 'workflow1.yaml' calls a secret 'NON_EXISTING_ONE' that does not exist in the secrets file"
````
