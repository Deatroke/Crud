package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type userResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Biography string `json:"Biography"`
}

func validateUser(u user) error {
	if len(strings.TrimSpace(u.FirstName)) < 2 || len(u.FirstName) > 20 {
		return errors.New("first_name must be between 2 and 20 characters")
	}
	if len(strings.TrimSpace(u.LastName)) < 2 || len(u.LastName) > 20 {
		return errors.New("last_name must be between 2 and 20 characters")
	}
	if len(strings.TrimSpace(u.Biography)) < 20 || len(u.Biography) > 450 {
		return errors.New("biography must be between 20 and 450 characters")
	}
	return nil
}

func decodeAndValidateUser(r *http.Request) (user, error) {
	var u user
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return u, errors.New("Erro ao decodificar JSON")
	}

	if strings.TrimSpace(u.FirstName) == "" ||
		strings.TrimSpace(u.LastName) == "" ||
		strings.TrimSpace(u.Biography) == "" {
		return u, errors.New("Please provide FirstName LastName and bio for the user")
	}

	if err := validateUser(u); err != nil {
		return u, err
	}

	return u, nil
}

func (app *application) getAllUsers() []userResponse {
	users := make([]userResponse, 0, len(app.data))
	for uid, u := range app.data {
		users = append(users, userResponse{
			ID:        uuid.UUID(uid).String(),
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Biography: u.Biography,
		})
	}
	return users
}

func (app *application) getUserByID(uid id) (user, bool) {
	u, exists := app.data[uid]
	return u, exists
}

func (app *application) updateUserByID(uid id, u user) (user, bool) {
	_, exists := app.data[uid]
	if !exists {
		return user{}, false
	}
	app.data[uid] = u
	return u, true
}

func (app *application) deleteUserByID(uid id) (user, bool) {
	u, exists := app.data[uid]
	if !exists {
		return user{}, false
	}
	delete(app.data, uid)
	return u, true
}
