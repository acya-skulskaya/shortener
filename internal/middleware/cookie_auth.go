package middleware

import (
	"context"
	"errors"
	"fmt"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	SECRET_KEY       = "x35k9f"
	AUTH_COOKIE_NAME = "auth"
)

// Выдавать пользователю симметрично подписанную куку, содержащую уникальный идентификатор пользователя, если такой куки не существует или она не проходит проверку подлинности.

func CookieAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieValue := ""
		cookie, err := r.Cookie(AUTH_COOKIE_NAME)
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
			token, err := BuildJWTString()
			if err != nil {
				logger.Log.Debug("could not create token string", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			cookie := http.Cookie{
				Name:     AUTH_COOKIE_NAME,
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				Expires:  time.Now().Add(time.Hour * 24 * 365),
			}

			http.SetCookie(w, &cookie)
			cookieValue = token
		}

		userID, err := GetUserID(cookieValue)
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
		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

		// передаём управление хендлеру
		next.ServeHTTP(w, r)
	})
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString() (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		UserID: uuid.New().String(),
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUserID(tokenString string) (userID string, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return "", errorsInternal.ErrTokenIsNotValid
	}

	if !token.Valid {
		return "", errorsInternal.ErrTokenIsNotValid
	}

	return claims.UserID, nil
}
