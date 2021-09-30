package main

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type FeedSource struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Success     bool   `json:"success"`
}

type FeedSourceModel struct {
	DB *sql.DB
}

func (m FeedSourceModel) Insert(feedSource *FeedSource) error {
	query := `
		INSERT INTO feed_sources (description, url)
		VALUES ($1, $2)
		RETURNING id`
	return m.DB.QueryRow(query, feedSource.Description, feedSource.URL).Scan(&feedSource.ID)
}

func (m FeedSourceModel) Get(id int64) (*FeedSource, error) {
	query := `
		SELECT
			id,
			description,
			url
		FROM feed_sources
		WHERE id = $1`

	var feedSource FeedSource

	err := m.DB.QueryRow(query, id).Scan(
		&feedSource.ID,
		&feedSource.Description,
		&feedSource.URL,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &feedSource, nil
}

func (m FeedSourceModel) GetAll() ([]*FeedSource, error) {
	query := `
		SELECT
			id,
			description,
			url
		FROM feed_sources
		ORDER BY description`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feedSources := []*FeedSource{}
	for rows.Next() {
		var feedSource FeedSource
		err := rows.Scan(
			&feedSource.ID,
			&feedSource.Description,
			&feedSource.URL,
		)
		if err != nil {
			return nil, err
		}
		feedSources = append(feedSources, &feedSource)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feedSources, nil
}

func (m FeedSourceModel) Update(feedSource *FeedSource) error {
	query := `
		UPDATE feed_sources
		SET
			description = $2,
			url = $3
		WHERE id = $1`

	_, err := m.DB.Exec(
		query,
		&feedSource.ID,
		&feedSource.Description,
		&feedSource.URL,
	)

	return err
}

func (m FeedSourceModel) Delete(id int64) error {
	query := `
		DELETE FROM feed_sources
		WHERE id = $1`

	_, err := m.DB.Exec(query, id)
	return err
}
