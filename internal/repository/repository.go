package repository

import (
	"fmt"
	"strconv"

	"github.com/LoaltyProgramm/to-do-app/internal/models"

	"github.com/jmoiron/sqlx"
)

type TaskRepository interface {
	AddTask(task models.Task) (int64, error)
	GetTasks(limit int) ([]models.Task, error)
	GetTask(id string) (models.Task, error)
	UpdateTask(task models.Task) error
	UpdateDateTask(task models.Task) error
	DeleteTask(id string) error
	SearchTasksDates(date string, limit int) ([]models.Task, error)
	SearchTasks(data string, limit int) ([]models.Task, error)
}

type taskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}	
}

func (r *taskRepository) AddTask(task models.Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat)
	VALUES (?, ?, ?, ?);
	`

	res, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("error insert task: %w", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error get last id: %w", err)
	}

	return id, err
}

func (r *taskRepository) GetTasks(limit int) ([]models.Task, error) {
	query := `SELECT * FROM scheduler ORDER BY date ASC LIMIT ?;`

	tasks := []models.Task{}

	err := r.db.Select(&tasks, query, limit)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) GetTask(id string) (models.Task, error) {
	query := `SELECT * FROM scheduler WHERE id=?;`

	task := models.Task{}

	if id == "" {
		return models.Task{}, fmt.Errorf("id not specified")
	}

	err := r.db.Get(&task, query, id)
	if err != nil {
		return models.Task{}, fmt.Errorf("there is no given task: %w", err)
	}

	return task, nil
}

func (r *taskRepository) UpdateTask(task models.Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? 
	WHERE id = ?;`

	if task.Id == "" {
		return fmt.Errorf("id not specified")
	}

	taskInt, err := strconv.Atoi(task.Id)
	if err != nil {
		return fmt.Errorf("the Id can only be a number: %w", err)
	}

	if taskInt > 1000000 {
		return fmt.Errorf("the ID is too large")
	}

	if task.Comment == "" || task.Title == "" {
		return fmt.Errorf("required fields must be filled in")
	}

	_, err = r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		return fmt.Errorf("error update task: %w", err)
	}

	return nil
}

func (r *taskRepository) UpdateDateTask(task models.Task) error {
	query := `UPDATE scheduler SET date = ? 
	WHERE id = ?;`

	if task.Id == "" {
		return fmt.Errorf("id not specified")
	}

	taskInt, err := strconv.Atoi(task.Id)
	if err != nil {
		return fmt.Errorf("the Id can only be a number: %w", err)
	}

	if taskInt > 1000000 {
		return fmt.Errorf("the ID is too large")
	}

	if task.Title == "" {
		return fmt.Errorf("required fields must be filled in")
	}

	_, err = r.db.Exec(query, task.Date, task.Id)
	if err != nil {
		return fmt.Errorf("error update task: %w", err)
	}

	return nil
}

func (r *taskRepository) DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("fail delete task for id")
	}

	return nil
}

func (r *taskRepository) SearchTasksDates(date string, limit int) ([]models.Task, error) {
	query := `SELECT * FROM scheduler WHERE date = ? LIMIT ?;`

	var tasks []models.Task

	err := r.db.Select(&tasks, query, date, limit)
	if err != nil {
		return nil, fmt.Errorf("the request to get rows by date could not be completed: %w", err)
	}

	return tasks, nil
}

func (r *taskRepository) SearchTasks(data string, limit int) ([]models.Task, error) {
	query := `SELECT * FROM scheduler WHERE title LIKE '%' || ? || '%'
			OR comment LIKE '%' || ? ||'%' LIMIT ?;`

	var tasks []models.Task

	err := r.db.Select(&tasks, query, data, data, limit)
	if err != nil {
		return nil, fmt.Errorf("the request to get rows by date could not be completed: %w", err)
	}

	return tasks, nil
}
