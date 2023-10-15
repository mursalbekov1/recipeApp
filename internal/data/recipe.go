package data

import (
	"encoding/json"
	"fmt"
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
