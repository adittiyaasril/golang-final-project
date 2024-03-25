package helpers

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	SECRETKEY = "my-gram"
	EXPIRY    = 24 * time.Hour
)

func GenerateToken(id uint64, email string) (string, error) {
	expirationTime := time.Now().Add(EXPIRY)

	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(ctx *gin.Context) (jwt.MapClaims, error) {
	err := errors.New("please login to get the token")
	auth := ctx.Request.Header.Get("Authorization")
	bearer := strings.HasPrefix(auth, "Bearer ")

	if !bearer {
		return nil, err
	}

	tokenStr := strings.Split(auth, "Bearer ")[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRETKEY), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("token is not valid")
	}
}
