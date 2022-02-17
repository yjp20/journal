package main

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Media struct {
	ID            int64      `json:"id"`
	Description   string     `json:"description"`
	MediaType     string     `json:"media_type"`
	Rating        *float32   `json:"rating"`
	Notes         *string    `json:"notes"`
	RelatedLink   *string    `json:"related_link"`
	Comments      *string    `json:"comments"`
	Cart          bool       `json:"cart"`
	Completed     bool       `json:"completed"`
	Progress      float64    `json:"progress"`
	CompletedDate *time.Time `json:"completed_date"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type MediaModel struct {
	DB *sql.DB
}

func (m MediaModel) Insert(media *Media) error {
	query := `
		INSERT INTO media (description, media_type, related_link, comments, cart, completed, progress)
		VALUES ($1, $2, $3, $4, false, false, 0)
		RETURNING id, created_at, updated_at`
	return m.DB.QueryRow(query, media.Description, media.MediaType, media.RelatedLink, media.Comments).Scan(&media.ID, &media.CreatedAt, &media.UpdatedAt)
}

func (m MediaModel) Get(id int64) (*Media, error) {
	query := `
		SELECT
			id,
			description,
			media_type,
			rating,
			notes,
			related_link,
			comments,
			cart,
			completed,
			progress,
			completed_date,
			created_at,
			updated_at
		FROM media
		WHERE id = $1`

	var media Media

	err := m.DB.QueryRow(query, id).Scan(
		&media.ID,
		&media.Description,
		&media.MediaType,
		&media.Rating,
		&media.Notes,
		&media.RelatedLink,
		&media.Comments,
		&media.Cart,
		&media.Completed,
		&media.Progress,
		&media.CompletedDate,
		&media.CreatedAt,
		&media.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &media, nil
}

func (m MediaModel) GetAll() ([]*Media, error) {
	query := `
		SELECT
			id,
			description,
			media_type,
			rating,
			notes,
			related_link,
			comments,
			cart,
			completed,
			progress,
			completed_date,
			created_at,
			updated_at
		FROM media
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	medias := []*Media{}
	for rows.Next() {
		var media Media
		err := rows.Scan(
			&media.ID,
			&media.Description,
			&media.MediaType,
			&media.Rating,
			&media.Notes,
			&media.RelatedLink,
			&media.Comments,
			&media.Cart,
			&media.Completed,
			&media.Progress,
			&media.CompletedDate,
			&media.CreatedAt,
			&media.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		medias = append(medias, &media)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return medias, nil
}

func (m MediaModel) Update(media *Media) error {
	query := `
		UPDATE media
		SET
			description = $1,
			media_type = $2,
			rating = $3,
			notes = $4,
			related_link = $5,
			comments = $6,
			cart = $7,
			completed = $8,
			progress = $9,
			completed_date = $10
		WHERE id = $11`

	_, err := m.DB.Exec(
		query,
		&media.Description,
		&media.MediaType,
		&media.Rating,
		&media.Notes,
		&media.RelatedLink,
		&media.Comments,
		&media.Cart,
		&media.Completed,
		&media.Progress,
		&media.CompletedDate,
		&media.ID,
	)

	return err
}

func (m MediaModel) Delete(id int64) error {
	query := `
		DELETE FROM media
		WHERE id = $1`

	_, err := m.DB.Exec(query, id)
	return err
}
