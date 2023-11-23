package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vitalis-virtus/news-telegram-bot/internal/model"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit int) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return nil, nil
}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}
