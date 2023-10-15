package data

import (
	"encoding/json"
	"fmt"
	"go_recipe/internal/validator"
)

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
	//v.Check(recipe.Year != 0, "year", "must be provided")
	//v.Check(recipe.Year >= 1888, "year", "must be greater than 1888")
	//v.Check(recipe.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(recipe.Runtime != 0, "runtime", "must be provided")
	v.Check(recipe.Runtime > 0, "runtime", "must be a positive integer")
	//v.Check(recipe.Genres != nil, "genres", "must be provided")
	//v.Check(len(recipe.Genres) >= 1, "genres", "must contain at least 1 genre")
	//v.Check(len(recipe.Genres) <= 5, "genres", "must not contain more than 5 genres")
	//v.Check(validator.Unique(recipe.Genres), "genres", "must not contain duplicate values")
}
