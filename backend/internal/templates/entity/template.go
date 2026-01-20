package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Struktur Item di dalam JSON
type TemplateItem struct {
	Name string `json:"name"` // Contoh: "Hafalan Baru", "Materi"
	Type string `json:"type"` // "class_activity" atau "student_activity"
}

// Custom Type agar bisa masuk ke Database sebagai JSONB
type TemplateStructure []TemplateItem

// Value: Mengubah Struct Go -> JSON string untuk Database
func (ts TemplateStructure) Value() (driver.Value, error) {
	return json.Marshal(ts)
}

// Scan: Mengubah JSON Database -> Struct Go
func (ts *TemplateStructure) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &ts)
}

// Entity Utama
type Template struct {
	ID            string            `db:"id" json:"id"`
	InstitutionID string            `db:"institution_id" json:"institution_id"`
	Name          string            `db:"name" json:"name"`
	Structure     TemplateStructure `db:"structure" json:"structure"` // Tipe Custom tadi
	CreatedAt     time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time         `db:"updated_at" json:"updated_at"`
}
