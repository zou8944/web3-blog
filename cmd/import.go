package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/project5e/web3-blog/bootstrap"
	"github.com/project5e/web3-blog/pkg/import/mywords"
	"github.com/urfave/cli/v2"
)

func ImportArticles(ctx *cli.Context) error {
	g := gin.New()
	bootstrap.SetupAll(g)
	gitAddress := ctx.Args().First()
	return mywords.Import(gitAddress)
}
