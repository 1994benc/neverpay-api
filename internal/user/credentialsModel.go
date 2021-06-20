package user

import (
	"encoding/json"
	"io"
	"net/http"
)

type CredentialsModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *CredentialsModel) FromJSON(body io.ReadCloser) error {
	err := json.NewDecoder(body).Decode(a)
	return err
}

func (a *CredentialsModel) ToJSON(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(a)
}
