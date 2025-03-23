package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dark-warlord14/CVENotifier/internal/errors"
	"github.com/dark-warlord14/CVENotifier/internal/slack"

	// importing sqlite3.
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, &errors.DatabaseError{Message: fmt.Sprintf("InitDB: %s", err.Error())}
	}

	if err := dbConn.Ping(); err != nil {
		dbConn.Close()
		return nil, &errors.DatabaseError{Message: fmt.Sprintf("InitDB: %s", err.Error())}
	}

	//log.Println("Connected to DB")

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
		return &errors.DatabaseError{Message: fmt.Sprintf("CreateTable: %s", err.Error())}
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		return &errors.DatabaseError{Message: fmt.Sprintf("CreateTable: %s", err.Error())}
	}

	//log.Println("vulndb table created successfully")

	return nil
}

func InsertData(dbConn *sql.DB, vulnTitle string, link string, published string, categories string, description string, slackWebhook string) error {
	insertQuery := `INSERT INTO vulndb (vuln_title, link, published, categories) VALUES (?, ?, ?, ?)`

	statement, err := dbConn.Prepare(insertQuery)
	if err != nil {
		return &errors.DatabaseError{Message: fmt.Sprintf("InsertData: %s", err.Error())}
	}
	defer statement.Close()

	_, err = statement.Exec(vulnTitle, link, published, categories)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return &errors.DatabaseError{Message: "InsertData: entry already exists in DB, skipping"}
		}

		return &errors.DatabaseError{Message: fmt.Sprintf("InsertData: %s", err.Error())}
	}

	//log.Println("Insert Completed")
	err = slack.NotifySlack(vulnTitle, link, published, categories, description, slackWebhook)
	if err != nil {
		//log.Println("InsertData: Failed to send Slack notification:", err)
		return err
	}

	return nil
}
