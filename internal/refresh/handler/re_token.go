package handler

// import (
// 	"errors"
// 	dtoRefresh "github.com/DVV-15324/witches-auth/internal/refresh/dto/request"
// 	"time"

// 	w_resp "github.com/DVV-15324/witches/pkg/core/response"
// 	"github.com/gin-gonic/gin"
// )

// func (h *RefreshHandle) HandleRefreshToken() func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		var req dtoRefresh.RefreshTokenRequest

// 		// Bind request
// 		if err := c.ShouldBindJSON(&req); err != nil {
// 			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, w_resp.NewAppError(400, errors.New("invalid request body"), time.Now()))
// 			return
// 		}

// 		//	isBlacklisted := h.Core.Blacklist.IsTokenBlacklisted(c.Request.Context(), req.RefreshToken)
// 		//if isBlacklisted {
// 			resp := w_resp.NewAppError(401, errors.New("Error: token has been revoked"), time.Now())
// 			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, resp)
// 			c.Abort()
// 			return
// 		}
// 		dpopProof := c.GetHeader("DPoP")
// 		tokenResp, err := h.RefreshUsecase.Refresh(c.Request.Context(), req.RefreshToken, dpopProof)
// 		if err != nil {
// 			w_resp.WriteErrorWithLog(c, h.Core.Logger, h.Core.Config, err)
// 			return
// 		}

// 		w_resp.WriteSuccess(c, tokenResp)
// 	}
// }
