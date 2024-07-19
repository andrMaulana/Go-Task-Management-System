package handler

import (
	"net/http"
	"strconv"

	"github.com/andrMaulana/Go-Task-Management-System/internal/models"
	"github.com/andrMaulana/Go-Task-Management-System/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Meta: models.Meta{
				Code:    http.StatusBadRequest,
				Status:  "Bad Request",
				Message: "Invalid input data",
			},
		})
		return
	}

	projectID, _ := strconv.Atoi(c.Param("projectId"))
	task.ProjectID = uint(projectID)

	if err := h.taskService.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Meta: models.Meta{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to create task",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Meta: models.Meta{
			Code:    http.StatusCreated,
			Status:  "Created",
			Message: "Task created successfully",
		},
		Data: task,
	})
}

func (h *TaskHandler) GetTasksByProject(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("projectId"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	tasks, total, err := h.taskService.GetTasksByProject(uint(projectID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Meta: models.Meta{
				Code:    http.StatusInternalServerError,
				Status:  "Internal Server Error",
				Message: "Failed to get tasks",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Meta: models.Meta{
			Code:    http.StatusOK,
			Status:  "OK",
			Message: "Tasks retrieved successfully",
		},
		Data: gin.H{
			"tasks":     tasks,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
