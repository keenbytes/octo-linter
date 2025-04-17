# octo-linter

[![Go Reference](https://pkg.go.dev/badge/gopkg.pl/mikogs/octo-linter.svg)](https://pkg.go.dev/gopkg.pl/mikogs/octo-linter) [![Go Report Card](https://goreportcard.com/badge/gopkg.pl/mikogs/octo-linter)](https://goreportcard.com/report/gopkg.pl/mikogs/octo-linter)

![octo-linter](octo-linter.png "octo-linter")

A tool that validates GitHub Actions workflow and action YAML files. It checks for syntax errors, such as
invalid inputs and outputs, and lints for missing descriptions, invalid rules, and other best practice
violations, ensuring your workflows are error-free and adhere to GitHub Actions standards.

This application is a refactored and enhanced
[github-actions-validator](https://github.com/mikogs/github-actions-validator), another software that I
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

