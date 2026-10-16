package handler

import (
	"errors"
	"strings"
	"time"

	w_easyjson "github.com/DVV-15324/witches/pkg/core/easyjson"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/gin-gonic/gin"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandle) HandleLogout() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req LogoutRequest
		if err := w_easyjson.BindJSON(c, &req); err != nil {
			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, w_resp.NewAppError(400, errors.New("invalid request body"), time.Now()))
			return
		}
		authHeader := c.GetHeader("Authorization")
		var accessToken string
		if len(authHeader) > 7 && strings.ToLower(authHeader[:7]) == "bearer " {
			accessToken = authHeader[7:]
		}

		if accessToken == "" {
			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, w_resp.NewAppError(401, errors.New("access token is required"), time.Now()))
			return
		}
		//_, _, _, _ := h.Core.Helper.GetInfo(c)
		if err := h.UsecaseAuth.Logout(c.Request.Context(), accessToken, req.RefreshToken); err != nil {
			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, err)
			return
		}
		w_resp.WriteSuccess(c, gin.H{"message": "logout success"})
	}
}
