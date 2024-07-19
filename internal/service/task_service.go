package service

import (
	"errors"

	"github.com/andrMaulana/Go-Task-Management-System/internal/models"
	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) CreateTask(task *models.Task) error {
	if task.Title == "" {
		return errors.New("task title cannot be empty")
	}
	return s.db.Create(task).Error
}

func (s *TaskService) GetTasksByProject(projectID uint, page, pageSize int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	offset := (page - 1) * pageSize

	query := s.db.Model(&models.Task{}).Where("project_id = ?", projectID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (s *TaskService) UpdateTask(task *models.Task) error {
	if task.Title == "" {
		return errors.New("task title cannot be empty")
	}
	return s.db.Save(task).Error
}

func (s *TaskService) DeleteTask(taskID uint) error {
	return s.db.Delete(&models.Task{}, taskID).Error
}
