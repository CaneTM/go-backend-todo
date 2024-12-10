package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/canetm/go-backend-todo/common"
	"github.com/canetm/go-backend-todo/reps"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	serviceName string
}

func NewUserService() *UserService {
	return &UserService{
		serviceName: common.APIRoot + "/users",
	}
}

// Gets all users from the db
func (us *UserService) GetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Run the query to get all the users from the users table
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		re := common.NewResponseError(http.StatusInternalServerError, "Error executing query")
		re.WriteError(w)
		return
	}

	defer rows.Close()

	var users []reps.User

	// Iterate through each user that the query returned
	for rows.Next() {
		var u reps.User
		// Map each column value to the appropriate user field
		if err := rows.Scan(&u.ID, &u.Username, &u.Hash); err != nil {
			re := common.NewResponseError(http.StatusInternalServerError, "Error reading user values")
			re.WriteError(w)
			return
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		re := common.NewResponseError(http.StatusInternalServerError, "Error reading from database")
		re.WriteError(w)
		return
	}

	// Prepare a JSON response
	resp, err := json.Marshal(users)
	if err != nil {
		re := common.NewResponseError(http.StatusInternalServerError, common.FailedToCreateResponseMsg)
		re.WriteError(w)
		return
	}

	w.Header().Set(common.ContentType, common.ApplicationJSON)
	w.Write(resp)
}

// Checks that the required fields are filled
func isValidUser(u reps.User) bool {
	return u.Username != "" && u.Hash != ""
}

// Creates a new user and stores it in the db
func (us *UserService) CreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		re := common.NewResponseError(http.StatusBadRequest, common.BadRequestMsg)
		re.WriteError(w)
		return
	}

	// Store the info from the request body as a User
	var u reps.User
	err = json.Unmarshal(reqBody, &u)
	if err != nil {
		re := common.NewResponseError(http.StatusBadRequest, common.BadRequestMsg)
		re.WriteError(w)
		return
	}

	if !isValidUser(u) {
		re := common.NewResponseError(http.StatusBadRequest, "Missing one or more fields: username, hash")
		re.WriteError(w)
		return
	}

	// The password is still in plaintext, so we hash it first
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Hash), bcrypt.DefaultCost)
	if err != nil {
		re := common.NewResponseError(http.StatusInternalServerError, "Error creating user")
		re.WriteError(w)
		return
	}
	u.Hash = string(hash)

	// Run the query and handle any errors
	_, err = db.Exec(
		"INSERT INTO users (username, pwhash) VALUES (?, ?)",
		u.Username,
		u.Hash,
	)
	if err != nil {
		var re *common.ResponseError
		if strings.Contains(err.Error(), "Duplicate entry") {
			re = common.NewResponseError(http.StatusConflict, "Username already exists")
		} else {
			re = common.NewResponseError(http.StatusInternalServerError, err.Error())
		}
		re.WriteError(w)
		return
	}

	resp, _ := json.Marshal(fmt.Sprintf("User %s created", u.Username))
	w.Header().Set(common.ContentType, common.ApplicationJSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (us *UserService) GetServiceName() string {
	return us.serviceName
}
