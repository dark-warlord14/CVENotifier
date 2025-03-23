package errors

import "fmt"

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("Config error: %s", e.Message)
}

type RSSFeedError struct {
	Message string
}

func (e *RSSFeedError) Error() string {
	return fmt.Sprintf("RSS feed error: %s", e.Message)
}

type DatabaseError struct {
	Message string
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("Database error: %s", e.Message)
}

type SlackNotificationError struct {
	Message string
}

func (e *SlackNotificationError) Error() string {
	return fmt.Sprintf("Slack notification error: %s", e.Message)
}
