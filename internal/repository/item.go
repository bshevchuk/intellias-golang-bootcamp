package repository

import (
	"context"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/database"
	"github.com/bshevchuk/intellias-golang-bootcamp/internal/models"
)

type ItemRepository interface {
	Create(ctx context.Context, feed models.Item) error
	DeleteAll(ctx context.Context) error
	GetAllInFeed(ctx context.Context, feedId int) ([]models.Item, error)
}

type ItemRepositoryImpl struct {
	db database.Database
}

func NewItemRepository(db database.Database) ItemRepository {
	return &ItemRepositoryImpl{db: db}
}

func (i ItemRepositoryImpl) Create(ctx context.Context, item models.Item) error {
	_, err := i.db.Exec(ctx,
		`INSERT INTO "items"("feedId", "title", "link", "description") VALUES ($1, $2, $3, $4)`,
		item.FeedID, item.Title, item.Link, item.Description,
	)
	return err
}

func (i ItemRepositoryImpl) DeleteAll(ctx context.Context) error {
	_, err := i.db.Exec(ctx, `DELETE FROM "items";`)
	return err
}

func (i ItemRepositoryImpl) GetAllInFeed(ctx context.Context, feedId int) ([]models.Item, error) {
	var items []models.Item

	rows, err := i.db.Query(ctx,
		`SELECT "title", "link", "description" FROM "items" WHERE "feedId"=$1;`,
		feedId,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.Title, &item.Link, &item.Description)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
