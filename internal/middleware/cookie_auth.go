package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	authService "github.com/acya-skulskaya/shortener/internal/service/auth"
	"go.uber.org/zap"
)

func CookieAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieValue := ""
		cookie, err := r.Cookie(authService.AuthCookieName)
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
			token, err := authService.BuildJWTString()
			if err != nil {
				logger.Log.Debug("could not create token string", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			authCookie := http.Cookie{
				Name:     authService.AuthCookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				Expires:  time.Now().Add(time.Hour * 24 * 365),
			}

			http.SetCookie(w, &authCookie)
			cookieValue = token
		}

		userID, err := authService.GetUserID(cookieValue)
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
		ctx := context.WithValue(r.Context(), authService.AuthContextKey(authService.AuthContextKeyUserID), userID)
		r = r.WithContext(ctx)

		// передаём управление хендлеру
		next.ServeHTTP(w, r)
	})
}
