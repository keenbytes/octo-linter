# workflow_runners

Group of rules checking various things related to runner.

## Rules

```yaml
version: '3'
rules:
  workflow_runner:
    not_latest: true
```

|Rule|Description|
|----|-----------|
|not_latest|Checks whether 'runs-on' does not contain the 'latest' string. In some case, runner version (image) should be frozen, instead of using the latest.|
