package mysql

import (
	"database/sql"
	"errors"
	"github.com/dapetoo/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

// Insert to add a new record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	//Use Bcrypt to hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 24)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, hashed_password, created)
			  VALUES (?,?,?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(query, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
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
