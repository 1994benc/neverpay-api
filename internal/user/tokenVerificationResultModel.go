package user

import (
	"encoding/json"
	"io"
	"net/http"
)

type TokenVerificationResultModel struct {
	IsTokenValid bool `json:"isTokenValid"`
}

func (t *TokenVerificationResultModel) FromJSON(body io.ReadCloser) error {
	err := json.NewDecoder(body).Decode(t)
	return err
}

func (t *TokenVerificationResultModel) ToJSON(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(t)
}
