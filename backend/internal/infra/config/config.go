package config

import (
	"os"
	"strconv"
)

type Config struct {
	Env                  string
	HTTPPort             string // Akan diisi oleh env PORT
	DatabaseURL          string
	RunMigrations        bool
	MigrationsDir        string
	JWTSecret            string
	JWTTTLMinutes        int
	DBMaxOpenConns       int
	DBMaxIdleConns       int
	DBConnMaxLifetimeMin int
	DBConnMaxIdleTimeMin int
	LogLevel             string
	LogPretty            bool
}

func getenv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getbool(key string, def bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}

func Load() Config {
	return Config{
		Env: getenv("APP_ENV", "dev"),

		// [FIX] Prioritaskan "PORT" (Standar Cloud: Render/Railway)
		// Jika tidak ada PORT, gunakan "8080" (Lokal)
		HTTPPort: getenv("PORT", "8080"),

		DatabaseURL:   getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/notes?sslmode=disable"),
		RunMigrations: getbool("RUN_MIGRATIONS", true),
		MigrationsDir: getenv("MIGRATIONS_DIR", "migrations"),
		JWTSecret:     getenv("JWT_SECRET", "devsecret"),
		JWTTTLMinutes: func() int {
			v := getenv("JWT_TTL_MINUTES", "60")
			i, _ := strconv.Atoi(v)
			if i <= 0 {
				i = 60
			}
			return i
		}(),
		DBMaxOpenConns: func() int {
			v := getenv("DB_MAX_OPEN_CONNS", "20")
			i, _ := strconv.Atoi(v)
			if i <= 0 {
				i = 20
			}
			return i
		}(),
		DBMaxIdleConns: func() int {
			v := getenv("DB_MAX_IDLE_CONNS", "5")
			i, _ := strconv.Atoi(v)
			if i < 0 {
				i = 5
			}
			return i
		}(),
		DBConnMaxLifetimeMin: func() int {
			v := getenv("DB_CONN_MAX_LIFETIME_MIN", "60")
			i, _ := strconv.Atoi(v)
			if i <= 0 {
				i = 60
			}
			return i
		}(),
		DBConnMaxIdleTimeMin: func() int {
			v := getenv("DB_CONN_MAX_IDLE_TIME_MIN", "10")
			i, _ := strconv.Atoi(v)
			if i < 0 {
				i = 10
			}
			return i
		}(),
		LogLevel:  getenv("LOG_LEVEL", ""),
		LogPretty: getbool("LOG_PRETTY", false),
	}
}
