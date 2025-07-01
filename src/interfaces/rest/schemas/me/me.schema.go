package meschema

import "time"

type UserAccount struct {
	ID           string    `json:"id"`
	UserIdentity string    `json:"user_identity"`
	Email        string    `json:"email"`
	Fullname     string    `json:"fullname"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserProfile struct {
	FirstName   string     `json:"first_name"`
	MiddleName  string     `json:"middle_name"`
	LastName    string     `json:"last_name"`
	Address     string     `json:"address"`
	Phone       string     `json:"phone"`
	Province    string     `json:"province"`
	City        string     `json:"city"`
	District    string     `json:"district"`
	SubDistrict string     `json:"sub_district"`
	Avatar      string     `json:"avatar"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type UserRole struct {
	RoleID    string    `json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
}

type UserSummary struct {
	TotalForm       int64 `json:"total_form"`
	TotalSubmit     int64 `json:"total_submit"`
	TotalSendSubmit int64 `json:"total_send_submit"`
}

type MeResponse struct {
	UserAccount UserAccount  `json:"user_account"`
	UserProfile *UserProfile `json:"user_profile"`
	UserRoles   []UserRole   `json:"user_roles"`
	UserSummary UserSummary  `json:"user_summary"`
}

type UserProfilePayload struct {
	FirstName   string `json:"first_name" validate:"required"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	SubDistrict string `json:"sub_district"`
	Avatar      string `json:"avatar"`
}

type ChangePasswordPayload struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
