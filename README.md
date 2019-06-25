# jiraq

This is a command to dump lifecycle of the backlogs on Jira.

## Getting Started

```sh
go run main.go -url https://jira.example.com \
  -username USER -password PASS \
  -jql "project = NAME AND issuetype = Story ORDER BY updated DESC"
```
