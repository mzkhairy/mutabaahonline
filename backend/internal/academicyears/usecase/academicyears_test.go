package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mutabaahapi/internal/academicyears/entity"
	"mutabaahapi/internal/academicyears/usecase"
	userEntity "mutabaahapi/internal/users/entity"
)

// --- MOCKS ---

type MockAcademicRepo struct {
	mock.Mock
}

func (m *MockAcademicRepo) Create(ctx context.Context, e *entity.AcademicYear) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockAcademicRepo) Update(ctx context.Context, e *entity.AcademicYear) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockAcademicRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockAcademicRepo) GetByID(ctx context.Context, id string) (*entity.AcademicYear, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AcademicYear), args.Error(1)
}
func (m *MockAcademicRepo) List(ctx context.Context, institutionID string) ([]entity.AcademicYear, error) {
	args := m.Called(ctx, institutionID)
	return args.Get(0).([]entity.AcademicYear), args.Error(1)
}
func (m *MockAcademicRepo) ResetStatusByInstitution(ctx context.Context, institutionID string) error {
	return m.Called(ctx, institutionID).Error(0)
}

// [FIX] Mock User Repo LENGKAP (Copy-Paste dari users_test.go)
type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) Create(ctx context.Context, u *userEntity.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepo) FindByUsername(ctx context.Context, iID, u string) (*userEntity.User, error) {
	args := m.Called(ctx, iID, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}
func (m *MockUserRepo) FindByIdentifier(ctx context.Context, iID, id string) (*userEntity.User, error) {
	args := m.Called(ctx, iID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}
func (m *MockUserRepo) FindBySerialNumber(ctx context.Context, iID, sn string) (*userEntity.User, error) {
	args := m.Called(ctx, iID, sn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}
func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*userEntity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}
func (m *MockUserRepo) ListByInstitution(ctx context.Context, iID, r, s string) ([]userEntity.User, error) {
	args := m.Called(ctx, iID, r, s)
	return args.Get(0).([]userEntity.User), args.Error(1)
}
func (m *MockUserRepo) UpdateStatus(ctx context.Context, uid, s string) error {
	return m.Called(ctx, uid, s).Error(0)
}
func (m *MockUserRepo) CountByInstitution(ctx context.Context, iID string) (int, error) {
	args := m.Called(ctx, iID)
	return args.Int(0), args.Error(1)
}
func (m *MockUserRepo) GetStats(ctx context.Context, iID string) (map[string]int, error) {
	args := m.Called(ctx, iID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}
func (m *MockUserRepo) UpdatePassword(ctx context.Context, uid, h string) error {
	return m.Called(ctx, uid, h).Error(0)
}
func (m *MockUserRepo) Update(ctx context.Context, u *userEntity.User) error {
	return m.Called(ctx, u).Error(0)
}

// --- TEST CASES ---

func TestCreateAcademicYear(t *testing.T) {
	mockRepo := new(MockAcademicRepo)
	mockUserRepo := new(MockUserRepo)
	uCase := usecase.NewUseCase(mockRepo, mockUserRepo)

	ctx := context.Background()
	adminID := "admin-1"
	instID := "inst-1"
	adminUser := &userEntity.User{ID: adminID, InstitutionID: &instID}

	t.Run("Success Create Active Year", func(t *testing.T) {
		input := usecase.CreateInput{
			Name:      "2025/2026",
			StartDate: "2025-07-01",
			EndDate:   "2026-06-30",
			IsActive:  true,
		}

		mockUserRepo.On("GetByID", ctx, adminID).Return(adminUser, nil).Once()
		mockRepo.On("ResetStatusByInstitution", ctx, instID).Return(nil).Once()
		mockRepo.On("Create", ctx, mock.MatchedBy(func(e *entity.AcademicYear) bool {
			return e.Name == "2025/2026" && e.IsActive == true
		})).Return(nil).Once()

		err := uCase.Create(ctx, adminID, input)
		assert.NoError(t, err)
	})

	t.Run("Fail Invalid Date Format", func(t *testing.T) {
		input := usecase.CreateInput{
			Name:      "Invalid Date",
			StartDate: "01-07-2025",
			EndDate:   "2026-06-30",
		}

		mockUserRepo.On("GetByID", ctx, adminID).Return(adminUser, nil).Once()

		err := uCase.Create(ctx, adminID, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid start_date format")
	})

	t.Run("Fail Unauthorized", func(t *testing.T) {
		mockUserRepo.On("GetByID", ctx, "unknown").Return(nil, errors.New("not found")).Once()

		err := uCase.Create(ctx, "unknown", usecase.CreateInput{})
		assert.Error(t, err)
		assert.Equal(t, "unauthorized", err.Error())
	})
}
