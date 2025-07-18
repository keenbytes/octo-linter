# Overview

![octo-linter](assets/logo.png "octo-linter")

**octo-linter** is a tool that validates GitHub Actions **workflow and action YAML files**.  It checks for **syntax errors**, such as
**calling invalid inputs and outputs**, and lints for **missing descriptions, invalid naming conventions, and other best practice
violations**, ensuring your workflows are error-free and adhere to GitHub Actions and your company standards.

## Motivation

The tool was developed during a large-scale refactor of existing GitHub Actions code, which was scattered across multiple repositories with no consistent standards in place.  To streamline the process and reduce manual effort, it made sense to automate many of the checks that would otherwise fall to the reviewer.  Notably, GitHub does not raise errors in several cases — for example, when referencing a non-existent input, it simply substitutes it with an empty string.  This behaviour can be difficult to detect, particularly when code is being moved or restructured during refactoring. 

## Demo

Please navigate to [Demo](demo.md) to see example usage.
