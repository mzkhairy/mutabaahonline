package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mutabaahapi/internal/mutabaah/entity"
	"mutabaahapi/internal/mutabaah/usecase"

	// Alias userEntity agar tidak bentrok
	userEntity "mutabaahapi/internal/users/entity"
)

// --- MOCKS ---

type MockMutabaahRepo struct{ mock.Mock }

func (m *MockMutabaahRepo) Upsert(ctx context.Context, e *entity.Entry) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockMutabaahRepo) FindBySessionID(ctx context.Context, id string) ([]entity.Entry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Entry), args.Error(1)
}
func (m *MockMutabaahRepo) FindByStudentID(ctx context.Context, id string) ([]entity.Entry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Entry), args.Error(1)
}
func (m *MockMutabaahRepo) GetSummaryByStudentID(ctx context.Context, id string) (map[string]int, error) {
	return nil, nil
}

// Reuse MockUserRepo dari users_test logic (simplified version here)
type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*userEntity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

// Stub method lain user repo biar ga error interface
func (m *MockUserRepo) Create(ctx context.Context, u *userEntity.User) error { return nil }
func (m *MockUserRepo) FindByUsername(ctx context.Context, iID string, u string) (*userEntity.User, error) {
	return nil, nil
}
func (m *MockUserRepo) FindByIdentifier(ctx context.Context, iID string, id string) (*userEntity.User, error) {
	return nil, nil
}
func (m *MockUserRepo) FindBySerialNumber(ctx context.Context, iID string, sn string) (*userEntity.User, error) {
	return nil, nil
}
func (m *MockUserRepo) ListByInstitution(ctx context.Context, iID, r, s string) ([]userEntity.User, error) {
	return nil, nil
}
func (m *MockUserRepo) UpdateStatus(ctx context.Context, uid, s string) error { return nil }
func (m *MockUserRepo) CountByInstitution(ctx context.Context, iID string) (int, error) {
	return 0, nil
}
func (m *MockUserRepo) GetStats(ctx context.Context, iID string) (map[string]int, error) {
	return nil, nil
}
func (m *MockUserRepo) UpdatePassword(ctx context.Context, uid, h string) error { return nil }
func (m *MockUserRepo) Update(ctx context.Context, u *userEntity.User) error    { return nil }

func TestListMutabaah(t *testing.T) {
	mockMutabaah := new(MockMutabaahRepo)
	mockUser := new(MockUserRepo)

	// Constructor butuh 2 repo
	uCase := usecase.NewListMutabaah(mockMutabaah, mockUser)

	ctx := context.Background()

	t.Run("Teacher View List", func(t *testing.T) {
		// Data Mock
		teacher := &userEntity.User{ID: "guru-1", Role: "GURU"}
		entries := []entity.Entry{
			{ID: "e1", StudentID: "std-1", StudentActivityData: map[string]interface{}{"nilai": "A"}},
			{ID: "e2", StudentID: "std-1", StudentActivityData: map[string]interface{}{"nilai": "B"}},
		}

		// 1. Cek User
		mockUser.On("GetByID", ctx, "guru-1").Return(teacher, nil).Once()

		// 2. Repo Call (FindByStudentID karena filter by student)
		mockMutabaah.On("FindByStudentID", ctx, "std-1").Return(entries, nil).Once()

		filter := usecase.ListFilter{
			UserID:    "guru-1",
			StudentID: "std-1",
		}

		res, err := uCase.Execute(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})
}
