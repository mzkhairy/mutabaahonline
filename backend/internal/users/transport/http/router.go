package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	auth := rg.Group("/auth")
	auth.POST("/login", h.Login)
	auth.POST("/register-institution", h.RegisterInstitution)
}
