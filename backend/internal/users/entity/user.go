package entity

import "time"

const (
	RoleAdmin = "ADMIN"
	RoleGuru  = "GURU"
	RoleMurid = "MURID"
)

type User struct {
	ID                 string    `db:"id" json:"id"`
	InstitutionID      *string   `db:"institution_id" json:"institution_id,omitempty"`
	Role               string    `db:"role" json:"role"`
	Username           string    `db:"username" json:"username"`
	Name               string    `db:"name" json:"name"`
	Email              string    `db:"email" json:"email"`
	PasswordHash       string    `db:"password_hash" json:"-"`
	SerialNumber       *string   `db:"serial_number" json:"serial_number"`
	MustChangePassword bool      `db:"must_change_password" json:"must_change_password"`
	Status             string    `db:"status" json:"status"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}
