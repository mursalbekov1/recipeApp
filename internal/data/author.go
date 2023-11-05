package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type AuthorModel struct {
	DB *sql.DB
}

func (a AuthorModel) Insert(author *Author) error {
	createTime := time.Now()

	query := `
		INSERT INTO authors (name, email, password, recipeaccesses)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name
	`

	args := []interface{}{author.Name, author.Email, author.Password, pq.Array(author.RecipeAccesses), createTime}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&author.ID)
}

func (a AuthorModel) Get(id int64) (*Author, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT name, email, recipeaccesses
		FROM authors WHERE id = $1 `

	var author Author

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, id).Scan(
		&author.ID,
		&author.Name,
		&author.Email,
		pq.Array(&author.RecipeAccesses),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &author, nil
}

func (a AuthorModel) Update(author *Author) error {

	query := `
		UPDATE authors 
		SET name = $1, email = $2, password = $3, recipeaccesses = $4
		WHERE id = $6
		RETURNING name`

	args := []interface{}{
		author.Name,
		author.Email,
		pq.Array(author.Recipes),
		pq.Array(author.RecipeAccesses),
		author.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, args...).Scan(&author.Name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (a AuthorModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM authors 
       		  WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := a.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
