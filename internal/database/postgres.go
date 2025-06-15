package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
)

type pool interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type Postgres struct {
	pool pool
}

func NewPostgres(pool pool) *Postgres {
	return &Postgres{pool}
}

func (p *Postgres) CreateItem(ctx context.Context, title, link, description string) error {
	_, err := p.pool.Exec(ctx,
		`INSERT INTO "items"("title", "link", "description") VALUES ($1, $2, $3)`,
		title, link, description,
	)
	return err
}

func (p *Postgres) DeleteAll(ctx context.Context) error {
	_, err := p.pool.Exec(ctx, `DELETE FROM "items";`)
	return err
}
