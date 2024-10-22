package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type Category struct {
	UserID     int64
	Title      string
	Total      int
	TotalSpent int
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS category 
		 (
		 	 id INTEGER PRIMARY KEY AUTOINCREMENT,
			 user_id INTEGER NOT NULL,
			 title TEXT NOT NULL,
			 total INTEGER NOT NULL,
			 total_spent INTEGER DEFAULT 0
		 )`

	_, err := s.db.ExecContext(ctx, q)

	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}

func (s *Storage) IsExists(ctx context.Context, c *Category) (bool, error) {
	q := `SELECT COUNT (*) FROM category WHERE title = ? AND user_id = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, c.Title, c.UserID).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if category exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Save(ctx context.Context, c *Category) error {
	q := `INSERT INTO category (user_id, title, total) VALUES (?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, c.UserID, c.Title, c.Total); err != nil {
		return fmt.Errorf("can't save data: %w", err)
	}
	return nil
}

func (s *Storage) Update(ctx context.Context, c *Category) error {
	q := `UPDATE category SET total_spent = ? WHERE user_id = ? AND title = ?`

	alreadySpent, err := s.getTotalSpent(ctx, c)
	if err != nil {
		return fmt.Errorf("can't get category spent: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, q, alreadySpent+c.TotalSpent, c.UserID, c.Title); err != nil {
		return fmt.Errorf("can't update category: %w", err)
	}
	return nil
}

func (s *Storage) Get(ctx context.Context, UserID int64) ([]Category, error) {
	q := `SELECT * FROM category WHERE user_id = ?`

	rows, err := s.db.QueryContext(ctx, q, UserID)
	if err != nil {
		return nil, fmt.Errorf("can't return query: %w", err)
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var cat Category
		var id int
		if err := rows.Scan(&id, &cat.UserID, &cat.Title, &cat.Total, &cat.TotalSpent); err != nil {
			return nil, fmt.Errorf("can't return rows: %w", err)
		}
		categories = append(categories, cat)
	}
	if err = rows.Err(); err != nil {
		return categories, err
	}
	return categories, nil
}

func (s *Storage) getTotalSpent(ctx context.Context, c *Category) (int, error) {
	q := `SELECT total_spent FROM category WHERE user_id = ? AND title = ?`

	row := s.db.QueryRowContext(ctx, q, c.UserID, c.Title)
	var totalSpent int
	err := row.Scan(&totalSpent)

	if err != nil {
		return 0, err
	}
	return totalSpent, nil
}

func (s *Storage) Delete(ctx context.Context, c *Category) error {
	q := `DELETE FROM category WHERE title = ? AND user_id = ?`

	if _, err := s.db.ExecContext(ctx, q, c.Title, c.UserID); err != nil {
		return fmt.Errorf("can't remove category: %w", err)
	}

	return nil
}

func (s *Storage) DeleteAll(ctx context.Context, c *Category) error {
	q := `DELETE FROM category WHERE user_id = ?`

	if _, err := s.db.ExecContext(ctx, q, c.UserID); err != nil {
		return fmt.Errorf("can't remove: %w", err)
	}

	return nil
}
