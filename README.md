[![Go Reference](https://pkg.go.dev/badge/github.com/dark-warlord14/CVENotifier.svg)](https://pkg.go.dev/github.com/dark-warlord14/CVENotifier)
[![Go Report Card](https://goreportcard.com/badge/github.com/dark-warlord14/CVEnotifier)](https://goreportcard.com/report/github.com/dark-warlord14/CVEnotifier)

# Customized CVE FEED Notifier

- This tool scrapes the CVE feed from [vuldb.com](https://vuldb.com/?), filters it based on keywords, and notifies via Slack about latest CVE only for the technology or the products you have listed as keywords.

## What it does?

- Parses the RSS feed from [vuldb.com](https://vuldb.com/?rss.recent) using [gofeed](https://github.com/mmcdole/gofeed).
- Filters the feed based on the defined keywords.
- Stores filtered CVEs in a database.
- Sends a Slack notification for each new CVE inserted into the database.

## Installation

Make sure go environment is properly configured
```
go install github.com/dark-warlord14/CVENotifier/cmd/CVENotifier@latest
```
## How to use?

1.  Set the `SLACK_WEBHOOK` environment variable with your Slack webhook URL. For example:

    ```bash
    export SLACK_WEBHOOK=https://hooks.slack.com/services/<id>/<id>
    ```

2.  Set up keywords in `config.yaml`:

    ```yaml
    keywords:
      - Floodlight
      - wordpress
    ```

3.  Run the tool on a regular interval (e.g., every few hours) to fetch the latest feeds and receive notifications for new CVEs. It's recommended to set up a cron job for this.

    ```bash
    CVENotifier -config config.yaml
    ```

cronjob example
```
0 * * * * user CVENotifier -config config.yaml 2>&1 | tee -a CVENotifier.log
```

## Slack Notification
![Slack notification](slack.png)

## To-do

- [x] Fetch RSS feed from  https://vuldb.com/?rss.recent
- [x] Filter the feed if any keyword is present in the title
- [x] Store the data in a database if a keyword is found in the title
- [x] Send a Slack message if the insert operation is successful

## Package Structure

The project is now organized into the following packages:

*   `cmd/CVENotifier`: Contains the main application logic.
*   `internal/config`: Contains the configuration loading logic.
*   `internal/rss`: Contains the RSS feed parsing logic.
*   `internal/slack`: Contains the Slack notification logic.
*   `internal/util`: Contains utility functions such as HTML tag removal.
*   `internal/db`: Contains the database operations logic.
*   `internal/errors`: Contains custom error types.

## Error Handling

The project now uses custom error types to provide more descriptive error messages.
