package controllers

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/ShangRui-hash/go-wafw00f/logger"
	"github.com/ShangRui-hash/go-wafw00f/model"
	"github.com/ShangRui-hash/go-wafw00f/settings"
	"github.com/ShangRui-hash/go-wafw00f/util"
	jsoniter "github.com/json-iterator/go"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var (
	nameRe    = regexp.MustCompile(`.*?NAME[\s]=[\s]'(.*?)'.*?`)
	headerRe  = regexp.MustCompile(`.*?self.matchHeader\(\((.*?)\)\).*?`)
	contentRe = regexp.MustCompile(`.*?self.matchContent\(r'(.*?)'\).*?`)
	cookieRe  = regexp.MustCompile(`self.matchCookie\(r'(.*?)'\).*?`)
)

func Parse(c *cli.Context) error {
	var Waf model.Waf
	successCount := 0
	failCount := 0
	//0.初始化logger配置
	logger.Init(settings.CurrentParseConf.Debug)
	//1.读取lib文件
	filenames, err := util.GetAllFile(settings.CurrentParseConf.LibDirPath)
	if err != nil {
		logrus.Error("util.GetAllFile failed,err:", err)
		return err
	}
	//解析lib文件都Waf中
	for i := range filenames {
		if !strings.HasSuffix(filenames[i], ".lib") && !strings.HasSuffix(filenames[i], ".py") {
			continue
		}
		content, err := ioutil.ReadFile(filenames[i])
		if err != nil {
			logrus.Error("ioutil.ReadFile failed,err:", err)
			return err
		}
		//todo:有些文件中schema有两个，需要特殊处理
		if strings.Contains(string(content), "schema1") || strings.Contains(string(content), "schema2") {
			logrus.Warn("请手动拆分:", filenames[i])
			failCount++
			continue
		}
		//解析lib文件
		item, err := doParse(string(content))
		if err != nil {
			logrus.Error("doParse failed,err:", err)
			return err
		}
		Waf.Items = append(Waf.Items, item)
		successCount++
	}
	//2.将Waf中的数据写入到json文件中
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(Waf); err != nil {
		logrus.Error("json.Marshal failed,err:", err)
		return err
	}
	if err := ioutil.WriteFile(settings.CurrentParseConf.RuleFilePath, buffer.Bytes(), 0777); err != nil {
		logrus.Error("WriteFile failed,err:", err)
		return err
	}
	logrus.Info("Create WAF Json File success!")
	logrus.Infof("filepath:%s\tsuccess:%d\tfail:%d\ttotal:%d\n", settings.CurrentParseConf.RuleFilePath, successCount, failCount, len(filenames))
	return nil
}

//doParse 解析lib文件
func doParse(input string) (model.WafItems, error) {
	result := model.WafItems{
		ReHeaders: make(map[string]string),
	}
	//1.解析NAME
	nameMatch := nameRe.FindAllStringSubmatch(input, -1)
	result.Name = nameMatch[0][1]
	//2.解析HEADER
	headerMatch := headerRe.FindAllStringSubmatch(input, -1)
	for _, v := range headerMatch {
		header := v[1]
		key := strings.TrimLeft(strings.Split(header, "',")[0], "'")
		key = strings.TrimSpace(key)
		value := strings.TrimSpace(strings.Split(header, "',")[1])
		value = strings.TrimLeft(value, "r")
		value = strings.Trim(value, "'")
		result.ReHeaders[key] = value
	}
	//3.解析CONTENT
	contentMatch := contentRe.FindAllStringSubmatch(input, -1)
	logrus.Debug("contentMatch:", contentMatch)
	for _, v := range contentMatch {
		content := v[1]
		logrus.Debug("content:", content)
		result.ReContent = append(result.ReContent, content)
	}
	//4.解析COOKIE
	cookieMatch := cookieRe.FindAllStringSubmatch(input, -1)
	for _, v := range cookieMatch {
		cookie := v[1]
		cookie = strings.TrimLeft(cookie, "r")
		cookie = strings.Trim(cookie, "'")
		result.ReCookies = append(result.ReContent, cookie)
	}
	//5.解析满足条件
	if strings.Contains(input, "any") {
		result.Condition = "any"
	} else {
		result.Condition = "all"
	}
	return result, nil

}
