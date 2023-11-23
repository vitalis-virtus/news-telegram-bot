package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/vitalis-virtus/news-telegram-bot/internal/model"
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

type dbSource struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var sources []dbSource
	if err := conn.SelectContext(ctx, &sources, `SELECT * FROM sources`); err != nil {
		return nil, err
	}

	return lo.Map(sources, func(s dbSource, _ int) model.Source {
		return model.Source(s)
	}), nil
}

func (s *SourcePostgresStorage) SourceByID(ctx context.Context, id int) (*model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var source dbSource
	if err := conn.GetContext(ctx, &source, `SELECT * FROM sources WHERE id = $1`, id); err != nil {
		return nil, err
	}

	return (*model.Source)(&source), nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source model.Source) (int, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var id int

	row := conn.QueryRowContext(
		ctx,
		`INSERT INTO resources (name, feed_url, created_at) VALUES ($1, $2, $3) RETURNING id`,
		source.Name,
		source.FeedURL,
		source.CreatedAt,
	)
	if err := row.Err(); err != nil { // query was not successfull
		return 0, err
	}

	if err := row.Scan(&id); err != nil { // scan id value
		return 0, err
	}

	return id, nil
}

func (s *SourcePostgresStorage) DeleteByID(ctx context.Context, id int) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `DELETE FROM resources WHERE id = $1`, id); err != nil {
		return err
	}

	return nil
}
