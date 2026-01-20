package http

import (
	"mutabaahapi/internal/platform/apperror"
	"mutabaahapi/internal/templates/entity"
	"mutabaahapi/internal/templates/usecase"
	"mutabaahapi/internal/transport/http/contextx"
	"mutabaahapi/internal/transport/http/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{ UC *usecase.UseCase }

func NewHandler(uc *usecase.UseCase) *Handler { return &Handler{UC: uc} }

func (h *Handler) Create(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	var req entity.Template
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	if err := h.UC.Create(c.Request.Context(), userID, &req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "Template created", req)
}

func (h *Handler) List(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	list, err := h.UC.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Templates loaded", list)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req entity.Template
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	if err := h.UC.Update(c.Request.Context(), id, &req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Template updated", nil)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.UC.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Template deleted", nil)
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/templates")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
