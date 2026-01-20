package router

import (
	"mutabaahapi/internal/platform/apperror"
	appjwt "mutabaahapi/internal/platform/jwt"
	"mutabaahapi/internal/transport/http/middleware"

	// Import Handler Module
	academicyearshttp "mutabaahapi/internal/academicyears/transport/http"
	classeshttp "mutabaahapi/internal/classes/transport/http"
	mutabaahhttp "mutabaahapi/internal/mutabaah/transport/http"
	sessionshttp "mutabaahapi/internal/sessions/transport/http"
	templateshttp "mutabaahapi/internal/templates/transport/http"
	usershttp "mutabaahapi/internal/users/transport/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Fungsi Custom Middleware CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Daftar Domain Frontend yang diizinkan
		allowedOrigins := map[string]bool{
			"http://localhost:5173":    true, // Frontend Lokal
			"http://localhost":         true,
			"http://103.127.133.92":    true,
			"http://103.127.133.92:80": true,
		}

		// Jika origin ada di daftar allow, set header
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// Opsional: Izinkan semua untuk tahap development awal (HATI-HATI untuk production)
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// New menerima SEMUA dependency dari main.go
func New(
	l zerolog.Logger,
	tokenizer *appjwt.Tokenizer,
	userHandler *usershttp.Handler,
	sessionHandler *sessionshttp.Handler,
	mutabaahHandler *mutabaahhttp.Handler,
	ayHandler *academicyearshttp.Handler,
	clsHandler *classeshttp.Handler,
	templateHandler *templateshttp.Handler,
) *gin.Engine {

	r := gin.New()

	// 0. Pasang CORS Middleware
	r.Use(CORSMiddleware())

	// 1. Middleware Global
	r.Use(middleware.Recover(l))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(l))

	// 2. Health Check
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// 3. API V1
	api := r.Group("/api/v1")
	{
		// A. Public Routes
		usershttp.RegisterRoutes(api, userHandler)

		// B. Protected Routes
		protected := api.Group("/", middleware.AuthRequired(l, tokenizer))
		{
			// User Management
			protected.GET("/users", userHandler.ListUsers)
			protected.GET("/users/:id", userHandler.GetUser)
			protected.PATCH("/users/:id/status", userHandler.UpdateUserStatus)
			protected.PATCH("/users/:id/password", userHandler.ResetPassword)
			protected.PUT("/users/:id", userHandler.UpdateUser)
			protected.POST("/users/import", userHandler.ImportUsers)
			protected.POST("/users", userHandler.CreateUser)
			protected.POST("/auth/setup-account", userHandler.SetupAccount)

			// Test Auth
			protected.GET("/me", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				role, _ := c.Get("user_role")
				c.JSON(200, gin.H{"user_id": userID, "role": role})
			})

			// Module Routes
			sessionshttp.RegisterRoutes(protected, sessionHandler)
			mutabaahhttp.RegisterRoutes(protected, mutabaahHandler)
			academicyearshttp.RegisterRoutes(protected, ayHandler)
			classeshttp.RegisterRoutes(protected, clsHandler)
			templateshttp.RegisterRoutes(protected, templateHandler)
		}
	}

	// 4. Handle 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, apperror.ErrNotFound)
	})

	return r
}
