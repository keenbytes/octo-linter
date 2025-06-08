# Running locally

## Syntax
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

## Using binary
Tweak below command with a path pointing to `.github` and configuration file:

````
./octo-linter lint -p /path/to/.github -l WARN -c config.yaml
````

## Using docker image
````
docker run --rm --name octo-linter \
  -v /path/to/.github:/dot-github -v $(pwd):/config \
  keenbytes/octo-linter:v1.2.3 \
  lint -p /dot-github -l WARN -c /config/config.yml
````

## Checking secrets
First, create a file with list of secrets that are defined within the repository.

````
% cat ~/secrets-list.txt 
MY_SECRET_1
MY_SECRET_2
````

Use that list with `-s` argument so that octo-linter scans will the secrets in code and compares them with a list.  If there is any secret
that is not on the list, tool will output info about it.

````
    % ./octo-linter validate -p /path/to/.github -s ~/secrets-list.txt -l WARN 2>&1
````

TODO
