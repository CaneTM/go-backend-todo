package handlers

//rc 12/22/24
import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/canetm/go-backend-todo/services"
)

const (
	// Patterns
	taskPattern = "taskname"
)

type taskHandler struct {
	db  *sql.DB
	svc *services.TaskService
}

func NewTaskHandler(db *sql.DB) *taskHandler {
	return &taskHandler{
		db:  db,
		svc: services.NewTaskService(),
	}
}

func (th *taskHandler) HandleService() {
	svcRoot := th.svc.GetServiceName()
	http.HandleFunc(svcRoot, th.handleTasks)
	http.HandleFunc(fmt.Sprintf("%s/{%s}", svcRoot, taskPattern), th.handleTasksByTaskName)
}

func (th *taskHandler) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		th.svc.GetTasks(w, r, th.db)
	case http.MethodPost:
		th.svc.CreateTask(w, r, th.db)
	}
}

func (th *taskHandler) handleTasksByTaskName(w http.ResponseWriter, r *http.Request) {
	taskName := r.URL.Query().Get(taskPattern)
	switch r.Method {
	//Not sure if this first case needs to be implemented, since it was ignored for users as well
	// case http.MethodGet:
	// 	th.svc.GetTaskByTitle(w, r, th.db, taskName)

	case http.MethodDelete:
		th.svc.DeleteTask(w, r, th.db, taskName)

		// Could implement to make a dynamic task
		// case http.MethodPut:
		// 	th.svc.UpdateTask(w, r, th.db, taskName)
	}
}
