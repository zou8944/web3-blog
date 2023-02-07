package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/project5e/web3-blog/app/route"
	"github.com/project5e/web3-blog/bootstrap"
	"github.com/project5e/web3-blog/config"
	"github.com/urfave/cli/v2"
)

func RunServer(_ *cli.Context) error {
	g := gin.New()
	bootstrap.SetupAll(g)
	route.RegisterRoutes(g)
	address := fmt.Sprintf("0.0.0.0:%d", config.Server.Port)
	return g.Run(address)
}
