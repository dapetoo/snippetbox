package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dapetoo/snippetbox/pkg/models"
	"time"
)

// SnippetModel Define a SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert This will insert a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	query := `INSERT INTO snippets (title, content, created, expires) 
			VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	//DB.Exec()
	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}
	//Use
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Get This will return a specific snippet from the database
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	query := `
			SELECT id, title, content, created, expires 
			FROM snippets 
			WHERE expires > UTC_TIMESTAMP() AND id = ?
			`

	row := m.DB.QueryRow(query, id)

	//Initialize a pointer to a new zeroed snippet struct
	s := &models.Snippet{}

	//Use row.Scan() to copy the values from each field in sql.Row to the corresponding filed in the struct
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	//Return the snippet object
	return s, nil
}

// GetByID Another implementation of Get
func (m *SnippetModel) GetByID(id int) (*models.Snippet, error) {
	s := &models.Snippet{}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err := m.DB.QueryRowContext(ctx, `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?
	`, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	//Return the snippet object
	return s, nil
}

// Latest This will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	query := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP() 
		ORDER BY created DESC
		LIMIT 10
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := make([]*models.Snippet, 0)

	for rows.Next() {
		s := &models.Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
