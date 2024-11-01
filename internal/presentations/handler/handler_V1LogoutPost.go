package handler

import (
	"errors"
	"github.com/mini-e-commerce-microservice/auth-service/internal/services/auth"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"net/http"
)

func (h *handler) V1LogoutPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(primitive.CookieRefreshTokenKey)
	if err != nil {
		h.httpOtel.Err(w, r, http.StatusUnauthorized, err)
		return
	}

	err = h.authService.Logout(r.Context(), auth.LogoutInput{
		RefreshToken: cookie.Value,
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			h.httpOtel.Err(w, r, http.StatusUnauthorized, err)
		} else {
			h.httpOtel.Err(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
