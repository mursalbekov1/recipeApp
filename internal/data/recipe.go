package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_recipe/internal/validator"
)

type RecipeModel struct {
	DB *sql.DB
}

func (r RecipeModel) Insert(movie *Recipe) error {
	return nil
}

func (r RecipeModel) Get(id int64) (*Recipe, error) {
	return nil, nil
}

func (r RecipeModel) Update(movie *Recipe) error {
	return nil
}

func (r RecipeModel) Delete(id int64) error {
	return nil
}

func (r Recipe) MarshalJSON() ([]byte, error) {
	var runtime string

	if r.Runtime != 0 {
		runtime = fmt.Sprintf("%d mins", r.Runtime)
	}

	type RecipeAlias Recipe

	aux := struct {
		RecipeAlias
		Runtime string `json:"runtime,omitempty"`
	}{
		RecipeAlias: RecipeAlias(r),
		Runtime:     runtime,
	}

	return json.Marshal(aux)
}

func ValidateRecipe(v *validator.Validator, recipe Recipe) {
	v.Check(recipe.Title != "", "title", "must be provided")
	v.Check(len(recipe.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(recipe.Runtime != 0, "runtime", "must be provided")
	v.Check(recipe.Runtime > 0, "runtime", "must be a positive integer")
}
