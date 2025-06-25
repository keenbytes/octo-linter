# referenced_variables_*

Groups of rules checking variables referenced in action or workflow steps, eg. `${{ var }}`.

## Rules

```yaml
version: '3'
rules:
  referenced_variables_in_actions:
    not_one_word: true
    not_in_double_quotes: true

  referenced_variables_in_workflows:
    not_one_word: true
    not_in_double_quotes: true
```

|Rule|Description|Value|
|----|-----------|-----|
|not_one_word|Checks for variable references that are single-word or single-level, e.g. `${{ something }}` instead of `${{ inputs.something }}`. Only the values `true` and `false` are permitted in this form; all other variables are considered invalid.|`bool`|
|not_in_double_quotes|Scans for all variable references enclosed in double quotes. It is safer to use single quotes, as double quotes expand certain characters and may allow the execution of sub-commands.|`bool`|
