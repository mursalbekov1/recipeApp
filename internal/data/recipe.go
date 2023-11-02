package data

import (
	"database/sql"
	"encoding/json"
	"errors"
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
		RETURNING id, time
	`

	args := []interface{}{recipe.Title, recipe.Description, pq.Array(recipe.Ingredients), pq.Array(recipe.Steps), recipe.Author, pq.Array(recipe.Collaborators), createTime}

	return r.DB.QueryRow(query, args...).Scan(&recipe.ID, &recipe.Time)
}

func (r RecipeModel) Get(id int64) (*Recipe, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, time, title, description, ingredients, steps 
		FROM recipes WHERE id = $1 `

	var recipe Recipe

	err := r.DB.QueryRow(query, id).Scan(
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
		SET title = $1, description = $2, ingredients = $3, steps = $4, collaborators = $5
		WHERE id = $6
		RETURNING time`

	args := []interface{}{
		recipe.Title,
		recipe.Description,
		pq.Array(recipe.Ingredients),
		pq.Array(recipe.Steps),
		pq.Array(recipe.Collaborators),
		recipe.ID,
	}

	return r.DB.QueryRow(query, args...).Scan(&recipe.Time)
}

func (r RecipeModel) Delete(id int64) error {
	return nil
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
