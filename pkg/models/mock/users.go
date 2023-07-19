package mock

import (
	"github.com/dapetoo/snippetbox/pkg/models"
	"time"
)

var mockUser = &models.User{
	ID:      1,
	Name:    "Peter",
	Email:   "abc@example.com",
	Created: time.Now(),
	Active:  true,
}

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "someone@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
