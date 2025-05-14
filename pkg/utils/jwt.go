package utils

import (
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GenerateToken(userID string, userName string, userType string, expire int64) (string, error) {
	err := godotenv.Load("/home/ashith/Edwins/trial/echo_auth/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"name":      userName,
		"user_type": userType,
		"exp":       expire,
	})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
