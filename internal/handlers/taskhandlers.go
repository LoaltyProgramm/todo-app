package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/LoaltyProgramm/to-do-app/internal/config"
	"github.com/LoaltyProgramm/to-do-app/internal/models"
	"github.com/LoaltyProgramm/to-do-app/internal/service/authservice"
	"github.com/LoaltyProgramm/to-do-app/internal/service/cookieservice"
	"github.com/LoaltyProgramm/to-do-app/internal/service/middlewareservice"
	"github.com/LoaltyProgramm/to-do-app/internal/service/taskservice"
	"github.com/LoaltyProgramm/to-do-app/internal/utils"
)

type TaskHandler struct {
	taskService       service.TasksService
	authService       authservice.JWTService
	cookieService     cookieservice.CookieService
	middlewareService middlewareservice.MiddlewareService
	cfg               config.Config
}

func NewTaskHandlers(taskService service.TasksService, cfg *config.Config, authService authservice.JWTService, cookieService cookieservice.CookieService, middlewareService middlewareservice.MiddlewareService) TaskHandler {
	return TaskHandler{
		taskService:       taskService,
		cfg:               *cfg,
		authService:       authService,
		cookieService:     cookieService,
		middlewareService: middlewareService,
	}
}

func (h *TaskHandler) NextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(utils.Layout, nowStr)
	if err != nil {
		fmt.Printf("Error parse now time: %v", err)
		return
	}

	next, err := utils.NextDate(now, date, repeat)
	if err != nil {
		fmt.Printf("Couldn't find the next date: %v", err)
		return
	}

	w.Write([]byte(next))
}

func (h *TaskHandler) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	task, err := h.taskService.GetTaskByID(id)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	utils.WriteJson(w, task, http.StatusOK)
}

func (h *TaskHandler) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parse body", http.StatusBadRequest)
	}
	defer r.Body.Close()

	var task models.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": "Unmarshal error"}, http.StatusBadRequest)
		return
	}

	id, err := h.taskService.CreateTask(task)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, map[string]any{"id": id}, http.StatusOK)
}

func (h *TaskHandler) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var task models.Task

	err = json.Unmarshal(body, &task)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	err = h.taskService.UpdateTask(task)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, task, http.StatusOK)
}

func (h *TaskHandler) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := h.taskService.DeleteTaskByID(id)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, struct{}{}, http.StatusOK)
}

func (h *TaskHandler) TaskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		h.getTaskHandler(w, r)
	case http.MethodPost:
		h.createTaskHandler(w, r)
	case http.MethodPut:
		h.updateTaskHandler(w, r)
	case http.MethodDelete:
		h.deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) TasksHandler(w http.ResponseWriter, r *http.Request) {
	valueSearch := r.FormValue("search")
	limit := 50
	if valueSearch == "" {
		tasks, err := h.taskService.ListTasks(limit)
		if err != nil {
			utils.WriteJson(w, map[string]string{"error": "Error select from table"}, http.StatusInternalServerError)
			return
		}

		utils.WriteJson(w, models.TasksResp{Tasks: tasks}, http.StatusOK)
		return
	}

	if valueSearch != "" {
		if utils.CheckingTheDateUsingATemplate(valueSearch) {
			layoutInput := "02.01.2006"
			dateTime, _ := time.Parse(layoutInput, valueSearch)
			date := dateTime.Format(utils.Layout)

			tasks, err := h.taskService.FindTasksByDate(date, 50)
			if err != nil {
				utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
				return
			}

			utils.WriteJson(w, models.TasksResp{Tasks: tasks}, http.StatusOK)
			return
		}

		tasks, err := h.taskService.SearchTasks(valueSearch, 50)
		if err != nil {
			utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}

		utils.WriteJson(w, models.TasksResp{Tasks: tasks}, http.StatusOK)
		return
	}
}

func (h *TaskHandler) CompletedHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	task, err := h.taskService.GetTaskByID(id)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err := h.taskService.DeleteTaskByID(id)
		if err != nil {
			utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		}
		utils.WriteJson(w, struct{}{}, http.StatusOK)
	}

	if task.Repeat != "" {
		nextDate, err := utils.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		task.Date = nextDate

		err = h.taskService.UpdateTaskDate(task)
		if err != nil {
			utils.WriteJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		utils.WriteJson(w, struct{}{}, http.StatusOK)
	}
}

func (h *TaskHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": "Error parse body"}, http.StatusBadRequest)
		return
	}

	password := models.TaskRequestPassword{}

	err = json.Unmarshal(body, &password)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": "Error unmarshal body"}, http.StatusInternalServerError)
		return
	}

	if password.Password != h.cfg.TodoPassword {
		utils.WriteJson(w, map[string]string{"error": "Password invalid"}, http.StatusBadRequest)
		return
	}

	token, err := h.authService.CreateToken(password.Password)
	if err != nil {
		utils.WriteJson(w, map[string]string{"error": "Fail create token"}, http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, map[string]string{"token": token}, http.StatusOK)
}

func (h *TaskHandler) InitHandler() {
	http.HandleFunc("/api/nextdate", h.middlewareService.MiddlewareAuth(h.NextDayHandler))
	http.HandleFunc("/api/task", h.middlewareService.MiddlewareAuth(h.TaskHandler))
	http.HandleFunc("/api/tasks", h.middlewareService.MiddlewareAuth(h.TasksHandler))
	http.HandleFunc("/api/task/done", h.middlewareService.MiddlewareAuth(h.CompletedHandler))
	http.HandleFunc("/api/signin", h.RegistrationHandler)
}
