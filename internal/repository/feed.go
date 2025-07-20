package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/database"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/models"
)

type FeedRepository interface {
	Create(ctx context.Context, feed models.Feed) (int, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]models.Feed, error)
	GetById(ctx context.Context, id int) (models.Feed, error)
	GetByUrl(ctx context.Context, url string) (models.Feed, error)
}

type FeedRepositoryImpl struct {
	db database.Database
}

func NewFeedRepository(db database.Database) FeedRepository {
	return &FeedRepositoryImpl{db: db}
}

func (f FeedRepositoryImpl) Create(ctx context.Context, feed models.Feed) (int, error) {
	var id int
	err := f.db.QueryRow(ctx,
		`INSERT INTO "feeds"("url") VALUES ($1) RETURNING id`,
		feed.Url,
	).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, err
}

func (f FeedRepositoryImpl) Delete(ctx context.Context, id int) error {
	_, err := f.db.Exec(ctx,
		`DELETE FROM "feeds" WHERE id=$1`,
		id,
	)
	return err
}

func (f FeedRepositoryImpl) GetById(ctx context.Context, id int) (models.Feed, error) {
	var feed models.Feed
	err := f.db.QueryRow(ctx,
		`SELECT "id","url" FROM "feeds" WHERE id=$1`,
		id,
	).Scan(&feed.ID, &feed.Url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return feed, models.ErrNoRecord
		}
		return feed, err
	}

	return feed, err
}

func (f FeedRepositoryImpl) GetAll(ctx context.Context) ([]models.Feed, error) {
	var feeds []models.Feed

	rows, err := f.db.Query(ctx, `SELECT "id","url" FROM "feeds"`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var feed models.Feed
		err := rows.Scan(&feed.ID, &feed.Url)
		if err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, err
}

func (f FeedRepositoryImpl) GetByUrl(ctx context.Context, url string) (models.Feed, error) {
	var feed models.Feed
	err := f.db.QueryRow(ctx,
		`SELECT "id","url" FROM "feeds" WHERE url=$1`,
		url,
	).Scan(&feed.ID, &feed.Url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return feed, models.ErrNoRecord
		}
		return feed, err
	}

	return feed, err
}
