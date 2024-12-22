package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/canetm/go-backend-todo/services"
)

const (
	// Patterns
	usernamePattern = "username"
)

type userHandler struct {
	db  *sql.DB
	svc *services.UserService
}

func NewUserHandler(db *sql.DB) *userHandler {
	return &userHandler{
		db:  db,
		svc: services.NewUserService(),
	}
}

func (uh *userHandler) HandleService() {
	svcRoot := uh.svc.GetServiceName()
	http.HandleFunc(svcRoot, uh.handleUsers)
	http.HandleFunc(fmt.Sprintf("%s/{%s}", svcRoot, usernamePattern), uh.handleUsersUsername)
}

func (uh *userHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uh.svc.GetUsers(w, r, uh.db)
	case http.MethodPost:
		uh.svc.CreateUser(w, r, uh.db)
	}
}

func (uh *userHandler) handleUsersUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue(usernamePattern)
	switch r.Method {
	// case http.MethodGet:
	// 	uh.userService.GetUserByUsername(w, r, uh.db)
	case http.MethodDelete:
		uh.svc.DeleteUser(w, r, uh.db, username)
	}
}
