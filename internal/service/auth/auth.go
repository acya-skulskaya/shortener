package auth

import (
	"fmt"
	"time"

	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	SecretKey            = "x35k9f"
	AuthCookieName       = "auth_shortener"
	AuthMetadataKeyName  = "auth-shortener"
	AuthContextKeyUserID = "userID"
)

type AuthContextKey string

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString() (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		UserID: uuid.New().String(),
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUserID(tokenString string) (userID string, err error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return "", errorsInternal.ErrTokenIsNotValid
	}

	if !token.Valid {
		return "", errorsInternal.ErrTokenIsNotValid
	}

	return claims.UserID, nil
}
