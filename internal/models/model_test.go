package models_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

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
	senderId := 1
	recipientId := 2
	count := 2
	offset := 0
	store := &mocks.DataStore{}
	now := 1524365200
	later := 1524365205
	var first storage.Record = []byte(fmt.Sprintf(`{"id":1, "timestamp":%d, "sender_id": %d, "recipient_id": %d, "content": "hello", "media_type": "text", "metadata": {}}`, now, senderId, recipientId))
	var second storage.Record = []byte(fmt.Sprintf(`{"id":2, "timestamp":%d, "sender_id": %d, "recipient_id": %d, "content": "hello", "media_type": "text", "metadata": {}}`, later, senderId, recipientId))
	records := []storage.Record{second, first}
	store.On("ReadMessages", senderId, recipientId, count, offset).Return(records, nil)
	model := models.New(store)

	messages, err := model.GetMessages(senderId, recipientId, count, offset)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)

	assert.Equal(t, messages[0].Id, 2)
	assert.Equal(t, messages[1].Id, 1)

	store.AssertExpectations(t)
}

func TestSendMessage(t *testing.T) {
	senderId := 1
	recipientId := 2
	content := "howdy"
	mediaType := "text"
	var metadata map[string]string
	store := &mocks.DataStore{}
	var record storage.Record = []byte("anything")
	store.On("CreateMessage", senderId, recipientId, content, mediaType, metadata).Return(record, nil)
	model := models.New(store)

	err := model.SendMessage(senderId, recipientId, content, mediaType)
	assert.NoError(t, err)

	store.AssertExpectations(t)
}
