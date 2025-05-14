package service

import (
	"fmt"

	"github.com/LoaltyProgramm/to-do-app/internal/models"
	"github.com/LoaltyProgramm/to-do-app/internal/repository"
	"github.com/LoaltyProgramm/to-do-app/internal/utils"
	"strconv"
)

type TasksService interface {
	CreateTask(task models.Task) (int64, error)
	ListTasks(limit int) ([]models.Task, error)
	GetTaskByID(id string) (models.Task, error)
	UpdateTask(task models.Task) error
	UpdateTaskDate(task models.Task) error
	DeleteTaskByID(id string) error
	FindTasksByDate(date string, limit int) ([]models.Task, error)
	SearchTasks(data string, limit int) ([]models.Task, error)
}

type tasksSevice struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TasksService {
	return &tasksSevice{
		repo: repo,
	}
}

func (s *tasksSevice) CreateTask(task models.Task) (int64, error) {
	if task.Title == "" {
		return 0, fmt.Errorf("the title cannot be empty")
		
	}

	err := utils.CheckDate(&task)
	if err != nil {
		return 0, fmt.Errorf("the date was not validated") 
		
	}

	return s.repo.AddTask(task)
}

func (s *tasksSevice) ListTasks(limit int) ([]models.Task, error) {
	return s.repo.GetTasks(limit)
}

func (s *tasksSevice) GetTaskByID(id string) (models.Task, error) {
		return s.repo.GetTask(id)
}

func (s *tasksSevice) UpdateTask(task models.Task) error {
	return s.repo.UpdateTask(task)
}

func (s *tasksSevice) UpdateTaskDate(task models.Task) error {
	return s.repo.UpdateDateTask(task)
}

func (s *tasksSevice) DeleteTaskByID(id string) error {
	if id == "" {
		return fmt.Errorf("id is not empty")
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid id. The ID is passed as an integer")
		
	}

	return s.repo.DeleteTask(id)
}

func (s *tasksSevice) FindTasksByDate(date string, limit int) ([]models.Task, error) {
	return s.repo.SearchTasksDates(date, limit)
}

func (s *tasksSevice) SearchTasks(data string, limit int) ([]models.Task, error) {
	return s.repo.SearchTasks(data, limit)
}

