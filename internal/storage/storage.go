package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Record []byte

type DataStore interface {
	CreateUser(username string, hash []byte) (Record, error)
	CreateMessage(senderID int, recipientID int, content string, mediaType string, metadata map[string]string) (Record, error)
	ReadMessages(senderID int, recipientID int, count int, offset int) ([]Record, error)
}

type SQLDataStore struct {
	database *sql.DB
}

func NewSQLDataStore(database *sql.DB) *SQLDataStore {
	return &SQLDataStore{
		database: database,
	}
}

func (s *SQLDataStore) CreateUser(username string, hash []byte) (Record, error) {
	record := []byte("")

	log.Debug("preparing insert user statement")
	statement, err := s.database.Prepare("INSERT INTO user (username, hash) VALUES (?, ?)")
	if err != nil {
		return record, err
	}

	log.Debug("executing statement")
	response, err := statement.Exec(username, hash)
	if err != nil {
		return record, err
	}

	id, err := response.LastInsertId()
	if err != nil {
		return record, err
	}

	record = []byte(fmt.Sprintf(`{"id": %d, "username": "%s"}`, id, username))

	log.WithFields(log.Fields{
		"record": string(record),
	}).Debug("record to return")

	return record, nil
}

func (s *SQLDataStore) CreateMessage(senderID int, recipientID int, content string, mediaType string, metadata map[string]string) (Record, error) {
	record := []byte("")

	log.Debug("preparing insert message statement")
	statement, err := s.database.Prepare("INSERT INTO message (timestamp, sender_id, recipient_id, content, media_type, metadata) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return record, err
	}

	encodedMetadata, err := json.Marshal(metadata)
	if err != nil {
		return record, err
	}

	log.Debug("executing statement")
	response, err := statement.Exec(time.Now().UTC().Unix(), senderID, recipientID, content, mediaType, encodedMetadata)
	if err != nil {
		return record, err
	}

	id, err := response.LastInsertId()
	if err != nil {
		return record, err
	}

	record = []byte(fmt.Sprintf(`{"id": %d}`, id))

	log.WithFields(log.Fields{
		"record": string(record),
	}).Debug("record to return")

	return record, nil
}

func (s *SQLDataStore) ReadMessages(senderID int, recipientID int, count int, offset int) ([]Record, error) {
	var records []Record
	rows, err := s.database.Query("SELECT id, timestamp, content, media_type, metadata FROM message WHERE sender_id <> recipient_id AND sender_id in (?, ?) AND recipient_id in (?, ?) ORDER BY timestamp DESC LIMIT ? OFFSET ?", senderID, recipientID, senderID, recipientID, count, offset)
	if err != nil {
		return records, err
	}
	defer rows.Close()

	var (
		id           int
		rawTimestamp []byte
		content      string
		mediaType    string
		rawMetadata  []byte
	)

	for rows.Next() {
		err := rows.Scan(&id, &rawTimestamp, &content, &mediaType, &rawMetadata)
		if err != nil {
			log.Errorf("unable to scan row into record")
		}

		// TODO: verify we are parsing timestamp and metadata correctly here

		record := []byte(fmt.Sprintf(`{"id": %d, "timestamp": "%s", "sender_id": %d, "recipient_id": %d, "content": "%s", "media_type": "%s", "metadata": "%s"}`, id, rawTimestamp, senderID, recipientID, content, mediaType, rawMetadata))

		log.WithFields(log.Fields{
			"record": string(record),
		}).Debug("record to return")

		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return records, err
	}

	return records, nil
}
