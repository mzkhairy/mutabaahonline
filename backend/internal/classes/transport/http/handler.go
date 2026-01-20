package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"mutabaahapi/internal/classes/entity"
	"mutabaahapi/internal/classes/usecase"
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

// POST /classes
func (h *Handler) Create(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	// Struct sementara untuk validasi input JSON
	var req struct {
		Name           string `json:"name" binding:"required"`
		Level          string `json:"level" binding:"required"`
		TeacherID      string `json:"teacher_id" binding:"required"`
		AcademicYearID string `json:"academic_year_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	// [FIX] Map ke Entity menggunakan Pointer
	// Agar aman, kita buat variable penampung dulu jika string kosong
	teacherID := req.TeacherID
	academicYearID := req.AcademicYearID

	input := entity.Class{
		Name:           req.Name,
		Level:          req.Level,
		TeacherID:      &teacherID,
		AcademicYearID: &academicYearID,
	}

	if err := h.UC.Create(c.Request.Context(), userID, &input); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Class created", nil)
}

// GET /admin/classes
func (h *Handler) ListClasses(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}
	ayID := c.Query("academic_year_id")
	res, err := h.UC.ListByInstitution(c.Request.Context(), userID, ayID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Classes retrieved", res)
}

// GET /guru/classes
func (h *Handler) ListByTeacher(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}
	ayID := c.Query("academic_year_id")
	res, err := h.UC.ListByTeacher(c.Request.Context(), userID, ayID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Teacher classes loaded", res)
}

// GET /student/classes
func (h *Handler) ListByStudent(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}
	res, err := h.UC.ListByStudent(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Student classes loaded", res)
}

// GET /classes/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.UC.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Class detail", res)
}

// PUT /classes/:id
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	// Gunakan struct req yang sama dengan Create untuk binding
	var req struct {
		Name           string `json:"name"`
		Level          string `json:"level"`
		TeacherID      string `json:"teacher_id"`
		AcademicYearID string `json:"academic_year_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	teacherID := req.TeacherID
	academicYearID := req.AcademicYearID

	input := entity.Class{
		Name:           req.Name,
		Level:          req.Level,
		TeacherID:      &teacherID,
		AcademicYearID: &academicYearID,
	}

	if err := h.UC.Update(c.Request.Context(), id, &input); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Class updated", nil)
}

// DELETE /classes/:id
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.UC.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Class deleted", nil)
}

// POST /classes/:id/students
func (h *Handler) AddStudent(c *gin.Context) {
	classID := c.Param("id")
	var req struct {
		StudentID string `json:"student_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}
	if err := h.UC.AddStudent(c.Request.Context(), classID, req.StudentID); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Student added to class", nil)
}

// DELETE /classes/:id/students/:student_id
func (h *Handler) RemoveStudent(c *gin.Context) {
	classID := c.Param("id")
	studentID := c.Param("student_id")
	if err := h.UC.RemoveStudent(c.Request.Context(), classID, studentID); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Student removed from class", nil)
}

// GET /classes/:id/students
func (h *Handler) GetStudents(c *gin.Context) {
	classID := c.Param("id")
	students, err := h.UC.ListStudents(c.Request.Context(), classID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Students loaded", students)
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/classes")
	g.POST("", h.Create)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.POST("/:id/students", h.AddStudent)
	g.GET("/:id/students", h.GetStudents)
	g.DELETE("/:id/students/:student_id", h.RemoveStudent)

	// List Routes
	rg.GET("/admin/classes", h.ListClasses)
	rg.GET("/guru/classes", h.ListByTeacher)
	rg.GET("/student/classes", h.ListByStudent)
}
