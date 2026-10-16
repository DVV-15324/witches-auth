package refresh

import (
	core "example/cmd/server/core"
	usecase_auth "example/internal/auth/usecase"
	dtoRefresh "example/internal/refresh/dto/request"
	handler "example/internal/refresh/handler"
	"example/internal/refresh/repository"
	usecase_refresh "example/internal/refresh/usecase"
	usecase_user "example/internal/user/usecase"
	w_handl "github.com/DVV-15324/witches/pkg/core/handler"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"
	"github.com/ulule/limiter/v3"
)

type RefreshModule struct {
	Core    *core.CoreServices
	Handler *handler.RefreshHandle
	Usecase *usecase_refresh.RefreshUseCase
	Repo    *repository.RefreshTokenRepository
}

func NewRefreshModule(
	core *core.CoreServices,
	userUsecase *usecase_user.UserUseCase,
) *RefreshModule {
	repo := repository.NewRefreshTokenRepository(core)
	usecase := usecase_refresh.NewRefreshUseCase(core, userUsecase, repo)
	handler := handler.NewRefreshHandle(core, usecase)
	return &RefreshModule{
		Core:    core,
		Handler: handler,
		Usecase: usecase,
		Repo:    repo,
	}
}
func (m *RefreshModule) SetAuthUseCase(authUsecase *usecase_auth.AuthUseCase) {
	m.Usecase.SetAuthUseCase(authUsecase)
}

func (m *RefreshModule) RegisterPublicRoutes(gen *w_handl.SwaggerGenerator) {

	gen.POST("/v1/auth/refresh").
		Summary("Refresh access token").
		Description("Get new access token using refresh token").
		Tags("auth").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		RateLimit(limiter.Rate{Period: m.Core.Limiter.Rate.Period, Limit: m.Core.Limiter.Rate.Limit}).
		Body(dtoRefresh.RefreshTokenRequest{}, "Refresh token").
		Response(200, w_resp.AppResponse{}, "Token refreshed successfully").
		Response(401, w_resp.AppResponse{}, "Invalid or expired refresh token").
		Handler(m.Handler.HandleRefreshToken()).
		Build()
}

func (m *RefreshModule) RegisterProtectedRoutes(gen *w_handl.SwaggerGenerator) {
	// Protected routes của auth
	// ... (giữ nguyên từ protected.go)
}
