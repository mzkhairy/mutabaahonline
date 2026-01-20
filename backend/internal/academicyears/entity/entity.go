package entity

import "time"

type AcademicYear struct {
	ID            string    `db:"id" json:"id"`
	InstitutionID string    `db:"institution_id" json:"institution_id"`
	Name          string    `db:"name" json:"name"`             // "2025/2026 Ganjil"
	StartDate     string    `db:"start_date" json:"start_date"` // YYYY-MM-DD (String biar gampang diparse di frontend)
	EndDate       string    `db:"end_date" json:"end_date"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
