package main

import (
	"fmt"
	"go/build"
	"os"

	"github.com/ShangRui-hash/go-wafw00f/controllers"
	"github.com/ShangRui-hash/go-wafw00f/settings"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func main() {
	authors := []*cli.Author{
		{
			Name:  "无在无不在",
			Email: "2227627947@qq.com",
		},
		{
			Name: "EmYiQing(原作者)",
		},
	}
	app := &cli.App{
		Name:      "go-wafw00f",
		Usage:     "go-wafw00f",
		UsageText: "go-wafw00f",
		Version:   "v0.0.1",
		Authors:   authors,
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "探测waf类型",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "url",
						Aliases:     []string{"u"},
						Usage:       "目标url",
						Required:    true,
						Destination: &settings.CurrentRunConf.URL,
					},
					&cli.StringFlag{
						Name:        "rule_file_path",
						Aliases:     []string{"r"},
						Usage:       "规则文件路径",
						Value:       fmt.Sprintf("%s/pkg/mod/github.com/!shang!rui-hash/go-wafw00f@v0.0.1/rule.json", build.Default.GOPATH),
						Destination: &settings.CurrentRunConf.RuleFilePath,
					},
					&cli.BoolFlag{
						Name:        "debug",
						Aliases:     []string{"d"},
						Usage:       "调试模式",
						Destination: &settings.CurrentRunConf.Debug,
					},
				},
				Action: controllers.Run,
			},
			{
				Name:  "parse",
				Usage: "解析wafw00f的规则库为json格式",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "rule_file_path",
						Aliases:     []string{"dst"},
						Usage:       "规则文件路径",
						DefaultText: "./rule.json",
						Destination: &settings.CurrentParseConf.RuleFilePath,
					},
					&cli.StringFlag{
						Name:        "lib_dir_path",
						Aliases:     []string{"lib"},
						Usage:       "wafw00f规则库文件夹的路径",
						Required:    true,
						Destination: &settings.CurrentParseConf.LibDirPath,
					},
					&cli.BoolFlag{
						Name:        "debug",
						Aliases:     []string{"d"},
						Usage:       "调试模式",
						Destination: &settings.CurrentParseConf.Debug,
					},
				},
				Action: controllers.Parse,
			},
		},
	}
	//启动app
	if err := app.Run(os.Args); err != nil {
		logrus.Error(err)
		return
	}
}
