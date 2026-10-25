package routers

import (
	"fmt"
	"log"

	"github.com/DVV-15324/witches-auth/cmd/server/core"
	pkg_middleware "github.com/DVV-15324/witches-auth/internal/auth/middleware"
	"github.com/DVV-15324/witches-auth/internal/shared/middleware"
	w_handl "github.com/DVV-15324/witches/pkg/core/handler"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, core *core.CoreServices, modules *Modules) {
	// Middleware CORS
	r.Use(w_utils.Cors(core.Config))
	r.Use(middleware.TimingMiddleware(core.Logger, core.Config))

	// Rate limit middleware
	rateLimitMiddleware := &middleware.LimitMiddleWare{}

	// Auth middleware
	authMiddleware := pkg_middleware.AuthMiddleware(
		pkg_middleware.AuthDependencies{
			UsecaseRefresh: modules.Refresh.Usecase,
			JWT:            core.JWT,
			Blacklist:      core.Blacklist,
			Session:        core.Session,
			Logger:         core.Logger,
			Config:         core.Config,
		},
	)
	core.Middleware = authMiddleware

	// Khởi tạo Swagger
	gen := initSwagger(r, core, rateLimitMiddleware)
	initModule(modules, gen)
	if err := gen.Save("swagger.json"); err != nil {
		log.Printf("Error saving swagger.json: %v", err)
	} else {
		log.Println("swagger.json generated successfully!")
	}

	// Swagger UI
	r.GET("/swagger/*any", w_handl.SwaggerUI())
	r.GET("/swagger.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(200, gen.GenerateJSON())
	})
}

func initSwagger(r *gin.Engine, core *core.CoreServices, rateLimitMiddleware w_handl.IRateLimitMiddleware) *w_handl.SwaggerGenerator {
	address := fmt.Sprintf("%s:%d", core.Config.Host, core.Config.Port)
	gen := w_handl.NewSwaggerGenerator("github.com/DVV-15324/witches-auth API", "1.0", address, "/")
	gen.SetEngine(r)
	gen.SetRedisClient(core.Redis.GetClient())
	gen.SetRateLimitMiddleware(rateLimitMiddleware)
	return gen
}

func initModule(modules *Modules, gen *w_handl.SwaggerGenerator) {
	gen.AddTag("v1", "API Version 1")
	gen.AddTag("auth", "Auth endpoints")
	modules.Auth.RegisterProtectedRoutes(gen)
	modules.Auth.RegisterPublicRoutes(gen)
	gen.AddTag("user", "User endpoints")
	modules.User.RegisterProtectedRoutes(gen)
	modules.User.RegisterPublicRoutes(gen)
	gen.AddTag("refresh", "Refresh endpoints")
	modules.Refresh.RegisterPublicRoutes(gen)
	modules.Refresh.RegisterProtectedRoutes(gen)

}
