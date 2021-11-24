package main

import (
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
						Destination: &settings.CurrentRunConf.URL,
					},
					&cli.StringFlag{
						Name:        "rule_file_path",
						Aliases:     []string{"r"},
						Usage:       "规则文件路径",
						Value:       "./rule.json",
						Destination: &settings.CurrentRunConf.RuleFilePath,
					},
					&cli.StringFlag{
						Name:        "proxy",
						Aliases:     []string{"p"},
						Usage:       "代理,例如:socks://127.0.0.1:7890",
						Destination: &settings.CurrentRunConf.Proxy,
					},
					&cli.BoolFlag{
						Name:        "debug",
						Aliases:     []string{"d"},
						Usage:       "调试模式",
						Destination: &settings.CurrentRunConf.Debug,
					},
					&cli.BoolFlag{
						Name:        "silent",
						Usage:       "静默模式，只输出探测结果和错误",
						Destination: &settings.CurrentRunConf.Silent,
					},
					&cli.BoolFlag{
						Name:        "stdin",
						Usage:       "从标准输入中读取url",
						Destination: &settings.CurrentRunConf.Stdin,
					},
					&cli.IntFlag{
						Name:        "routine",
						Aliases:     []string{"t"},
						Usage:       "协程数",
						Value:       1000,
						Destination: &settings.CurrentRunConf.Routine,
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
