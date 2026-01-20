package main

import (
	"time"

	"mutabaahapi/internal/infra/config"
	"mutabaahapi/internal/infra/db"
	"mutabaahapi/internal/infra/log"
	appjwt "mutabaahapi/internal/platform/jwt"
	"mutabaahapi/internal/transport/http/router"

	// Repositories
	acadRepoPG "mutabaahapi/internal/academicyears/repository/postgres"
	classRepoPG "mutabaahapi/internal/classes/repository/postgres"
	instRepoPG "mutabaahapi/internal/institutions/repository/postgres"
	mutabaahRepoPG "mutabaahapi/internal/mutabaah/repository/postgres"
	sessionRepoPG "mutabaahapi/internal/sessions/repository/postgres"
	tplRepoPG "mutabaahapi/internal/templates/repository/postgres"
	userRepoPG "mutabaahapi/internal/users/repository/postgres"

	// UseCases
	acadUC "mutabaahapi/internal/academicyears/usecase"
	classUC "mutabaahapi/internal/classes/usecase"
	mutabaahUC "mutabaahapi/internal/mutabaah/usecase"
	sessionUC "mutabaahapi/internal/sessions/usecase"
	tplUC "mutabaahapi/internal/templates/usecase"
	userUC "mutabaahapi/internal/users/usecase"

	// Handlers
	acadHttp "mutabaahapi/internal/academicyears/transport/http"
	classHttp "mutabaahapi/internal/classes/transport/http"
	mutabaahHttp "mutabaahapi/internal/mutabaah/transport/http"
	sessHttp "mutabaahapi/internal/sessions/transport/http"
	tplHttp "mutabaahapi/internal/templates/transport/http"
	userHttp "mutabaahapi/internal/users/transport/http"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	// 1. Config & Logger
	cfg := config.Load()
	logger := log.New(cfg)

	// 2. Database
	dbConn, err := db.New(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to DB")
	}
	defer dbConn.Close()

	// 3. MIGRASI
	logger.Info().Msg("Checking database migrations...")
	if err := dbConn.Migrate("migrations"); err != nil {
		// Jika migrasi gagal, matikan server agar kita sadar ada masalah
		logger.Fatal().Err(err).Msg("Failed to run database migrations")
	}
	logger.Info().Msg("Database migrations up to date")

	// 4. Tokenizer
	tokenizer := appjwt.New(appjwt.Config{
		Secret: cfg.JWTSecret,
		TTL:    time.Duration(cfg.JWTTTLMinutes) * time.Minute,
	})

	// --- REPOSITORIES ---
	userRepo := userRepoPG.NewUserRepository(dbConn.DB)
	instRepo := instRepoPG.NewInstitutionRepository(dbConn.DB)
	acadRepo := acadRepoPG.NewRepository(dbConn.DB)
	classRepo := classRepoPG.NewRepository(dbConn.DB)
	sessionRepo := sessionRepoPG.NewSessionRepository(dbConn.DB)
	mutabaahRepo := mutabaahRepoPG.NewMutabaahRepository(dbConn.DB)
	tplRepo := tplRepoPG.NewRepository(dbConn.DB)

	// --- USECASES ---
	registerUC := userUC.NewRegister(userRepo, instRepo)
	loginUC := userUC.NewLogin(userRepo, instRepo, cfg.JWTSecret, time.Duration(cfg.JWTTTLMinutes)*time.Minute)
	manageUserUC := userUC.NewManageUsers(userRepo)
	importUserUC := userUC.NewImportUsers(userRepo)
	setupUserUC := userUC.NewSetupAccount(userRepo)

	ayUC := acadUC.NewUseCase(acadRepo, userRepo)
	classUC := classUC.NewUseCase(classRepo, userRepo)
	sessionCreateUC := sessionUC.NewCreateSession(sessionRepo)
	sessionListUC := sessionUC.NewListSession(sessionRepo)
	sessionGetUC := sessionUC.NewGetSession(sessionRepo)
	sessionUpdateUC := sessionUC.NewUpdateSession(sessionRepo)
	sessionListStudentUC := sessionUC.NewListStudentSession(sessionRepo)

	mutabaahInputUC := mutabaahUC.NewInputMutabaah(mutabaahRepo, sessionRepo, userRepo, classRepo)
	mutabaahListUC := mutabaahUC.NewListMutabaah(mutabaahRepo, userRepo)
	mutabaahBulkUC := mutabaahUC.NewBulkAttendance(mutabaahRepo, userRepo)
	mutabaahSummaryUC := mutabaahUC.NewMutabaahSummary(mutabaahRepo, userRepo)

	tplUC := tplUC.NewUseCase(tplRepo, userRepo)

	// --- HANDLERS ---
	// Perhatikan urutan parameter sesuai constructor Handler
	userHandler := userHttp.NewHandler(registerUC, loginUC, manageUserUC, importUserUC, setupUserUC, tokenizer, logger)
	ayHandler := acadHttp.NewHandler(ayUC, logger)
	clsHandler := classHttp.NewHandler(classUC, logger)
	sessionHandler := sessHttp.NewHandler(sessionCreateUC, sessionListUC, sessionGetUC, sessionUpdateUC, sessionListStudentUC)
	mutabaahHandler := mutabaahHttp.NewHandler(mutabaahInputUC, mutabaahListUC, mutabaahBulkUC, mutabaahSummaryUC)
	templateHandler := tplHttp.NewHandler(tplUC)

	// --- ROUTER ---
	// Masukkan semua dependency ke Router
	r := router.New(
		logger,
		tokenizer,
		userHandler,
		sessionHandler,
		mutabaahHandler,
		ayHandler,
		clsHandler,
		templateHandler,
	)

	// Start Server
	logger.Info().Msgf("Server starting on port %s...", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
