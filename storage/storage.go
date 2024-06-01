package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"movies-adviser-bot/lib/e"
)

type Storage interface {
	Save(ctx context.Context, p *Movie) error
	PickRandom(ctx context.Context, userName string) (*Movie, error)
	PickLast(ctx context.Context, userName string, ftype string) (*Movie, error)
	FindOne(ctx context.Context, userName string, searchText string, ftype string) (*Movie, error)
	GetAll(ctx context.Context, userName string, ftype string) ([]*Movie, error)
	Remove(ctx context.Context, p *Movie) error
	IsExists(ctx context.Context, p *Movie) (bool, error)
}

var ErrNoSavedMovies = errors.New("no saved movies")
var ErrNofoundMovies = errors.New("no movies found")

type Movie struct {
	Title      string
	UserName string
	Type 	 string
}

func (m Movie) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, m.Title ); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, m.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
