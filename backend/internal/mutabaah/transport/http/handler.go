package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mutabaahapi/internal/mutabaah/usecase"
	"mutabaahapi/internal/platform/apperror"
	"mutabaahapi/internal/transport/http/contextx"
	"mutabaahapi/internal/transport/http/response"
)

type Handler struct {
	InputUC   *usecase.InputMutabaah
	ListUC    *usecase.ListMutabaah
	BulkUC    *usecase.BulkAttendance
	SummaryUC *usecase.MutabaahSummary
}

func NewHandler(u *usecase.InputMutabaah, l *usecase.ListMutabaah, b *usecase.BulkAttendance, s *usecase.MutabaahSummary) *Handler {
	return &Handler{InputUC: u, ListUC: l, BulkUC: b, SummaryUC: s}
}

type inputRequest struct {
	SessionID           string                 `json:"session_id" binding:"required"`
	StudentID           string                 `json:"student_id"` // Optional jika Murid
	Attendance          string                 `json:"attendance"`
	StudentActivityData map[string]interface{} `json:"student_activity_data"`
	Status              string                 `json:"status"`
	Note                string                 `json:"note"`
}

func (h *Handler) Input(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	var req inputRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	res, err := h.InputUC.Execute(c.Request.Context(), usecase.InputData{
		UserID:              userID,
		SessionID:           req.SessionID,
		StudentID:           req.StudentID,
		Attendance:          req.Attendance,
		StudentActivityData: req.StudentActivityData,
		Status:              req.Status,
		Note:                req.Note,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Mutabaah saved", res)
}

func (h *Handler) List(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	// Ambil Query Param: ?session_id=...&student_id=...
	sessionID := c.Query("session_id")
	studentID := c.Query("student_id")

	res, err := h.ListUC.Execute(c.Request.Context(), usecase.ListFilter{
		UserID:    userID,
		SessionID: sessionID,
		StudentID: studentID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Data retrieved", res)
}

func (h *Handler) BulkAttendance(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c) // Asumsi udah ada import contextx
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	var req struct {
		SessionID string                   `json:"session_id" binding:"required"`
		Items     []usecase.AttendanceItem `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	err := h.BulkUC.Execute(c.Request.Context(), usecase.BulkInput{
		UserID:    userID,
		SessionID: req.SessionID,
		Items:     req.Items,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Attendance saved", nil)
}

func (h *Handler) Summary(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	res, err := h.SummaryUC.Execute(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Summary loaded", res)
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/mutabaah", h.Input)
	rg.GET("/mutabaah", h.List)
	rg.POST("/bulk-attendance", h.BulkAttendance)
	rg.GET("/mutabaah/summary", h.Summary)
}
