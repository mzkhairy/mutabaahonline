package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mutabaahapi/internal/sessions/entity"
	"mutabaahapi/internal/sessions/usecase"
)

// Mock Session Repo
type MockSessionRepo struct {
	mock.Mock
}

func (m *MockSessionRepo) Create(ctx context.Context, s *entity.Session) error {
	return m.Called(ctx, s).Error(0)
}
func (m *MockSessionRepo) IsClassTeacher(ctx context.Context, classID, teacherID string) (bool, error) {
	args := m.Called(ctx, classID, teacherID)
	return args.Bool(0), args.Error(1)
}

// Stub method lain (GetByID, ListByClass, dll) biar implement interface
func (m *MockSessionRepo) GetByID(ctx context.Context, id string) (*entity.Session, error) {
	return nil, nil
}
func (m *MockSessionRepo) ListByClass(ctx context.Context, id string) ([]entity.Session, error) {
	return nil, nil
}
func (m *MockSessionRepo) ListByStudent(ctx context.Context, sid string, mode bool, cid string) ([]entity.Session, error) {
	return nil, nil
}
func (m *MockSessionRepo) Update(ctx context.Context, s *entity.Session) error { return nil }

func TestCreateSession(t *testing.T) {
	mockRepo := new(MockSessionRepo)
	// Gunakan NewCreateSession sesuai kode asli
	uCase := usecase.NewCreateSession(mockRepo)

	t.Run("Success Create Session", func(t *testing.T) {
		input := usecase.CreateSessionInput{
			UserID:                "guru-1",
			ClassID:               "class-A",
			Name:                  "Setoran Pagi",
			Date:                  "2025-01-27", // String, bukan time.Time
			IsStudentInputAllowed: true,
		}

		// 1. Cek Otoritas Guru
		mockRepo.On("IsClassTeacher", mock.Anything, "class-A", "guru-1").Return(true, nil).Once()

		// 2. Create
		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(s *entity.Session) bool {
			return s.Name == "Setoran Pagi" && s.ClassID == "class-A"
		})).Return(nil).Once()

		res, err := uCase.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Fail - Not Class Teacher", func(t *testing.T) {
		input := usecase.CreateSessionInput{
			UserID:  "guru-asing",
			ClassID: "class-A",
		}

		// Return False
		mockRepo.On("IsClassTeacher", mock.Anything, "class-A", "guru-asing").Return(false, nil).Once()

		res, err := uCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Equal(t, "you are not the teacher of this class", err.Error())
		assert.Nil(t, res)
	})
}
