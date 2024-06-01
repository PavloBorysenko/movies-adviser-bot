package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"movies-adviser-bot/storage"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
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

// Save saves movie to storage.
func (s *Storage) Save(ctx context.Context, m *storage.Movie) error {
	q := `INSERT INTO movies (title, user_name, type) VALUES (?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, m.Title, m.UserName, m.Type); err != nil {
		return fmt.Errorf("can't save movie: %w", err)
	}

	return nil
}
// PickLast picks last film or series from storage.
func (s *Storage) GetAll(ctx context.Context, userName string, ftype string) ([]*storage.Movie, error) {
	q := `SELECT title FROM movies WHERE user_name = ? AND type = ? ORDER BY Timestamp`

	var muvies []*storage.Movie

	rows, err := s.db.QueryContext(ctx, q, userName, ftype)
	if err != nil {
		return nil, fmt.Errorf("can't get all  %w: %w", ftype, err)
	}
	defer rows.Close()
	for rows.Next() {
		var muvie string
		if err := rows.Scan(&muvie); err != nil {
			return nil, fmt.Errorf("can't parse film  %w: %w", ftype, err)
		}
		muvies = append(
			muvies, 
			&storage.Movie{
				Title:      muvie,
				UserName: userName,
				Type: ftype,
			},
		)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get all  %w: %w", ftype, err)
	}

	return muvies, nil
}
// PickLast picks last film or series from storage.
func (s *Storage) PickLast(ctx context.Context, userName string, ftype string) (*storage.Movie, error) {
	q := `SELECT title FROM movies WHERE user_name = ? AND type = ? ORDER BY Timestamp LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName, ftype).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedMovies
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick Last %w: %w", ftype, err)
	}

	return &storage.Movie{
		Title:      url,
		UserName: userName,
		Type: ftype,
	}, nil
}
// FindOne find film or series from storage.
func (s *Storage) FindOne(ctx context.Context, userName string, searchText string, ftype string) (*storage.Movie, error) {
	q := `SELECT title FROM movies WHERE user_name = ? AND type = ? AND LOWER(title) LIKE '%' || ? || '%'  ORDER BY Timestamp LIMIT 1`

	var muvie string
	
	err := s.db.QueryRowContext(ctx, q, userName, ftype, searchText).Scan(&muvie)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNofoundMovies
	}
	if err != nil {
		return nil, fmt.Errorf("there is no such  %w: %w", ftype, err)
	}

	return &storage.Movie{
		Title:      muvie,
		UserName: userName,
		Type: ftype,
	}, nil
}

// PickRandom picks random movie from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Movie, error) {
	q := `SELECT title, type FROM movies WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var muvie storage.Movie

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&muvie.Title, &muvie.Type)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedMovies
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random : %w", err)
	}

	return &storage.Movie{
		Title:      muvie.Title,
		UserName: userName,
		Type:  muvie.Type,
		
	}, nil
}

// Remove removes movie from storage.
func (s *Storage) Remove(ctx context.Context, movie *storage.Movie) error {
	q := `DELETE FROM movies WHERE title = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, q, movie.Title, movie.UserName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

// IsExists checks if movie exists in storage.
func (s *Storage) IsExists(ctx context.Context, movie *storage.Movie) (bool, error) {
	q := `SELECT COUNT(*) FROM movies WHERE title = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, movie.Title, movie.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS movies (title TEXT, user_name TEXT, type TEXT, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
