package mysql

import (
	"database/sql"
	"github.com/dapetoo/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// Insert to add a new record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate method to verify that a user exists with the provided email address
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get fetch a details for a specific user
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
