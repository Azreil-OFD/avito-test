package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

type Token struct {
	mySecret string
}

func NewToken(mySecret string) Token {
	return Token{mySecret: mySecret}
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (t Token) GenerateJWT(userID int) (string, error) {
	claims := Claims{
		UserID: strconv.Itoa(userID),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "avitotest",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Токен действует 24 часа
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(t.mySecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t Token) validateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(t.mySecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

func (t Token) GetUserIDFromJWT(tokenString string) (string, error) {
	claims, err := t.validateJWT(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
