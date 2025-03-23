package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	// importing sqlite3.
	_ "github.com/mattn/go-sqlite3"
)


type SlackMessage struct {
	Text string `json:"text"`
}

func InitDB(dbPath string) (*sql.DB, error) {
	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("InitDB: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("InitDB: %w", err)
	}

	log.Println("Connected to DB")

	if err := CreateTable(dbConn); err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("InitDB: %w", err)
	}

	return dbConn, nil
}

func CreateTable(dbConn *sql.DB) error {
	createVulnDBTableSQL := `CREATE TABLE IF NOT EXISTS vulndb (
		vuln_title TEXT NOT NULL,
		link TEXT NOT NULL,
		published TEXT NOT NULL,
		categories TEXT NOT NULL,
		PRIMARY KEY (vuln_title, link)
	)`

	statement, err := dbConn.Prepare(createVulnDBTableSQL)
	if err != nil {
		return fmt.Errorf("CreateTable: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("CreateTable: %w", err)
	}

	log.Println("vulndb table created successfully")

	return nil
}

func InsertData(dbConn *sql.DB, vulnTitle string, link string, published string, categories string, description string, slackWebhook string) error {
	insertQuery := `INSERT INTO vulndb (vuln_title, link, published, categories) VALUES (?, ?, ?, ?)`

	statement, err := dbConn.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("InsertData: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(vulnTitle, link, published, categories)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("InsertData: entry already exists in DB, skipping")
		}

		return fmt.Errorf("InsertData: %w", err)
	}

	log.Println("Insert Completed")
	err = notifySlack(vulnTitle, link, published, categories, description, slackWebhook)
	if err != nil {
		log.Println("InsertData: Failed to send Slack notification:", err)
	}

	return nil
}

func notifySlack(vulnTitle string, link string, published string, categories string, description string, slackWebhook string) error {
	re := regexp.MustCompile(`<a href=".*?">(.*?)</a>`)
	description = re.ReplaceAllString(description, "$1")

	description = strings.ReplaceAll(description, "<code>", "*")
	description = strings.ReplaceAll(description, "</code>", "*")

	description = strings.ReplaceAll(description, "<em>", "*")
	description = strings.ReplaceAll(description, "</em>", "*")

	message := SlackMessage{
		Text: "Title: " + vulnTitle + "\nLink: " + link + "\nDate Published: " + published + "\nDescription: " + description,
	}

	// Encode message payload as JSON
	payload, err := json.Marshal(message)
	if err != nil {
		log.Println("notifySlack: Failed to marshal message:", err)
		return fmt.Errorf("notifySlack: Failed to marshal message: %w", err)
	}

	// Make POST request to Slack webhook
	resp, err := http.Post(slackWebhook, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("notifySlack: Failed to send message, check if slack webhook is valid:", err)
		return fmt.Errorf("notifySlack: Failed to send message, check if slack webhook is valid: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("notifySlack: Failed to send message, status code:", resp.StatusCode)
		return fmt.Errorf("notifySlack: Failed to send message, status code: %d", resp.StatusCode)
	}

	log.Println("Message sent successfully!")
	return nil
}
