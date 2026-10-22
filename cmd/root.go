package cmd

import (
	"context"
	"fmt"
	"github.com/DVV-15324/witches-auth/cmd/server/config"
	routes "github.com/DVV-15324/witches-auth/cmd/server/routers"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "root",
	Short: "En: Start the software | Vi: Bắt đầu khởi động phần mềm",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		config := config.Load()
		r = w_utils.InstanceMetric(fmt.Sprintf("%d", config.MetrictPort), config.Host, r)
		routes.Start(r, config)
		w_utils.ShutdownServer(context.Background(), r, config.Host, fmt.Sprintf("%d", config.Port))
	},
}

func GetExcute() *cobra.Command {
	return root
}
