package wafw00f

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ShangRui-hash/go-wafw00f/retryhttp"
	"github.com/ShangRui-hash/go-wafw00f/util"
	jsoniter "github.com/json-iterator/go"

	"io/ioutil"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	nameRe         = regexp.MustCompile(`.*?NAME[\s]=[\s]'(.*?)'.*?`)
	headerRe       = regexp.MustCompile(`.*?self.matchHeader\(\((.*?)\)\).*?`)
	contentRe      = regexp.MustCompile(`.*?self.matchContent\(r?[',"](.*?)[',"]\).*?`)
	cookieRe       = regexp.MustCompile(`self.matchCookie\(r?'(.*?)'\).*?`)
	onlineJsonFile = "https://raw.githubusercontent.com/ShangRui-hash/go-wafw00f/master/rule.json"
)

type WafW00f struct {
	Wafs []Waf `json:"wafs"`
	json jsoniter.API
}

func NewWafW00f() *WafW00f {
	return &WafW00f{
		Wafs: make([]Waf, 0),
		json: jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

//Waf 一个WAF
type Waf struct {
	Name     string            `json:"name"`
	Relation string            `json:"relation"`
	Rules    map[string]string `json:"rules"`
}

func NewWaf() *Waf {
	return &Waf{
		Rules: make(map[string]string),
	}
}

//PraseJson 从json文件中读取规则
func (w *WafW00f) ParseJsonFile(jsonFilePath string) error {
	logrus.Debug("WAF Json Exist")
	data, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		logrus.Errorf("ioutil.ReadFile(%s) failed,err:%v", jsonFilePath, err)
		resp, err := retryhttp.Get(onlineJsonFile)
		if err != nil {
			logrus.Errorf("retryhttp.Get(%s) failed,err:%v", onlineJsonFile, err)
			return err
		}
		data = []byte(resp.Body)
	}
	if err := w.json.Unmarshal(data, &w.Wafs); err != nil {
		logrus.Error("WAF Unmarshal Error")
		return err
	}
	return nil
}

//ParseLib 解析对应目录下的所有规则文件到内存
func (w *WafW00f) ParseLib(libDirPath string) error {
	failCount := 0
	successCount := 0
	//1.读取lib文件
	filenames, err := util.GetAllFile(libDirPath)
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
			continue
		}
		//todo:有些文件中schema有两个，需要特殊处理
		if strings.Contains(string(content), "schema1") || strings.Contains(string(content), "schema2") {
			logrus.Warn("请手动拆分:", filenames[i])
			failCount++
			continue
		}
		//解析lib文件
		waf, err := w.parsePythonFile(string(content))
		if err != nil {
			logrus.Errorf("doParse(%s) failed,err:%v", content, err)
			continue
		}
		w.Wafs = append(w.Wafs, waf)
		successCount++
	}
	logrus.Debug("failCount:", failCount, "successCount:", successCount)
	return nil
}

//parsePythonFile 解析单个python文件
func (w *WafW00f) parsePythonFile(input string) (Waf, error) {
	waf := NewWaf()
	//1.解析NAME
	nameMatch := nameRe.FindAllStringSubmatch(input, -1)
	if len(nameMatch) == 0 {
		return *waf, fmt.Errorf("parsePythonFile failed,name not found")
	}
	waf.Name = nameMatch[0][1]
	//2.解析HEADER
	headerMatch := headerRe.FindAllStringSubmatch(input, -1)
	for _, v := range headerMatch {
		header := v[1]
		key := strings.TrimLeft(strings.Split(header, "',")[0], "'")
		key = strings.TrimSpace(key)
		value := strings.TrimSpace(strings.Split(header, "',")[1])
		value = strings.TrimLeft(value, "r")
		value = strings.Trim(value, "'")
		waf.Rules[value] = fmt.Sprintf("header::%s", key)
	}
	//3.解析CONTENT
	contentMatch := contentRe.FindAllStringSubmatch(input, -1)
	logrus.Debug("contentMatch:", contentMatch)
	for _, v := range contentMatch {
		content := v[1]
		logrus.Debug("content:", content)
		waf.Rules[content] = "content"
	}
	//4.解析COOKIE
	cookieMatch := cookieRe.FindAllStringSubmatch(input, -1)
	for _, v := range cookieMatch {
		cookie := v[1]
		cookie = strings.TrimLeft(cookie, "r")
		cookie = strings.Trim(cookie, "'")
		waf.Rules[cookie] = "header::cookie"
	}
	//5.解析满足条件
	if strings.Contains(input, "any") {
		waf.Relation = "any"
	} else {
		waf.Relation = "all"
	}
	if len(waf.Rules) == 0 {
		return *waf, errors.New("parsePythonFile failed,rules not found")
	}
	return *waf, nil
}

//DumpJson 将内存中的规则输出为json文件
func (w *WafW00f) DumpJson(dstFilePath string) error {
	buffer := &bytes.Buffer{}
	encoder := w.json.NewEncoder(buffer)
	//关闭HTML实体编码
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(w.Wafs); err != nil {
		logrus.Error("json.Marshal failed,err:", err)
		return err
	}
	if err := ioutil.WriteFile(dstFilePath, buffer.Bytes(), 0777); err != nil {
		logrus.Error("WriteFile failed,err:", err)
		return err
	}
	return nil
}
