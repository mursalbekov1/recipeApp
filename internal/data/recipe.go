package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"go_recipe/internal/validator"
	"time"
)

type RecipeModel struct {
	DB *sql.DB
}

func (r RecipeModel) Insert(recipe *Recipe) error {
	createTime := time.Now()

	query := `
		INSERT INTO recipes (title, description, ingredients, steps, author_id, collaborators, time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, time, version
	`

	args := []interface{}{recipe.Title, recipe.Description, pq.Array(recipe.Ingredients), pq.Array(recipe.Steps), recipe.Author, pq.Array(recipe.Collaborators), createTime}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&recipe.ID, &recipe.Time)
}

func (r RecipeModel) Get(id int64) (*Recipe, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, time, title, description, ingredients, steps 
		FROM recipes WHERE id = $1 `

	var recipe Recipe

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&recipe.ID,
		&recipe.Time,
		&recipe.Title,
		&recipe.Description,
		pq.Array(&recipe.Ingredients),
		pq.Array(&recipe.Steps),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &recipe, nil
}

func (r RecipeModel) Update(recipe *Recipe) error {

	query := `
		UPDATE recipes 
		SET title = $1, description = $2, ingredients = $3, steps = $4, collaborators = $5, version = version + 1
		WHERE id = $6 and version = $7
		RETURNING time, version`

	args := []interface{}{
		recipe.Title,
		recipe.Description,
		pq.Array(recipe.Ingredients),
		pq.Array(recipe.Steps),
		pq.Array(recipe.Collaborators),
		recipe.ID,
		recipe.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&recipe.Time, &recipe.Version)
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

func (r RecipeModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM recipes 
       		  WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, id)
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

func (r RecipeModel) GetAll(title string, ingredients []string, filters Filters) ([]*Recipe, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT  count(*) OVER(), id, title, description, ingredients, steps, author_id, collaborators, version
		FROM recipes
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (ingredients @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, pq.Array(ingredients), filters.limit(), filters.offset()}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	recipes := []*Recipe{}

	for rows.Next() {

		var recipe Recipe

		err := rows.Scan(
			&totalRecords,
			&recipe.ID,
			&recipe.Title,
			&recipe.Description,
			pq.Array(&recipe.Ingredients),
			pq.Array(&recipe.Steps),
			&recipe.Author,
			pq.Array(&recipe.Collaborators),
			&recipe.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		recipes = append(recipes, &recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return recipes, metadata, nil
}

func (r Recipe) MarshalJSON() ([]byte, error) {
	//var runtime string
	//
	//if r.Runtime != 0 {
	//	runtime = fmt.Sprintf("%d mins", r.Runtime)
	//}

	type RecipeAlias Recipe

	aux := struct {
		RecipeAlias
		Runtime string `json:"runtime,omitempty"`
	}{
		RecipeAlias: RecipeAlias(r),
		//Runtime:     runtime,
	}

	return json.Marshal(aux)
}

func ValidateRecipe(v *validator.Validator, recipe *Recipe) {
	v.Check(recipe.Title != "", "title", "must be provided")
	v.Check(len(recipe.Title) <= 500, "title", "must not be more than 500 bytes long")
	//v.Check(recipe.Runtime != 0, "runtime", "must be provided")
	//v.Check(recipe.Runtime > 0, "runtime", "must be a positive integer")
}
