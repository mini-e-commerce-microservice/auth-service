package handler

import (
	"errors"
	"github.com/mini-e-commerce-microservice/auth-service/generated/api"
	"github.com/mini-e-commerce-microservice/auth-service/internal/services/auth"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"net/http"
	"time"
)

func (h *handler) V1LoginPost(w http.ResponseWriter, r *http.Request) {
	req := api.V1LoginPostRequestBody{}
	if !h.httpOtel.BindBodyRequest(w, r, &req) {
		return
	}

	loginOutput, err := h.authService.Login(r.Context(), auth.LoginInput{
		Email:      req.Email,
		Password:   req.Password,
		RememberMe: req.RememberMe,
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidEmail) || errors.Is(err, auth.ErrInvalidPassword) {
			h.httpOtel.Err(w, r, http.StatusBadRequest, err, "invalid email or password")
		} else {
			h.httpOtel.Err(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     primitive.CookieRefreshTokenKey,
		Value:    loginOutput.RefreshToken.Token,
		Expires:  loginOutput.RefreshToken.ExpiredAt,
		MaxAge:   int(loginOutput.RefreshToken.ExpiredAt.Sub(time.Now().UTC()).Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	resp := api.V1LoginPostResponseBody{
		AccessToken: api.V1AuthTokenResponse{
			ExpiredAt: loginOutput.AccessToken.ExpiredAt,
			Token:     loginOutput.AccessToken.Token,
		},
		Email:           loginOutput.User.Email,
		Id:              loginOutput.User.ID,
		IsEmailVerified: loginOutput.User.IsEmailVerified,
	}

	h.httpOtel.WriteJson(w, r, http.StatusOK, resp)
}
