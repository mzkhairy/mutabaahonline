package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(b), nil // <--- Kirim sebagai STRING, bukan bytes
}
func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = make(JSONMap)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte or string failed")
	}

	return json.Unmarshal(bytes, &m)
}

type Entry struct {
	ID        string `db:"id" json:"id"`
	SessionID string `db:"session_id" json:"session_id"`
	StudentID string `db:"student_id" json:"student_id"`

	// JOIN Fields
	StudentName string `db:"student_name" json:"student_name,omitempty"`
	SessionName string `db:"session_name" json:"session_name,omitempty"`

	Attendance          string  `db:"attendance" json:"attendance"`
	StudentActivityData JSONMap `db:"student_activity_data" json:"student_activity_data"`
	Status              string  `db:"status" json:"status"` // LANJUT / ULANG
	Note                string  `db:"note" json:"note"`

	IsLocked      bool   `db:"is_locked" json:"is_locked"`
	LastUpdatedBy string `db:"last_updated_by" json:"last_updated_by"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
