// cmd/main.go
package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/dark-warlord14/CVENotifier/internal/db"
	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Keywords     []string `yaml:"keywords"`
	SlackWebhook []string `yaml:"slackWebhook"`
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Path to the configuration YAML file")
	flag.Parse()

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v.\nPlease provide the config file using -config flag.\ne.g. go run cmd/CVENotifier/main.go -config config.yaml", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to unmarshal config data: %v", err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://vuldb.com/?rss.recent")
	// feed, _ := fp.ParseURL("https://cvefeed.io/rssfeed/latest.xml")

	if feed == nil {
		log.Fatalf("Failed to parse RSS feed: %v. Please retry", err)
	}

	databasePath := "CVENotifier.db"
	dbConn, err := db.InitDB(databasePath)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer dbConn.Close()

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

				if err := db.InsertData(dbConn, item.Title, item.Link, item.Published, strings.Join(item.Categories, ","), cfg.SlackWebhook[0]); err != nil {
					log.Printf("Result: %v", err)
				}
			}
		}
	}

	if matchFound == 0 {
		log.Printf("Result: No CVE matches found in the vuldb RSS feed")
	}
}
