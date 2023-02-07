package cmd

import (
	"github.com/project5e/web3-blog/pkg/import/mywords"
	"github.com/urfave/cli/v2"
)

func ImportArticles(ctx *cli.Context) error {
	gitAddress := ctx.Args().First()
	return mywords.Import(gitAddress)
}
