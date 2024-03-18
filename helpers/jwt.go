package helpers

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = os.Getenv("JWT_SECRET")

func GenerateToken(id int, username string) (string, error) {
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      time.Now().Add(time.Minute * 10), // expire in 10 minutes
		// "exp":   time.Now().Add(time.Minute * 2), // expire in 2 minutes for testing in weekend
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := parseToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(ctx *gin.Context) (interface{}, error) {
	headerToken := ctx.Request.Header.Get("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

	if !bearer {
		return nil, errors.New("sign in to proceed")
	}

	stringToken := strings.Split(headerToken, " ")[1]
	token, _ := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("sign in to proceed")
		}
		return []byte(secretKey), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("sign in to proceed")
	}

	expClaim, exists := claims["exp"]
	if !exists {
		return nil, errors.New("expire claim is missing")
	}

	expStr, ok := expClaim.(string)
	if !ok {
		return nil, errors.New("expire claim is not a valid type")
	}

	expTime, err := time.Parse(time.RFC3339, expStr)
	if err != nil {
		return nil, errors.New("error parsing expiration time")
	}

	if time.Now().After(expTime) {
		return nil, errors.New("token is expired")
	}

	return token.Claims.(jwt.MapClaims), nil
}
