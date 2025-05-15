package authschema

type RegisterPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Fullname string `json:"fullname" validate:"required"`
}
