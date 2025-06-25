# octo-linter

[![Go Reference](https://pkg.go.dev/badge/github.com/keenbytes/octo-linter.svg)](https://pkg.go.dev/github.com/keenbytes/octo-linter) [![Go Report Card](https://goreportcard.com/badge/github.com/keenbytes/octo-linter)](https://goreportcard.com/report/github.com/keenbytes/octo-linter)

![octo-linter](octo-linter.png "octo-linter")

A tool that validates GitHub Actions workflow and action YAML files. It checks for syntax errors, such as
invalid inputs and outputs, and lints for missing descriptions, invalid rules, and other best practice
violations, ensuring your workflows are error-free and adhere to GitHub Actions standards.

This application is a refactored and enhanced software that I
have created few years ago.  In the new version, rules can be configured, checks are executed in parallel,
and log level command line argument has been introduced.  Also, each rule's source code is now extracted
to a separate file for better maintenance.

## Full documentation
[Link to docs](https://keenbytes.github.io/octo-linter/)

## Running
Check below help message for `lint` command:

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
[`internal/linter/dotgithub.yml`](internal/linter/dotgithub.yml).

**Use `init` command to create a default `dotgithub.yml` configuration file in current directory.**

### Using docker image
Note that the image has to be present, either built or pulled from the registry.
Replace path to the .github directory.

    docker run --rm --name tmp-octo-linter \
      -v /Users/me/my-repo/.github:/dotgithub \
      keenbytes/octo-linter:v2.0.0 \
	  validate -p /dotgithub


## Exit code
Tool exits with exit code `0` when everything is fine.  `1` when there are errors, `2` when there are only
warnings.  Additionally it may exit with a different code, eg. `22`.  These numbers indicate another error
whilst reading files.

