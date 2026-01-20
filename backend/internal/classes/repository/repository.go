package repository

import (
	"context"
	"mutabaahapi/internal/classes/entity"
)

type ClassRepository interface {
	Create(ctx context.Context, c *entity.Class) error
	GetByID(ctx context.Context, id string) (*entity.Class, error)

	// Update & Delete
	Update(ctx context.Context, c *entity.Class) error
	Delete(ctx context.Context, id string) error

	// List Functions
	ListByInstitution(ctx context.Context, institutionID string, academicYearID string) ([]entity.Class, error)
	ListByTeacher(ctx context.Context, teacherID string, academicYearID string) ([]entity.Class, error)

	// [FIX] Tambahkan ini agar Usecase tidak error "undefined"
	ListByStudent(ctx context.Context, studentID string) ([]entity.Class, error)

	// Student Management
	AddStudent(ctx context.Context, classID, studentID string) error
	RemoveStudent(ctx context.Context, classID, studentID string) error
	GetStudents(ctx context.Context, classID string) ([]entity.StudentInClass, error)
	IsStudentInClass(ctx context.Context, classID, studentID string) (bool, error)
}
