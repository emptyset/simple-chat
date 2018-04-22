package models_test

import (
	"github.com/stretchr/testify/assert"
	"fmt"
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
