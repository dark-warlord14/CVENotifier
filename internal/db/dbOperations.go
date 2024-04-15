package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	// importing sqlite3.
	_ "github.com/mattn/go-sqlite3"
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type SlackMessage struct {
	Text string `json:"text"`
}

func InitDB(dbPath string) (*sql.DB, error) {
	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("exception occurred: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("exception occurred: %w", err)
	}

	log.Println("Connected to DB")

	if err := CreateTable(dbConn); err != nil {
		dbConn.Close()
		return nil, err
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
		return fmt.Errorf("exception occurred: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("exception occurred: %w", err)
	}

	log.Println("vulndb table created successfully")

	return nil
}

func InsertData(dbConn *sql.DB, vulnTitle string, link string, published string, categories string, slackWebhook string) error {
	insertQuery := `INSERT INTO vulndb (vuln_title, link, published, categories) VALUES (?, ?, ?, ?)`

	statement, err := dbConn.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("exception occurred: %w", err)
	}
	defer statement.Close()

	_, err = statement.Exec(vulnTitle, link, published, categories)
	if err != nil {
		return fmt.Errorf("exception occurred: %w", err)
	}

	log.Println("Insert Completed")
	notifySlack(vulnTitle, link, published, categories, slackWebhook)

	return nil
}

func notifySlack(vulnTitle string, link string, published string, categories string, slackWebhook string) {
	message := SlackMessage{
		Text: "Title: " + vulnTitle + "\nLink: " + link + "\nDate Published: " + published + "\n" + strings.ReplaceAll(categories,
			",", "\n"),
	}

	// Encode message payload as JSON
	payload, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshal message:", err)
		return
	}

	// Make POST request to Slack webhook
	resp, err := http.Post(slackWebhook, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Failed to send message:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to send message, status code:", resp.StatusCode)
		return
	}

	log.Println("Message sent successfully!")
}
