package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	institutionsEntity "mutabaahapi/internal/institutions/entity"
	"mutabaahapi/internal/users/entity"
	"mutabaahapi/internal/users/usecase"
)

// ====================================================================
// 1. MOCK OBJECTS (Sesuai Interface Asli Anda)
// ====================================================================

// Mock User Repo (LENGKAP sesuai repository.go)
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, u *entity.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepo) FindByUsername(ctx context.Context, iID string, u string) (*entity.User, error) {
	args := m.Called(ctx, iID, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}
func (m *MockUserRepo) FindByIdentifier(ctx context.Context, iID string, id string) (*entity.User, error) {
	args := m.Called(ctx, iID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}
func (m *MockUserRepo) FindBySerialNumber(ctx context.Context, iID string, sn string) (*entity.User, error) {
	args := m.Called(ctx, iID, sn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}
func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}
func (m *MockUserRepo) ListByInstitution(ctx context.Context, iID, role, search string) ([]entity.User, error) {
	args := m.Called(ctx, iID, role, search)
	return args.Get(0).([]entity.User), args.Error(1)
}
func (m *MockUserRepo) UpdateStatus(ctx context.Context, uid, status string) error {
	return m.Called(ctx, uid, status).Error(0)
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
func (m *MockUserRepo) UpdatePassword(ctx context.Context, uid, hash string) error {
	return m.Called(ctx, uid, hash).Error(0)
}
func (m *MockUserRepo) Update(ctx context.Context, u *entity.User) error {
	return m.Called(ctx, u).Error(0)
}

// Mock Institution Repo
type MockInstitutionRepo struct {
	mock.Mock
}

func (m *MockInstitutionRepo) Create(ctx context.Context, i *institutionsEntity.Institution) error {
	return m.Called(ctx, i).Error(0)
}
func (m *MockInstitutionRepo) GetByCode(ctx context.Context, code string) (*institutionsEntity.Institution, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*institutionsEntity.Institution), args.Error(1)
}
func (m *MockInstitutionRepo) GetByID(ctx context.Context, id string) (*institutionsEntity.Institution, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*institutionsEntity.Institution), args.Error(1)
}
func (m *MockInstitutionRepo) Update(ctx context.Context, i *institutionsEntity.Institution) error {
	return m.Called(ctx, i).Error(0)
}

// ====================================================================
// 2. TEST CASES
// ====================================================================

// --- TEST REGISTER INSTITUTION (Fitur 1) ---
func TestRegisterInstitution(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockInstRepo := new(MockInstitutionRepo)

	// Gunakan NewRegister sesuai kode asli
	uCase := usecase.NewRegister(mockUserRepo, mockInstRepo)

	t.Run("Success Register", func(t *testing.T) {
		input := usecase.RegisterInstitutionInput{
			InstitutionName: "Ponpes Darul Quran",
			InstitutionCode: "pdq001",
			Name:            "Haji Fulan",
			Username:        "admin_pdq",
			Email:           "admin@test.com",
			Password:        "password123",
		}

		// 1. Cek Kode Institusi (Harus belum ada -> return nil, nil)
		mockInstRepo.On("GetByCode", mock.Anything, "pdq001").Return(nil, nil).Once()

		// 2. Create Institution
		mockInstRepo.On("Create", mock.Anything, mock.MatchedBy(func(i *institutionsEntity.Institution) bool {
			return i.Code == "pdq001"
		})).Return(nil).Once()

		// 3. Create Admin User
		mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
			return u.Username == "admin_pdq" && u.Role == "ADMIN"
		})).Return(nil).Once()

		res, err := uCase.RegisterInstitution(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "admin_pdq", res.Username)
	})

	t.Run("Fail - Code Already Exists", func(t *testing.T) {
		input := usecase.RegisterInstitutionInput{InstitutionCode: "exist01"}
		existingInst := &institutionsEntity.Institution{ID: "inst-1", Code: "exist01"}

		mockInstRepo.On("GetByCode", mock.Anything, "exist01").Return(existingInst, nil).Once()

		res, err := uCase.RegisterInstitution(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
		assert.Nil(t, res)
	})
}

// --- TEST MANAGE USERS (Fitur 3) ---
func TestCreateUserManual(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	// Gunakan NewManageUsers
	uCase := usecase.NewManageUsers(mockUserRepo)

	ctx := context.Background()
	adminID := "admin-1"
	instID := "inst-1"

	adminUser := &entity.User{
		ID:            adminID,
		Role:          "ADMIN",
		InstitutionID: &instID,
	}

	t.Run("Success Create Guru", func(t *testing.T) {
		input := usecase.CreateUserManualInput{
			Name:     "Ustadz Budi",
			Username: "guru_budi",
			Role:     "GURU",
			Password: "123",
		}

		// 1. Cek Admin
		mockUserRepo.On("GetByID", ctx, adminID).Return(adminUser, nil).Once()

		// 2. Cek Username Unik
		mockUserRepo.On("FindByUsername", ctx, instID, "guru_budi").Return(nil, nil).Once()

		// 3. Create User
		mockUserRepo.On("Create", ctx, mock.MatchedBy(func(u *entity.User) bool {
			return u.Name == "Ustadz Budi" && *u.InstitutionID == instID && u.Status == "ACTIVE"
		})).Return(nil).Once()

		err := uCase.CreateUser(ctx, adminID, input)
		assert.NoError(t, err)
	})
}
