package middleware

import (
	"context"
	"errors"
	"strings"
	"time"

	domainRefresh "github.com/DVV-15324/witches-auth/internal/shared/domain"
	utils "github.com/DVV-15324/witches-auth/internal/shared/utils"
	pkg "github.com/DVV-15324/witches-auth/pkg"
	// Thư viện dùng chung (Witches)
	pkg_dpop "github.com/DVV-15324/witches-auth/pkg/dpop"
	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"

	"github.com/gin-gonic/gin"
)

type UsecaseRefresh interface {
	GetByJKT(ctx context.Context, JKT string) (*domainRefresh.RefreshToken, *w_resp.AppError)
}

type JWTService interface {
	ParseToken(ctx context.Context, token string) (*pkg.JwtClaims, error)
}

type BlacklistService interface {
	IsTokenBlacklisted(ctx context.Context, token string) bool
}

type SessionService interface {
	ValidateSession(ctx context.Context, userID int64, jkt string, accessToken string) (*pkg.SessionCache, error)
	UpdateLastActive(ctx context.Context, userID int64, dpopProof string) error
}

type AuthDependencies struct {
	UsecaseRefresh UsecaseRefresh
	JWT            JWTService
	Blacklist      BlacklistService
	Session        SessionService
	Logger         *logger.ModelLogger
	Config         *wcmd_utils.Config
}

func AuthMiddleware(deps AuthDependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		// access token
		accessToken, err := ExtractDPoPTokenFromHeader(c.GetHeader("Authorization"))
		if err != nil {
			resp := w_resp.NewAppError(401, err, time.Now())
			w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, resp)
			c.Abort()
			return
		}

		claims, errResp := deps.JWT.ParseToken(c.Request.Context(), accessToken)
		if errResp != nil {
			resp := w_resp.NewAppError(401, errResp, time.Now())
			w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, resp)
			c.Abort()
			return
		}

		dpopProof := c.GetHeader("DPoP")
		jkt := claims.GetJKT()

		refreshToken, err_jkt := deps.UsecaseRefresh.GetByJKT(c, jkt)
		if err_jkt != nil {
			w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, err_jkt)
			c.Abort()
			return
		}

		proof_, err_verity := pkg_dpop.DPoPVerify(c, dpopProof, refreshToken.JKT, 5*time.Second, accessToken)
		if err_verity != nil {
			resp := w_resp.NewAppError(401, err_verity, time.Now())
			w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, resp)
			c.Abort()
			return
		}

		isBlacklisted := deps.Blacklist.IsTokenBlacklisted(c.Request.Context(), accessToken)
		if isBlacklisted {
			resp := w_resp.NewAppError(401, errors.New("Error: token has been revoked"), time.Now())
			w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, resp)
			c.Abort()
			return
		}

		UID := utils.DecodeFromBase58(claims.Subject)
		userID := UID.LocalID

		if userID == 0 {
			resp := w_resp.NewAppError(401, errors.New("Error: invalid user ID in token"), time.Now())
			w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, resp)
			c.Abort()
			return
		}

		reqCtx := w_utils.GetRequestContext(c.Request.Context(), deps.Config)

		_, err = deps.Session.ValidateSession(c.Request.Context(), userID, proof_.JKT, accessToken)
		if err != nil {
			c.Abort()
			resp := w_resp.NewAppError(401, err, time.Now())
			w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, resp)
			return
		}

		reqCtx.Sub = claims.Subject
		reqCtx.Tid = claims.ID
		ctx := w_utils.SaveRequestContext(c.Request.Context(), reqCtx, deps.Config)
		c.Request = c.Request.WithContext(ctx)

		// RBAC
		c.Next()

		// async response
		go func() {
			ctxBg, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			deps.Session.UpdateLastActive(ctxBg, userID, dpopProof)
		}()
	}
}

func ExtractDPoPTokenFromHeader(header string) (string, error) {
	if header == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Fields(header)

	if len(parts) != 2 || parts[0] != "DPoP" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
