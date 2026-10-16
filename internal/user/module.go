package user

import (
	core "example/cmd/server/core"

	handler "example/internal/user/handler"
	"example/internal/user/repository"
	usecase_user "example/internal/user/usecase"

	w_handl "github.com/DVV-15324/witches/pkg/core/handler"
	w_resp "github.com/DVV-15324/witches/pkg/core/response"

	"github.com/ulule/limiter/v3"
)

type UserModule struct {
	Core    *core.CoreServices
	Handler *handler.UserHandle
	Usecase *usecase_user.UserUseCase
	Repo    *repository.UserRepository
}

func NewUserModule(
	core *core.CoreServices,
) *UserModule {
	repo := repository.NewUserRepository(core)
	usecase := usecase_user.NewUserUsecase(core, repo)
	handler := handler.NewUserHandle(core, usecase)
	return &UserModule{
		Core:    core,
		Handler: handler,
		Usecase: usecase,
		Repo:    repo,
	}
}

func (m *UserModule) RegisterPublicRoutes(gen *w_handl.SwaggerGenerator) {}

func (m *UserModule) RegisterProtectedRoutes(gen *w_handl.SwaggerGenerator) {
	gen.POST("/v1/user/get_user").
		Summary("Get user by ID").
		Description("Get user details by user ID").
		Tags("user").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		HeaderParam("DPoP", "DPoP proof", false).
		RateLimit(limiter.Rate{Period: m.Core.Limiter.Rate.Period, Limit: m.Core.Limiter.Rate.Limit}).
		HeaderParam("Authorization", "DPoP access token for authentication", false).
		HeaderParam("DPoP", "DPoP proof JWT", false).
		Body(map[string]int{"id": 0}, "User ID").
		Response(200, w_resp.AppResponse{}, "User found").
		Response(404, w_resp.AppResponse{}, "User not found").
		Use(m.Core.Middleware).
		Handler(m.Handler.HandleGetUserById()).
		Build()

	gen.POST("/v1/user/get_user_all").
		Summary("Get all users").
		Description("Get list of all users").
		Tags("user").
		HeaderParam("User-Agent", "User-Agent", false).
		HeaderParam("Accept-Language", "Preferred language (e.g., vi-VN, en-US)", false).
		HeaderParam("X-Timezone", "Timezone (e.g., Asia/Ho_Chi_Minh)", false).
		HeaderParam("Authorization", "DPoP access token for authentication", false).
		HeaderParam("DPoP", "DPoP proof JWT", false).
		RateLimit(limiter.Rate{Period: m.Core.Limiter.Rate.Period, Limit: m.Core.Limiter.Rate.Limit}).
		QueryParam("page", "Pagenation request", false).
		QueryParam("limit", "Pagenation request", false).
		QueryParam("search", "Pagenation request", false).
		QueryParam("sort", "Pagenation request", false).
		QueryParam("order", "Pagenation request", false).
		Response(200, w_resp.AppResponse{}, "Users list").
		Use(m.Core.Middleware).
		Handler(m.Handler.HandleGetAllUser()).
		Build()
}
