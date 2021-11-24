package controllers

import (
	"github.com/ShangRui-hash/go-wafw00f/wafw00f"

	"github.com/ShangRui-hash/go-wafw00f/logger"
	"github.com/ShangRui-hash/go-wafw00f/settings"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func Parse(c *cli.Context) error {
	//0.初始化logger配置
	logger.Init(settings.CurrentParseConf.Debug)
	w00f := wafw00f.NewWafW00f()
	//1.读取LibDirPath下所有的规则文件，并将他们解析到内存中
	if err := w00f.ParseLib(settings.CurrentParseConf.LibDirPath); err != nil {
		logrus.Error("parse lib file error:", err)
		return err
	}
	//2.将Waf中的数据写入到json文件中
	if err := w00f.DumpJson(settings.CurrentParseConf.RuleFilePath); err != nil {
		logrus.Error("dump json file error:", err)
		return err
	}
	logrus.Info("Create WAF Json File success!")
	return nil
}
