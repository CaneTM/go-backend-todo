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
)

type TaskService struct {
	serviceName string
}

func NewTaskService() *TaskService {
	return &TaskService{
		serviceName: common.APIRoot + "/tasks",
	}
}

// Gets all tasks from the db
func (ts *TaskService) GetTasks(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Run the query to get all the tasks from the tasks table
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		re := common.NewResponseError(http.StatusInternalServerError, "Error executing query")
		re.WriteError(w)
		return
	}

	defer rows.Close()

	var tasks []reps.Task

	// Iterate through each task that the query returned
	for rows.Next() {
		var t reps.Task
		// Map each column value to the appropriate task field
		if err := rows.Scan(&t.Title, &t.Description, &t.DueDate, &t.Status); err != nil {
			re := common.NewResponseError(http.StatusInternalServerError, "Error reading task values")
			re.WriteError(w)
			return
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		re := common.NewResponseError(http.StatusInternalServerError, "Error reading from database")
		re.WriteError(w)
		return
	}

	// Prepare a JSON response
	resp, err := json.Marshal(tasks)
	if err != nil {
		re := common.NewResponseError(http.StatusInternalServerError, common.FailedToCreateResponseMsg)
		re.WriteError(w)
		return
	}

	w.Header().Set(common.ContentType, common.ApplicationJSON)
	w.Write(resp)
}

// Checks that the required fields are filled
func isValidTask(t reps.Task) bool {
	return t.Title != "" && t.Description != "" && t.Status != ""
}

// Creates a new task and stores it in the db
func (ts *TaskService) CreateTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Get the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		re := common.NewResponseError(http.StatusBadRequest, common.BadRequestMsg)
		re.WriteError(w)
		return
	}

	// Store the info from the request body as a Task
	var t reps.Task
	err = json.Unmarshal(reqBody, &t)
	if err != nil {
		re := common.NewResponseError(http.StatusBadRequest, common.BadRequestMsg)
		re.WriteError(w)
		return
	}

	if !isValidTask(t) {
		re := common.NewResponseError(http.StatusBadRequest, "Missing one or more fields: name, description, status")
		re.WriteError(w)
		return
	}

	// Run the query and handle any errors
	_, err = db.Exec(
		"INSERT INTO tasks (name, description, status, due_date) VALUES (?, ?, ?, ?)",
		t.Title, t.Description, t.DueDate, t.Status,
	)
	if err != nil {
		var re *common.ResponseError
		if strings.Contains(err.Error(), "Duplicate entry") {
			re = common.NewResponseError(http.StatusConflict, "Task name already exists")
		} else {
			re = common.NewResponseError(http.StatusInternalServerError, err.Error())
		}
		re.WriteError(w)
		return
	}

	resp, _ := json.Marshal(fmt.Sprintf("Task %s created", t.Title))
	w.Header().Set(common.ContentType, common.ApplicationJSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// Deletes the task with the given name
func (ts *TaskService) DeleteTask(w http.ResponseWriter, r *http.Request, db *sql.DB, taskName string) {
	var taskToDelete string
	row := db.QueryRow("SELECT name FROM tasks WHERE name = ?", taskName)
	if err := row.Scan(&taskToDelete); err != nil {
		if err == sql.ErrNoRows || taskToDelete == "" {
			re := common.NewResponseError(http.StatusNotFound, fmt.Sprintf("Task %s does not exist", taskName))
			re.WriteError(w)
			return
		}
	}

	_, err := db.Exec("DELETE FROM tasks WHERE name = ?", taskName)
	if err != nil {
		re := common.NewResponseError(http.StatusInternalServerError, fmt.Sprintf("Could not delete task %s", taskName))
		re.WriteError(w)
		return
	}

	resp, _ := json.Marshal(fmt.Sprintf("Task %s deleted", taskName))
	w.Header().Set(common.ContentType, common.ApplicationJSON)
	w.Write(resp)
}

func (ts *TaskService) GetServiceName() string {
	return ts.serviceName
}
