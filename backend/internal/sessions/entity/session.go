package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Helper untuk kolom JSON Object (Map) - Tetap ada untuk ClassActivityData
type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(b), nil
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
		return errors.New("type assertion failed")
	}
	return json.Unmarshal(bytes, &m)
}

// Helper untuk kolom JSON Array (Slice) - Untuk TemplateStructure
type JSONArray []interface{}

func (a JSONArray) Value() (driver.Value, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (a *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*a = make(JSONArray, 0)
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion failed")
	}
	return json.Unmarshal(bytes, &a)
}

type Session struct {
	ID         string `db:"id" json:"id"`
	ClassID    string `db:"class_id" json:"class_id"`
	TemplateID string `db:"template_id" json:"template_id"`

	Date                  string `db:"date" json:"date"`
	Name                  string `db:"name" json:"name"`
	IsStudentInputAllowed bool   `db:"is_student_input_allowed" json:"is_student_input_allowed"`

	ClassActivityData JSONMap `db:"class_activity_data" json:"class_activity_data"`
	TeacherID         *string `db:"teacher_id" json:"teacher_id,omitempty"`

	TemplateName      string    `db:"template_name" json:"template_name,omitempty"`
	TemplateStructure JSONArray `db:"template_structure" json:"template_structure,omitempty"`

	// [BARU] Field tambahan untuk Dashboard Murid
	ClassName string `db:"class_name" json:"class_name,omitempty"`
	HasEntry  bool   `db:"has_entry" json:"has_entry"` // True jika murid sudah mengisi

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
