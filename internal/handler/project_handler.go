package handler

import (
	"net/http"
	"strconv"

	"github.com/andrMaulana/Go-Task-Management-System/internal/models"
	"github.com/andrMaulana/Go-Task-Management-System/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Meta: models.Meta{
				Code:    http.StatusBadRequest,
				Status:  "Bad Request",
				Message: "Invalid input data",
			},
		})
		return
	}

	userIDFloat := c.MustGet("user_id").(float64)
	project.OwnerID = uint(userIDFloat)

	if err := h.projectService.CreateProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Meta: models.Meta{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to create project",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Meta: models.Meta{
			Code:    http.StatusCreated,
			Status:  "Created",
			Message: "Project created successfully",
		},
		Data: project,
	})
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	userID := c.MustGet("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	projects, total, err := h.projectService.GetProjects(userID.(float64), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Meta: models.Meta{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to get projects",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Meta: models.Meta{
			Code:    http.StatusOK,
			Status:  "OK",
			Message: "Projects retrieved successfully",
		},
		Data: gin.H{
			"projects":  projects,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *ProjectHandler) ShareProject(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("projectId"))
	var shareData struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&shareData); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Meta: models.Meta{
				Code:    http.StatusBadRequest,
				Status:  "Bad Request",
				Message: "Invalid input data",
			},
		})
		return
	}

	if err := h.projectService.ShareProject(uint(projectID), shareData.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Meta: models.Meta{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to share project",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Meta: models.Meta{
			Code:    http.StatusOK,
			Status:  "OK",
			Message: "Project shared successfully",
		},
	})
}
