package entity

import "time"

type Class struct {
	ID            string `json:"id" db:"id"`
	InstitutionID string `json:"institution_id" db:"institution_id"`
	// Pointer wajib untuk data yang bisa NULL
	AcademicYearID *string `json:"academic_year_id" db:"academic_year_id"`
	TeacherID      *string `json:"teacher_id" db:"teacher_id"`

	Name      string    `json:"name" db:"name"`
	Level     string    `json:"level" db:"level"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	TeacherName      *string `json:"teacher_name,omitempty" db:"teacher_name"`
	AcademicYearName *string `json:"academic_year_name,omitempty" db:"academic_year_name"`
	StudentCount     int     `json:"student_count" db:"student_count"`
}

type StudentInClass struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	SerialNumber string `json:"serial_number" db:"serial_number"`
	Status       string `json:"status" db:"status"`
}
