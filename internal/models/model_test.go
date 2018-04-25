package models_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/emptyset/simple-chat/internal/mocks"
	"github.com/emptyset/simple-chat/internal/models"
	"github.com/emptyset/simple-chat/internal/storage"
)

func TestSignupUser(t *testing.T) {
	username := "emptyset"
	password := "correct horse battery staple"
	store := &mocks.DataStore{}
	var record storage.Record = []byte(fmt.Sprintf(`{"id":1, "username":"%s"}`, username))
	store.On("CreateUser", username, mock.Anything).Return(record, nil)
	model := models.New(store)

	user, err := model.SignupUser(username, password)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, username)

	store.AssertExpectations(t)
}

func TestGetMessages(t *testing.T) {
	senderID := 1
	recipientID := 2
	count := 2
	offset := 0
	store := &mocks.DataStore{}
	now := 1524365200
	later := 1524365205
	var first storage.Record = []byte(fmt.Sprintf(`{"id":1, "timestamp":%d, "sender_id": %d, "recipient_id": %d, "content": "hello", "media_type": "text", "metadata": {}}`, now, senderID, recipientID))
	var second storage.Record = []byte(fmt.Sprintf(`{"id":2, "timestamp":%d, "sender_id": %d, "recipient_id": %d, "content": "hello", "media_type": "text", "metadata": {}}`, later, senderID, recipientID))
	records := []storage.Record{second, first}
	store.On("ReadMessages", senderID, recipientID, count, offset).Return(records, nil)
	model := models.New(store)

	messages, err := model.GetMessages(senderID, recipientID, count, offset)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)

	assert.Equal(t, messages[0].ID, 2)
	assert.Equal(t, messages[1].ID, 1)

	store.AssertExpectations(t)
}

func TestSendMessage(t *testing.T) {
	senderID := 1
	recipientID := 2
	content := "howdy"
	mediaType := "text"
	var metadata map[string]string
	store := &mocks.DataStore{}
	var record storage.Record = []byte("anything")
	store.On("CreateMessage", senderID, recipientID, content, mediaType, metadata).Return(record, nil)
	model := models.New(store)

	err := model.SendMessage(senderID, recipientID, content, mediaType)
	assert.NoError(t, err)

	store.AssertExpectations(t)
}
