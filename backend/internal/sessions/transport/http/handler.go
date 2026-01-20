package http

import (
	"mutabaahapi/internal/platform/apperror"
	"mutabaahapi/internal/sessions/entity"
	"mutabaahapi/internal/sessions/usecase"
	"mutabaahapi/internal/transport/http/contextx"
	"mutabaahapi/internal/transport/http/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	CreateUC      *usecase.CreateSession
	ListUC        *usecase.ListSession
	GetUC         *usecase.GetSession
	UpdateUC      *usecase.UpdateSession
	ListStudentUC *usecase.ListStudentSession
}

// [UPDATE] Constructor terima 2 UseCase
func NewHandler(createUC *usecase.CreateSession, listUC *usecase.ListSession, getUC *usecase.GetSession, updateUC *usecase.UpdateSession, listStudentUC *usecase.ListStudentSession) *Handler {
	return &Handler{
		CreateUC:      createUC,
		ListUC:        listUC,
		GetUC:         getUC,
		UpdateUC:      updateUC,
		ListStudentUC: listStudentUC,
	}
}

type createSessionRequest struct {
	ClassID               string                 `json:"class_id" binding:"required"`
	TemplateID            string                 `json:"template_id" binding:"required"`
	Date                  string                 `json:"date" binding:"required"`
	Name                  string                 `json:"name" binding:"required"`
	ClassActivityData     map[string]interface{} `json:"class_activity_data"`
	IsStudentInputAllowed bool                   `json:"is_student_input_allowed"`
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	var req createSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	res, err := h.CreateUC.Execute(c.Request.Context(), usecase.CreateSessionInput{
		UserID:                userID,
		ClassID:               req.ClassID,
		TemplateID:            req.TemplateID,
		Date:                  req.Date,
		Name:                  req.Name,
		ClassActivityData:     req.ClassActivityData,
		IsStudentInputAllowed: req.IsStudentInputAllowed,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Session created", res)
}

// [BARU] Handler List
func (h *Handler) List(c *gin.Context) {
	classID := c.Query("class_id")
	if classID == "" {
		response.Error(c, apperror.New(apperror.ErrCodeInvalidInput, "class_id is required"))
		return
	}

	list, err := h.ListUC.Execute(c.Request.Context(), classID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Sessions loaded", list)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.GetUC.Execute(c.Request.Context(), id)
	if err != nil {
		// Mapping error "resource not found" ke 404
		if err.Error() == "resource not found" {
			response.Error(c, apperror.ErrNotFound)
			return
		}
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Session detail", res)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	// Gunakan pointer untuk field optional agar bisa partial update
	var req struct {
		Name                  *string                `json:"name"`
		Date                  *string                `json:"date"`
		IsStudentInputAllowed *bool                  `json:"is_student_input_allowed"`
		ClassActivityData     map[string]interface{} `json:"class_activity_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	err := h.UpdateUC.Execute(c.Request.Context(), id, usecase.UpdateSessionInput{
		Name:                  req.Name,
		Date:                  req.Date,
		IsStudentInputAllowed: req.IsStudentInputAllowed,
		ClassActivityData:     req.ClassActivityData,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Session updated", nil)
}

func (h *Handler) ListForStudent(c *gin.Context) {
	actorID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}
	role, _ := c.Get("user_role")

	targetStudentID := actorID
	// Jika Guru yang akses, ambil student_id dari parameter
	if role == "GURU" {
		if q := c.Query("student_id"); q != "" {
			targetStudentID = q
		}
	}

	// [PENTING] Ambil class_id dari query params untuk filter
	classID := c.Query("class_id")

	var sessions []entity.Session
	var err error

	if c.Query("mode") == "report" {
		// Mode report (lihat semua history)
		sessions, err = h.ListStudentUC.ExecuteReport(c.Request.Context(), targetStudentID, classID)
	} else {
		// Mode dashboard (hanya aktif)
		sessions, err = h.ListStudentUC.Execute(c.Request.Context(), targetStudentID, classID)
	}

	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Student sessions retrieved", sessions)
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/sessions", h.Create)
	rg.GET("/sessions", h.List)
	rg.GET("/sessions/student", h.ListForStudent)
	rg.GET("/sessions/:id", h.GetByID)
	rg.PUT("/sessions/:id", h.Update)
}
