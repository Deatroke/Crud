package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type id uuid.UUID

type user struct {
	FirstName string
	LastName  string
	Biography string
}

type application struct {
	data map[id]user
}

func NewApplication() *application {
	return &application{
		data: make(map[id]user),
	}
}

func NewHandler(app *application) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Use(jsonMiddleware)

		r.Post("/api/users", app.postUserHandler)
		r.Get("/api/users", app.getAllUsersHandler)
		r.Get("/api/users/{id}", app.getUserHandler)
		r.Put("/api/users/{id}", app.putUserHandler)
		r.Delete("/api/users/{id}", app.deleteUserHandler)
	})

	return r
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (app *application) postUserHandler(w http.ResponseWriter, r *http.Request) {
	u, err := decodeAndValidateUser(r)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Please provide FirstName LastName and bio for the user")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			writeJSONError(w, http.StatusInternalServerError, "There was an error while saving the user to the database")
		}
	}()

	newID := id(uuid.New())
	app.data[newID] = u

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse{
		ID:        uuid.UUID(newID).String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Biography: u.Biography,
	})
}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			writeJSONError(w, http.StatusInternalServerError, "The users information could not be retrieved")
		}
	}()

	users := app.getAllUsers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	paramID := chi.URLParam(r, "id")
	uidParsed, err := uuid.Parse(paramID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "The user with the specified ID does not exist")
	}
	uid := id(uidParsed)

	userData, ok := app.getUserByID(uid)
	if !ok {
		writeJSONError(w, http.StatusInternalServerError, "The users information could not be retrieved")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	paramID := chi.URLParam(r, "id")
	uidParsed, err := uuid.Parse(paramID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "The user with the specified ID does not exist")
	}
	uid := id(uidParsed)

	deletedUser, ok := app.deleteUserByID(uid)
	if !ok {
		writeJSONError(w, http.StatusInternalServerError, "The user could not be removed")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"id":        uuid.UUID(uid).String(),
		"FirstName": deletedUser.FirstName,
		"LastName":  deletedUser.LastName,
		"Biography": deletedUser.Biography,
	})
}

func (app *application) putUserHandler(w http.ResponseWriter, r *http.Request) {
	paramID := chi.URLParam(r, "id")
	uidParsed, err := uuid.Parse(paramID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "The user with the specified ID does not exist")
	}
	uid := id(uidParsed)

	var updatedUser user
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "The user information could not be modified")
		return
	}

	if err := validateUser(updatedUser); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Please provide name and bio for the user")
	}

	result, ok := app.updateUserByID(uid, updatedUser)
	if ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"id":        uuid.UUID(uid).String(),
			"FirstName": updatedUser.FirstName,
			"LastName":  updatedUser.LastName,
			"Biography": updatedUser.Biography,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
