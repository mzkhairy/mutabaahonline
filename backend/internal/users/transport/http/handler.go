package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"mutabaahapi/internal/platform/apperror"
	appjwt "mutabaahapi/internal/platform/jwt"
	"mutabaahapi/internal/transport/http/contextx"
	"mutabaahapi/internal/transport/http/response"
	"mutabaahapi/internal/users/usecase"
)

type Handler struct {
	RegisterUC     *usecase.Register
	LoginUC        *usecase.Login
	ManageUserUC   *usecase.ManageUsers
	ImportUserUC   *usecase.ImportUsers // [PENTING] Field ini wajib ada
	Tokenizer      *appjwt.Tokenizer
	Logger         zerolog.Logger
	SetupAccountUC *usecase.SetupAccount
}

func NewHandler(
	registerUC *usecase.Register,
	loginUC *usecase.Login,
	mng *usecase.ManageUsers,
	imp *usecase.ImportUsers,
	setupUC *usecase.SetupAccount,
	tokenizer *appjwt.Tokenizer,
	logger zerolog.Logger,
) *Handler {
	return &Handler{
		RegisterUC:     registerUC,
		LoginUC:        loginUC,
		ManageUserUC:   mng,
		ImportUserUC:   imp,
		SetupAccountUC: setupUC,
		Tokenizer:      tokenizer,
		Logger:         logger,
	}
}

// --- Handler Functions ---

type loginRequest struct {
	InstitutionCode string `json:"institution_code" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.NewWithDetails(apperror.ErrCodeBadRequest, "invalid request body", err.Error()))
		return
	}
	input := usecase.LoginInput{
		InstitutionCode: req.InstitutionCode,
		Username:        req.Username,
		Password:        req.Password,
	}
	res, err := h.LoginUC.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, mapDomainError(err))
		return
	}
	response.Success(c, http.StatusOK, "Login successful", res)
}

type registerInstitutionRequest struct {
	InstitutionName string `json:"institution_name" binding:"required"`
	InstitutionCode string `json:"institution_code" binding:"required,min=3"`
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
}

func (h *Handler) RegisterInstitution(c *gin.Context) {
	var req registerInstitutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.NewWithDetails(apperror.ErrCodeBadRequest, "invalid request body", err.Error()))
		return
	}
	input := usecase.RegisterInstitutionInput{
		InstitutionName: req.InstitutionName,
		InstitutionCode: req.InstitutionCode,
		Username:        req.Username,
		Password:        req.Password,
		Name:            req.Name,
		Email:           req.Email,
	}
	user, err := h.RegisterUC.RegisterInstitution(c.Request.Context(), input)
	if err != nil {
		response.Error(c, mapDomainError(err))
		return
	}
	response.Success(c, http.StatusCreated, "Institution and Admin registered", user)
}

func (h *Handler) ListUsers(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}
	roleFilter := c.Query("role")
	search := c.Query("q")
	res, err := h.ManageUserUC.ExecuteList(c.Request.Context(), userID, roleFilter, search)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Users loaded", res)
}

func (h *Handler) ResetPassword(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}
	targetID := c.Param("id")
	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}
	if err := h.ManageUserUC.ResetPassword(c.Request.Context(), userID, targetID, req.NewPassword); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Password reset successfully", nil)
}

func (h *Handler) UpdateUserStatus(c *gin.Context) {
	targetUserID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}
	if err := h.ManageUserUC.ChangeStatus(c.Request.Context(), targetUserID, req.Status); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "User status updated", nil)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	adminID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}
	targetID := c.Param("id")
	var req struct {
		Name         string `json:"name" binding:"required"`
		Username     string `json:"username" binding:"required"`
		SerialNumber string `json:"serial_number"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}
	err := h.ManageUserUC.UpdateUser(c.Request.Context(), adminID, targetID, usecase.UpdateUserInput{
		Name:         req.Name,
		Username:     req.Username,
		SerialNumber: req.SerialNumber,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "User updated", nil)
}

func (h *Handler) ImportUsers(c *gin.Context) {
	adminID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, apperror.New(apperror.ErrCodeInvalidInput, "file is required (csv)"))
		return
	}
	defer file.Close()
	res, err := h.ImportUserUC.ExecuteWrapper(c.Request.Context(), adminID, file)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Import processed", res)
}

// [BARU] POST /users (Manual Create)
func (h *Handler) CreateUser(c *gin.Context) {
	adminID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	var req struct {
		Name         string `json:"name" binding:"required"`
		Username     string `json:"username" binding:"required"`
		SerialNumber string `json:"serial_number"`
		Role         string `json:"role" binding:"required,oneof=GURU MURID"`
		Password     string `json:"password"` // Optional
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	err := h.ManageUserUC.CreateUser(c.Request.Context(), adminID, usecase.CreateUserManualInput{
		Name:         req.Name,
		Username:     req.Username,
		SerialNumber: req.SerialNumber,
		Role:         req.Role,
		Password:     req.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "User created successfully", nil)
}

func (h *Handler) SetupAccount(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c) // Ambil dari token
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}

	var req struct {
		Username string `json:"username" binding:"required,min=4"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.ErrInvalidInput)
		return
	}

	err := h.SetupAccountUC.Execute(c.Request.Context(), userID, usecase.SetupInput{
		NewUsername: req.Username,
		NewPassword: req.Password,
		Email:       req.Email,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Account setup completed", nil)
}

func (h *Handler) GetUser(c *gin.Context) {
	userID, ok := contextx.UserIDFrom(c)
	if !ok {
		response.Error(c, apperror.ErrUnauthorized)
		return
	}
	targetID := c.Param("id")

	user, err := h.ManageUserUC.GetUser(c.Request.Context(), userID, targetID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "User detail retrieved", user)
}
