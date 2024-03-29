package main

import (
	"context"
	"database/sql"
	"time"
)

type FeedItem struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	MediaType   string    `json:"media_type"`
	RelatedLink *string   `json:"related_link"`
	Comments    *string   `json:"comments"`
	SourceID    int64     `json:"source_id"`
	Added       bool      `json:"added"`
	PostDate    time.Time `json:"created_at"`
}

type FeedItemModel struct {
	DB *sql.DB
}

func (m FeedItemModel) Get(id int) (*FeedItem, error) {
	query := `
		SELECT
			id,
			description,
			media_type,
			related_link,
			comments,
			source_id,
			post_date,
			added
		FROM feed_items
		WHERE id = $1
	`

	f := &FeedItem{}
	row := m.DB.QueryRow(query, id)
	err := row.Scan(&f.ID, &f.Description, &f.MediaType, &f.RelatedLink, &f.Comments, &f.SourceID, &f.PostDate, &f.PostDate)
	return f, err
}

func (m FeedItemModel) Insert(feedItem *FeedItem) error {
	query := `
		INSERT INTO feed_items (description, media_type, related_link, comments, source_id, post_date, added)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (related_link) DO NOTHING
		RETURNING id`
	return m.DB.QueryRow(query, feedItem.Description, feedItem.MediaType, feedItem.RelatedLink, feedItem.Comments, feedItem.SourceID, feedItem.PostDate, feedItem.Added).Scan(&feedItem.ID)
}

func (m FeedItemModel) Update(feedItem *FeedItem) error {
	query := `
		UPDATE feed_items
		SET
			description = $2,
			media_type = $3,
			related_link = $4,
			comments = $5,
			source_id = $6,
			post_date = $7,
			added = $8
		WHERE id = $1`

	_, err := m.DB.Exec(
		query,
		feedItem.ID,
		feedItem.Description,
		feedItem.MediaType,
		feedItem.RelatedLink,
		feedItem.Comments,
		feedItem.SourceID,
		feedItem.PostDate,
		feedItem.Added,
	)

	return err
}

func (m FeedItemModel) GetAll(start time.Time, end time.Time) ([]*FeedItem, error) {
	query := `
		SELECT
			id,
			description,
			media_type,
			related_link,
			comments,
			source_id,
			post_date,
			added
		FROM feed_items
		WHERE post_date >= $1 AND post_date < $2
		ORDER BY post_date
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, &start, &end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feedItems := []*FeedItem{}
	for rows.Next() {
		var feedItem FeedItem
		err := rows.Scan(
			&feedItem.ID,
			&feedItem.Description,
			&feedItem.MediaType,
			&feedItem.RelatedLink,
			&feedItem.Comments,
			&feedItem.SourceID,
			&feedItem.PostDate,
			&feedItem.Added,
		)
		if err != nil {
			return nil, err
		}
		feedItems = append(feedItems, &feedItem)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feedItems, nil
}
