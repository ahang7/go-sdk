package main

import (
	"github.com/ahang7/go-sdk/cli"
	"github.com/ahang7/go-sdk/cli/examples/options"
	"github.com/ahang7/go-sdk/log/zlog"
)

func App(prefix, appname string) *cli.App {
	ao := options.NewAppOption()
	app := cli.NewApp(
		prefix,
		appname,
		cli.WithFlags(ao),
		cli.WithDescription("This is a test case that provides sample code for the cli"),
		cli.WithConfig(true),
		cli.WithVersion(false),
		cli.WithRunFunc(func(app string) error {
			zlog.Info("The cli started successfully")
			return nil
		}),
		cli.WithCommand("ex",
			"short example",
			"this is a test case that provides sample code for the cli"),
	)
	return app
}

func main() {
	// 这里的prefix最终也为环境变量的前缀
	App("example", "app").Run()
}
