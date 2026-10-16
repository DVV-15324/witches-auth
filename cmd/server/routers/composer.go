package routers

import (
	"example/cmd/server/core"
	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
)

func Start(r *gin.Engine, cfg *wcmd_utils.Config) {
	core := core.NewCoreServices(cfg)
	modules := InitModules(core)
	RegisterRoutes(r, core, modules)
}
