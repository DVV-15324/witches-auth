package handler

import (
	"errors"
	dtoAuth "example/internal/auth/dto/request"
	mapping "example/internal/auth/mapping"
	"time"

	w_easyjson "github.com/DVV-15324/witches/pkg/core/easyjson"

	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/gin-gonic/gin"
)

func (h *AuthHandle) HandleLogin() func(c *gin.Context) {
	return func(c *gin.Context) {
		var data dtoAuth.Login

		if err := w_easyjson.BindJSON(c, &data); err != nil {
			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, w_resp.NewAppError(400, errors.New("invalid request body"), time.Now()))
			return
		}

		// Then validate
		err_valid := data.Validate()
		if err_valid != nil {
			resq := w_resp.NewAppError(401, err_valid, time.Now())
			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, resq)
			return
		}

		domain, jwk := mapping.FromLoginToDomainAuth(&data)

		ipAddress, userAgent, locale, timeZone := h.Core.Helper.GetInfo(c)

		claims, resq := h.UsecaseAuth.Login(c.Request.Context(), domain, jwk, ipAddress, userAgent, locale, timeZone)
		if resq != nil {
			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, resq)
			return
		}
		w_resp.WriteSuccess(c, claims)
	}
}
