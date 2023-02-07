package cmd

import "github.com/urfave/cli/v2"

func Run(args []string) {
	app := &cli.App{
		Name:  "web3-blog",
		Usage: "Web3 博客",
		Commands: []*cli.Command{
			{
				Name:   "runserver",
				Usage:  "启动博客服务",
				Action: RunServer,
			},
			{
				Name:   "import",
				Usage:  "博客导入",
				Action: ImportArticles,
			},
		},
	}
	if err := app.Run(args); err != nil {
		panic(err)
	}
}
