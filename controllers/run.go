package controllers

import (
	"fmt"
	"sync"

	"github.com/ShangRui-hash/go-wafw00f/logger"
	"github.com/ShangRui-hash/go-wafw00f/retryhttp"
	"github.com/ShangRui-hash/go-wafw00f/settings"
	"github.com/ShangRui-hash/go-wafw00f/wafw00f"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {
	urlCh := make(chan string, 10240)
	//0.初始化logger配置
	logger.Init(settings.CurrentRunConf.Debug, settings.CurrentRunConf.Silent)
	//0.初始化retryhttp客户端
	if err := retryhttp.Init(); err != nil {
		logrus.Error("retryhttp init error:", err)
		return err
	}
	w00f := wafw00f.NewWafW00f()
	//1.解析waf识别规则文件
	if err := w00f.ParseJsonFile(settings.CurrentRunConf.RuleFilePath); err != nil {
		logrus.Error("waf.Init failed,err:", err)
		return err
	}
	//2.开始探测waf类型
	//消费者：消费url
	var wg sync.WaitGroup
	w00f.DetectWaf(urlCh, &wg)
	//生产者：生产url
	URLgener := wafw00f.NewURLGenerator()
	if err := URLgener.GenURL(urlCh); err != nil {
		return err
	}
	wg.Wait()
	fmt.Println("Not Detected:", w00f.NotDetectedURL.ToSlice())
	return nil
}
