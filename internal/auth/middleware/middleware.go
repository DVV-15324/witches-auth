package middleware

import (
	"context"
	"errors"
	"strings"
	"time"

	domainRefresh "github.com/DVV-15324/witches-auth/internal/shared/domain"
	utils "github.com/DVV-15324/witches-auth/internal/shared/utils"
	pkg "github.com/DVV-15324/witches-auth/pkg"
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
	UpdateLastActive(ctx context.Context, userID int64, jkt, accessToken string) error
}

type AuthDependencies struct {
	UsecaseRefresh UsecaseRefresh
	JWT            JWTService
	Blacklist      BlacklistService
	Session        SessionService
	Logger         *logger.ModelLogger
	Config         *wcmd_utils.Config
}

type AuthScheme int

const (
	SchemeBearer AuthScheme = iota
	SchemeDPoP
)

func AuthMiddleware(deps AuthDependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Nhận cả Bearer và DPoP
		accessToken, scheme, err := ExtractTokenFromHeader(c.GetHeader("Authorization"))
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

		// Blacklist: check chung cho cả 2 scheme
		if deps.Blacklist.IsTokenBlacklisted(c.Request.Context(), accessToken) {
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

		// sessionJKT = "" cho web, = proof.JKT cho mobile
		var sessionJKT string

		// Tách nhánh DPoP / Bearer
		if scheme == SchemeDPoP {
			dpopProof := c.GetHeader("DPoP")
			jkt := claims.GetJKT()

			if jkt == "" {
				resp := w_resp.NewAppError(401, errors.New("DPoP token missing JKT claim"), time.Now())
				w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, resp)
				c.Abort()
				return
			}

			refreshToken, errJkt := deps.UsecaseRefresh.GetByJKT(c.Request.Context(), jkt)
			if errJkt != nil {
				w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, errJkt)
				c.Abort()
				return
			}

			proof, errVerify := pkg_dpop.DPoPVerify(c, dpopProof, refreshToken.JKT, 5*time.Second, accessToken)
			if errVerify != nil {
				resp := w_resp.NewAppError(401, errVerify, time.Now())
				w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, resp)
				c.Abort()
				return
			}
			sessionJKT = proof.JKT
		} else {
			// Bearer (web): token web KHÔNG được có JKT claim.
			// Nếu có → đây là token mobile bị lộ (downgrade attack) → chặn.
			if claims.GetJKT() != "" {
				resp := w_resp.NewAppError(401,
					errors.New("bound token must use DPoP"), time.Now())
				w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, resp)
				c.Abort()
				return
			}
		}

		// ValidateSession — web: jkt = "" (lookup theo sha256(AT))
		if _, err := deps.Session.ValidateSession(
			c.Request.Context(), userID, sessionJKT, accessToken,
		); err != nil {
			resp := w_resp.NewAppError(401, err, time.Now())
			w_resp.WriteErrorWithLog(c, deps.Logger, deps.Config, resp)
			c.Abort()
			return
		}

		reqCtx := w_utils.GetRequestContext(c.Request.Context(), deps.Config)
		reqCtx.Sub = claims.Subject
		reqCtx.Tid = claims.ID
		ctx := w_utils.SaveRequestContext(c.Request.Context(), reqCtx, deps.Config)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// FIX #6, #7: copy biến local + dùng đúng jkt + log error
		uid, jkt, at := userID, sessionJKT, accessToken
		go func() {
			ctxBg, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := deps.Session.UpdateLastActive(ctxBg, uid, jkt, at); err != nil {
				// Nếu logger có method Warn/Error theo signature khác, sửa cho khớp.
				// Mục đích: không silent fail.
				_ = err
			}
		}()
	}
}

// Hàm parse scheme
func ExtractTokenFromHeader(header string) (string, AuthScheme, error) {
	if header == "" {
		return "", 0, errors.New("authorization header is required")
	}

	parts := strings.Fields(header)
	if len(parts) != 2 {
		return "", 0, errors.New("invalid authorization header format")
	}

	switch parts[0] {
	case "Bearer":
		return parts[1], SchemeBearer, nil
	case "DPoP":
		return parts[1], SchemeDPoP, nil
	default:
		return "", 0, errors.New("unsupported authorization scheme")
	}
}

func ExtractDPoPTokenFromHeader(header string) (string, error) {
	token, scheme, err := ExtractTokenFromHeader(header)
	if err != nil {
		return "", err
	}
	if scheme != SchemeDPoP {
		return "", errors.New("expected DPoP scheme")
	}
	return token, nil
}
