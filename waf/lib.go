package waf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ShangRui-hash/go-wafw00f/model"
	"github.com/sirupsen/logrus"
)

var (
	Waf model.Waf
)

//Init 解析用于waf识别的lib文件
func Init(filename string) error {
	var ret model.Waf
	logrus.Debug("WAF Json Exist")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Errorf("ioutil.ReadFile(%s) failed,err:%v", filename, err)
		return err
	}
	if err := json.Unmarshal(data, &ret); err != nil {
		logrus.Error("WAF Unmarshal Error")
		return err
	}
	Waf = ret
	return nil
}
