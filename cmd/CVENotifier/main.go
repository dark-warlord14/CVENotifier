// cmd/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dark-warlord14/CVENotifier/internal/config"
	"github.com/dark-warlord14/CVENotifier/internal/db"
	"github.com/dark-warlord14/CVENotifier/internal/errors"
	"github.com/dark-warlord14/CVENotifier/internal/rss"
	"github.com/joho/godotenv"
)

type Config struct {
	Keywords []string `yaml:"keywords"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("main: Error loading .env file: %v", err)
	}

	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Path to the configuration YAML file")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	feed, err := rss.ParseFeed("https://vuldb.com/?rss.recent")
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	databasePath := "CVENotifier.db"
	dbConn, err := db.InitDB(databasePath)
	if err != nil {
		log.Fatalf("main: %v", err)
	}
	defer dbConn.Close()

	slackWebhook := os.Getenv("SLACK_WEBHOOK")
	if slackWebhook == "" {
		log.Fatalf("main: SLACK_WEBHOOK environment variable not set")
	}

	var matchFound = 0

	for _, item := range feed.Items {
		for _, keyword := range cfg.Keywords {
			if strings.Contains(strings.ToLower(item.Title), strings.ToLower(keyword)) {
				matchFound++

				log.Printf("Matched Keyword: " + keyword)
				log.Printf("Title: " + item.Title)
				log.Printf("Link: " + item.Link)
				log.Printf("Published Date: " + item.Published)
				log.Printf("Categories: " + strings.Join(item.Categories, ","))
				log.Printf("Description: " + item.Description)

				description := item.Description
				if description == "" {
					description = "No description available."
				}

				err = db.InsertData(dbConn, item.Title, item.Link, item.Published, strings.Join(item.Categories, ","), description, slackWebhook)
				if err != nil {
					if _, ok := err.(*errors.SlackNotificationError); ok {
						log.Printf("main: Failed to send Slack notification: %v", err)
					} else {
						log.Printf("main: Failed to insert data: %v", err)
					}
				}
			}
		}
	}

	if matchFound == 0 {
		fmt.Printf("main: Result: No CVE matches found in the vuldb RSS feed\n")
	}
}
