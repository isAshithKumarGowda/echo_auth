package models

import (
	"time"

	"github.com/labstack/echo/v4"
)

type AuthModel struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthVerifyModel struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

type AuthModelInterface interface {
	Register(echo.Context) (string, string, error)
	Login(echo.Context) (string, string, time.Time, error)
	VeirfyEmail(echo.Context) (string, string, error)
	Logout(echo.Context) error
	GetLoginHistory(echo.Context) ([]AuthLoginHistoryModel, error)
}

type AuthLoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLogoutModel struct {
	Email string    `json:"email"`
	Login time.Time `json:"login_time"`
}

type AuthLoginHistoryModel struct {
	Email  string     `json:"email"`
	Login  time.Time  `json:"login_time"`
	Logout *time.Time `json:"logout_time,omitempty"`
}
