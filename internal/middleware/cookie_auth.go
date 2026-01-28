package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	SecretKey            = "x35k9f"
	AuthCookieName       = "auth_shortener"
	AuthContextKeyUserID = "userID"
)

type AuthContextKey string

func CookieAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieValue := ""
		cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				logger.Log.Debug("auth cookie not found")
			default:
				logger.Log.Debug("could not get cookie", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			cookieValue = cookie.Value
		}

		if cookieValue == "" {
			logger.Log.Info("auth cookie is empty, will create a new one")
			token, err := buildJWTString()
			if err != nil {
				logger.Log.Debug("could not create token string", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			authCookie := http.Cookie{
				Name:     AuthCookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				Expires:  time.Now().Add(time.Hour * 24 * 365),
			}

			http.SetCookie(w, &authCookie)
			cookieValue = token
		}

		userID, err := getUserID(cookieValue)
		if err != nil {
			if errors.Is(err, errorsInternal.ErrTokenIsNotValid) {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				logger.Log.Debug("error getting user id from auth token", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if userID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		logger.Log.Info("got user id", zap.String("userID", userID))
		ctx := context.WithValue(r.Context(), AuthContextKey(AuthContextKeyUserID), userID)
		r = r.WithContext(ctx)

		// передаём управление хендлеру
		next.ServeHTTP(w, r)
	})
}

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

// buildJWTString создаёт токен и возвращает его в виде строки.
func buildJWTString() (string, error) {
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

func getUserID(tokenString string) (userID string, err error) {
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
