package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/isAshithKumarGowda/Echo_Auth/pkg/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := godotenv.Load("/home/ashith/Edwins/trial/echo_auth/.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		jwtSecret := os.Getenv("JWT_SECRET")

		AccessCookie, err := c.Cookie("access_token")
		if err == nil {
			accessToken, err := jwt.Parse(AccessCookie.Value, func(t *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err == nil && accessToken.Valid {
				return next(c)
			}
		}

		RefreshCookie, err := c.Cookie("refresh_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Unauthorized: Login required"})
		}

		token, err := jwt.Parse(RefreshCookie.Value, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || token == nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Unauthorized: Refresh token invalid"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Invalid token claims",
			})
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Invalid user_id",
			})
		}

		userName, ok := claims["name"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Invalid user_name",
			})
		}

		userType, ok := claims["user_type"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Invalid user_name",
			})
		}

		access_token, err := utils.GenerateToken(userID, userName, userType, time.Now().Local().Add(time.Hour).Unix())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Error generating access token",
			})
		}

		accessCookie := new(http.Cookie)
		accessCookie.Name = "access_token"
		accessCookie.Value = access_token
		accessCookie.HttpOnly = true
		accessCookie.Secure = false
		accessCookie.Expires = time.Now().Add(1 * time.Hour)
		c.SetCookie(accessCookie)

		return next(c)
	}
}
