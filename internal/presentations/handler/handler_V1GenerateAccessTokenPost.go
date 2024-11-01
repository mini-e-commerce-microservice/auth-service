package handler

import (
	"errors"
	"github.com/mini-e-commerce-microservice/auth-service/generated/api"
	"github.com/mini-e-commerce-microservice/auth-service/internal/services/auth"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"net/http"
	"strconv"
)

func (h *handler) V1GenerateAccessTokenPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(primitive.CookieRefreshTokenKey)
	if err != nil {
		h.httpOtel.Err(w, r, http.StatusUnauthorized, err)
		return
	}

	updateUserDataStr := r.URL.Query().Get("update_user_data")
	updateUserData := false
	if updateUserDataStr != "" {
		updateUserData, err = strconv.ParseBool(updateUserDataStr)
		if err != nil {
			h.httpOtel.Err(w, r, http.StatusBadRequest, err, "invalid query param update_user_data")
			return
		}
	}
	accessTokenOutput, err := h.authService.GenerateAccessToken(r.Context(), auth.GenerateAccessTokenInput{
		RefreshToken:          cookie.Value,
		UpdateUserDataInCache: updateUserData,
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrRefreshTokenNotExistsInRedis) {
			h.httpOtel.Err(w, r, http.StatusUnauthorized, err)
		} else {
			h.httpOtel.Err(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	resp := api.V1GenerateAccessTokenPostResponseBody{
		AccessToken: api.V1AuthTokenResponse{
			ExpiredAt: accessTokenOutput.ExpiredAt,
			Token:     accessTokenOutput.Token,
		},
	}

	h.httpOtel.WriteJson(w, r, http.StatusOK, resp)
}
