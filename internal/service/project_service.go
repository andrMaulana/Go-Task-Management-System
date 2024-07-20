package service

import (
	"errors"

	"github.com/andrMaulana/Go-Task-Management-System/internal/models"
	"gorm.io/gorm"
)

type ProjectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

func (s *ProjectService) CreateProject(project *models.Project) error {
	if project.Name == "" {
		return errors.New("project name cannot be empty")
	}
	return s.db.Create(project).Error
}

func (s *ProjectService) GetProjects(userID float64, page, pageSize int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	offset := (page - 1) * pageSize

	query := s.db.Model(&models.Project{}).Where("owner_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (s *ProjectService) GetProjectByID(projectID uint) (*models.Project, error) {
	var project models.Project
	if err := s.db.First(&project, projectID).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectService) UpdateProject(project *models.Project) error {
	if project.Name == "" {
		return errors.New("project name cannot be empty")
	}
	return s.db.Save(project).Error
}

func (s *ProjectService) DeleteProject(projectID uint) error {
	return s.db.Delete(&models.Project{}, projectID).Error
}

func (s *ProjectService) ShareProject(projectID, userID uint) error {
	project, err := s.GetProjectByID(projectID)
	if err != nil {
		return err
	}

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	return s.db.Model(project).Association("Users").Append(&user)
}
