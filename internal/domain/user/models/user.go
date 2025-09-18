package models

type User struct {
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Email    string `json:"email"    validate:"email"`
	Password string `json:"password" validate:"min=8"`
	Age      int    `json:"age"      validate:"gte=18"`
	Phone    string `json:"phone"    validate:"e164"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
