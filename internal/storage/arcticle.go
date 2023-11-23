package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/vitalis-virtus/news-telegram-bot/internal/model"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

type dbArticle struct {
	ID          int       `db:"id"`
	SourceID    int       `db:"source_id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Summary     string    `db:"summary"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	PostedAt    time.Time `db:"posted_at"`
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

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, link, summary, published_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit int) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticle
	if err := conn.SelectContext(ctx, &articles, `SELECT * FROM articles WHERE posted_at IS NULL AND published_at >= $1::timestamp ORDER BY published_at DESC LIMIT $2`, since.UTC().Format(time.RFC3339), limit); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(a dbArticle, _ int) model.Article {
		return model.Article(a)
	}), nil
}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `UPDATE articles SET posted_at=$1::timestamp WHERE id=$2`, time.Now().UTC().Format(time.RFC3339), id); err != nil {
		return err
	}

	return nil
}
