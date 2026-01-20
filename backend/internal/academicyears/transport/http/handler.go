package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"mutabaahapi/internal/academicyears/usecase"
	"mutabaahapi/internal/platform/apperror"
	"mutabaahapi/internal/transport/http/contextx"
	"mutabaahapi/internal/transport/http/response"
)

type Handler struct {
	UC     *usecase.UseCase
	Logger zerolog.Logger
}

func NewHandler(uc *usecase.UseCase, l zerolog.Logger) *Handler {
	return &Handler{UC: uc, Logger: l}
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	var req struct {
		Name      string `json:"name" binding:"required"`
		StartDate string `json:"start_date" binding:"required"`
		EndDate   string `json:"end_date" binding:"required"`
		IsActive  bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	input := usecase.CreateInput{
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  req.IsActive,
	}

	// Pass userID
	if err := h.UC.Create(c.Request.Context(), userID, input); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "Academic year created", nil)
}

func (h *Handler) List(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	res, err := h.UC.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "List retrieved", res)
}

func (h *Handler) Update(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}
	id := c.Param("id")
	var req struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		IsActive  bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	input := usecase.UpdateInput{
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  req.IsActive,
	}

	if err := h.UC.Update(c.Request.Context(), userID, id, input); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Updated", nil)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.UC.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Deleted", nil)
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.POST("/academic-years", h.Create)
	r.GET("/academic-years", h.List)
	r.PUT("/academic-years/:id", h.Update)
	r.DELETE("/academic-years/:id", h.Delete)
}
