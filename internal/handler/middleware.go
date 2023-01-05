package handler

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

const ctxKeyUser ctxKey = iota

type ctxKey int8

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_cookie")
		if err != nil {
			logrus.Errorf("cookie: %w", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userID, err := h.service.ParseToken(cookie.Value)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, userID)))
	})
}
