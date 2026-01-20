package usecase

import (
	"context"
	"mutabaahapi/internal/institutions/repository"
	"mutabaahapi/internal/platform/jwt"
	userRepo "mutabaahapi/internal/users/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	InstitutionCode string
	Username        string
	Password        string
}

type LoginOutput struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	User        struct {
		ID                 string `json:"id"`
		InstitutionID      string `json:"institution_id"`
		Username           string `json:"username"`
		Name               string `json:"name"`
		Role               string `json:"role"`
		MustChangePassword bool   `json:"must_change_password"`
	} `json:"user"`
}

type Login struct {
	UserRepo        userRepo.UserRepository
	InstitutionRepo repository.InstitutionRepository
	JWTSecret       string
	TokenTTL        time.Duration
}

func NewLogin(u userRepo.UserRepository, i repository.InstitutionRepository, secret string, ttl time.Duration) *Login {
	return &Login{
		UserRepo:        u,
		InstitutionRepo: i,
		JWTSecret:       secret,
		TokenTTL:        ttl,
	}
}

func (uc *Login) Execute(ctx context.Context, in LoginInput) (*LoginOutput, error) {
	institution, err := uc.InstitutionRepo.GetByCode(ctx, in.InstitutionCode)
	if err != nil || institution == nil {
		return nil, ErrInvalidCredentials
	}
	if institution == nil {
		return nil, ErrInstitutionNotFound
	}

	user, err := uc.UserRepo.FindByIdentifier(ctx, institution.ID, in.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		return nil, ErrInvalidCredentials
	}

	claims := jwt.Claims{
		UserID: user.ID,
		Role:   user.Role,
	}
	token, exp, err := jwt.GenerateToken(claims, uc.JWTSecret, uc.TokenTTL)
	if err != nil {
		return nil, err
	}

	res := &LoginOutput{AccessToken: token, ExpiresAt: exp}
	res.User.ID = user.ID
	res.User.InstitutionID = *user.InstitutionID
	res.User.Username = user.Username
	res.User.Name = user.Name
	res.User.Role = user.Role
	res.User.MustChangePassword = user.MustChangePassword
	return res, nil
}
