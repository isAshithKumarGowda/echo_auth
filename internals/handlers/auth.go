package handlers

import (
	"net/http"
	"time"

	"github.com/isAshithKumarGowda/Echo_Auth/internals/models"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	AuthRepo models.AuthModelInterface
}

func NewAuthHandler(AuthRepo models.AuthModelInterface) *AuthHandler {
	return &AuthHandler{
		AuthRepo: AuthRepo,
	}
}

func (ah *AuthHandler) RegisterHandler(e echo.Context) error {
	token, message, err := ah.AuthRepo.Register(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, echo.Map{
		"message": message,
		"token":   token,
	})
}

func (ah *AuthHandler) LoginHandler(e echo.Context) error {
	access_token, refresh_token, login_time, err := ah.AuthRepo.Login(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	accessCookie := new(http.Cookie)
	accessCookie.Name = "access_token"
	accessCookie.Value = access_token
	accessCookie.HttpOnly = true
	accessCookie.Secure = false
	accessCookie.Expires = time.Now().Add(1 * time.Hour)
	e.SetCookie(accessCookie)

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = refresh_token
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = false
	refreshCookie.Expires = time.Now().Add(5 * 24 * time.Hour)
	e.SetCookie(refreshCookie)

	return e.JSON(http.StatusOK, echo.Map{
		"message":    "Login Successful",
		"login_time": login_time,
	})
}

func (ah *AuthHandler) VeirfyEmail(e echo.Context) error {
	token, message, err := ah.AuthRepo.VeirfyEmail(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}
	return e.JSON(http.StatusOK, echo.Map{
		"message": message,
		"token":   token,
	})
}

func (ah *AuthHandler) Logout(e echo.Context) error {
	err := ah.AuthRepo.Logout(e)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	accessCookie := new(http.Cookie)
	accessCookie.Name = "access_token"
	accessCookie.Value = ""
	accessCookie.HttpOnly = true
	accessCookie.Secure = false
	accessCookie.Expires = time.Unix(0, 0)
	e.SetCookie(accessCookie)

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = ""
	refreshCookie.HttpOnly = true
	refreshCookie.Secure = false
	refreshCookie.Expires = time.Unix(0, 0)
	e.SetCookie(refreshCookie)

	return e.JSON(http.StatusOK, echo.Map{
		"message": "Logout Successful",
	})
}

func (ah *AuthHandler) GetLoginHistory(e echo.Context) error {
	users, err := ah.AuthRepo.GetLoginHistory(e)
	if err.Error() == "error bad request" {
		return e.JSON(http.StatusBadRequest, echo.Map{
			"error": "Page number exceeds total records",
		})
	} else if err != nil {
		return e.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}
	return e.JSON(http.StatusOK, echo.Map{
		"users": users,
	})
}
