package main

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Todo struct {
	ID            int64      `json:"id"`
	Description   string     `json:"description"`
	Cart          bool       `json:"cart"`
	Completed     bool       `json:"completed"`
	Blocked       bool       `json:"blocked"`
	Progress      float64    `json:"progress"`
	DueDate       *time.Time `json:"due_date"`
	CompletedDate *time.Time `json:"completed_date"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Private       bool       `json:"private"`
	Recur         *int       `json:"recur"`
}

type TodoModel struct {
	DB *sql.DB
}

func (m TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todos (description, due_date, cart, completed, progress, private, blocked, recur)
		VALUES ($1, $2, false, false, 0, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	return m.DB.QueryRow(
		query,
		todo.Description,
		todo.DueDate,
		todo.Private,
		todo.Blocked,
		todo.Recur,
	).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
}

func (m TodoModel) Get(id int64) (*Todo, error) {
	query := `
		SELECT
			id,
			description,
			cart,
			completed,
			progress,
			due_date,
			completed_date,
			created_at,
			updated_at,
			private,
			blocked,
			recur
		FROM todos
		WHERE id = $1`

	var todo Todo

	err := m.DB.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Description,
		&todo.Cart,
		&todo.Completed,
		&todo.Progress,
		&todo.DueDate,
		&todo.CompletedDate,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.Private,
		&todo.Blocked,
		&todo.Recur,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &todo, nil
}

func (m TodoModel) GetAll() ([]*Todo, error) {
	query := `
		SELECT
			id,
			description,
			cart,
			completed,
			progress,
			due_date,
			completed_date,
			created_at,
			updated_at,
			private,
			blocked,
			recur
		FROM todos
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []*Todo{}
	for rows.Next() {
		var todo Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Description,
			&todo.Cart,
			&todo.Completed,
			&todo.Progress,
			&todo.DueDate,
			&todo.CompletedDate,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.Private,
			&todo.Blocked,
			&todo.Recur,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (m TodoModel) Update(todo *Todo) error {
	query := `
		UPDATE todos
		SET
			description = $1,
			cart = $2,
			completed = $3,
			progress = $4,
			due_date = $5,
			completed_date = $6,
			private = $7,
			blocked = $8,
			recur = $9
		WHERE id = $10`

	_, err := m.DB.Exec(
		query,
		&todo.Description,
		&todo.Cart,
		&todo.Completed,
		&todo.Progress,
		&todo.DueDate,
		&todo.CompletedDate,
		&todo.Private,
		&todo.Blocked,
		&todo.Recur,
		&todo.ID,
	)

	return err
}

func (m TodoModel) Delete(id int64) error {
	query := `
		DELETE FROM todos
		WHERE id = $1`

	_, err := m.DB.Exec(query, id)
	return err
}
