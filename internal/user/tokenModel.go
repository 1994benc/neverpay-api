package user

import (
	"encoding/json"
	"net/http"
)

type TokenModel struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

func (t *TokenModel) ToJSON(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(t)
}
