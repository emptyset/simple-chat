package models

import (
	"encoding/json"
	"fmt"
	"github.com/emptyset/simple-chat/internal/storage"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

type UnixTime struct {
	time.Time
}

func (t *UnixTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	t.Time = time.Unix(i, 0)
	return nil
}

func (t *UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", t.Time.UTC().Unix())), nil
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type Message struct {
	Id          int         `json:"id"`
	Timestamp   UnixTime    `json:"timestamp"`
	SenderId    int         `json:"sender_id"`
	RecipientId int         `json:"recipient_id"`
	Content     string      `json:"content"`
	MediaType   string      `json:"media_type"`
	Metadata    interface{} `json:"metadata"`
}

// media types for metadata _type field
const (
	Text  string = "text"
	Image string = "image"
	Video string = "video"
)

type ChatModel interface {
	SignupUser(username string, password string) (*User, error)
	GetMessages(senderId int, recipientId int, count int, offset int) ([]Message, error)
	SendMessage(senderId int, recipientId int, content string, mediaType string) error
}

type Model struct {
	store storage.DataStore
}

func New(store storage.DataStore) *Model {
	return &Model{
		store: store,
	}
}

func (m *Model) SignupUser(username string, password string) (*User, error) {
	// TODO: check for existing username

	log.Debug("generating hash for password")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("error when generating hash for password: %s", err)
		return nil, err
	}

	log.Debug("creating user in the data store")
	record, err := m.store.CreateUser(username, hash)
	if err != nil {
		log.Errorf("error when creating user in the data store: %s", err)
		return nil, err
	}

	log.Debug("transcribing user from record")
	return transcribeUser(record)
}

func transcribeUser(record storage.Record) (*User, error) {
	user := &User{}
	err := json.Unmarshal(record, user)
	return user, err
}

func (m *Model) GetMessages(senderId int, recipientId int, count int, offset int) ([]Message, error) {
	var messages []Message

	log.Debug("reading messages from the data store")
	records, err := m.store.ReadMessages(senderId, recipientId, count, offset)
	if err != nil {
		return messages, err
	}

	log.Debug("transcribing messages from records")
	for _, record := range records {
		message, err := transcribeMessage(record)
		if err != nil {
			log.Errorf("error when transcribing message from record: %s", err)
			continue
		}

		messages = append(messages, *message)
	}

	return messages, nil
}

func (m *Model) SendMessage(senderId int, recipientId int, content string, mediaType string) error {
	log.Debug("getting metadata based on content and media type")
	metadata, err := getMetadata(content, mediaType)
	if err != nil {
		return err
	}

	log.Debug("creating message in the data store")
	_, err = m.store.CreateMessage(senderId, recipientId, content, mediaType, metadata)
	return err
}

func transcribeMessage(record storage.Record) (*Message, error) {
	message := &Message{
		Metadata: &json.RawMessage{},
	}
	err := json.Unmarshal(record, message)
	return message, err
}

func getMetadata(content string, mediaType string) (map[string]string, error) {
	var metadata map[string]string
	switch mediaType {
	case Text:
		break
	case Image:
		getImageMetadata(content, metadata)
	case Video:
		getVideoMetadata(content, metadata)
	default:
		return metadata, fmt.Errorf("unknown content media type: %s", mediaType)
	}

	return metadata, nil
}

func getImageMetadata(content string, metadata map[string]string) {
	// could parse content and query service for metadata
	metadata["width"] = "200px"
	metadata["height"] = "200px"
}

func getVideoMetadata(content string, metadata map[string]string) {
	// could parse content and query service for metadata
	metadata["length"] = "9h3m0.5s"
	metadata["source"] = "YouTube"
}
