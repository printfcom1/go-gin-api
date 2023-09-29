package service

import "github.com/golang-jwt/jwt/v5"

type AuthInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type JwtCustomClaims struct {
	Id       string `json:"_id"`
	UserName string `json:"username"`
	Admin    bool   `json:"admin"`
	jwt.RegisteredClaims
}

type RegisterUser struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Email           string `json:"email"`
}

type UserService interface {
	Login(AuthInput) (*string, error)
	RegisterUser(RegisterUser) (*string, error)
}
