package wafw00f

import (
	"regexp"
	"strings"

	"github.com/ShangRui-hash/go-wafw00f/retryhttp"

	"github.com/sirupsen/logrus"
)

var urlRe = regexp.MustCompile(`http.*?\..*?/`)

func (w *WafW00f) DetectWaf(url string, payloads []string) (results []string, err error) {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	match := urlRe.FindAllStringSubmatch(url, -1)
	rootUrl := strings.TrimRight(match[0][0], "/")
	for _, v := range payloads {
		url := rootUrl + v
		resp, err := retryhttp.Get(url)
		if err != nil {
			logrus.Error("retryhttp.Get failed,err:", err)
			return nil, err
		}
		logrus.Debugf("requestURL:%s,statusCode:%d\nresp headers:%s\nresp_body:%s\n", rootUrl+v, resp.StatusCode, resp.Headers, resp.Body)
		success, wafName := w.detect(resp.Headers, resp.Body)
		if success {
			results = append(results, wafName)
		}
	}
	results = removeRepeatedElement(results)
	return results, nil
}

func (w *WafW00f) matchContent(body, ruleReg string) bool {
	return regexp.MustCompile(ruleReg).MatchString(body)
}

func (w *WafW00f) matchHeader(headers map[string]string, location, ruleReg string) bool {
	headerKey := strings.Split(location, "::")[1]
	headerValue, ok := headers[headerKey]
	if !ok {
		return false
	}
	return regexp.MustCompile(ruleReg).MatchString(headerValue)
}

//detect 对响应报文的headers 和body进行匹配
func (w *WafW00f) detect(headers map[string]string, body string) (bool, string) {
	for _, waf := range w.Wafs {
		matched := 0
		for ruleReg, location := range waf.Rules {
			if location == "content" {
				if w.matchContent(body, ruleReg) {
					matched++
				}
			} else if strings.HasPrefix(location, "header::") {
				if w.matchHeader(headers, location, ruleReg) {
					matched++
				}
			}
		}

		if waf.Relation == "all" {
			if matched == len(waf.Rules) {
				return true, waf.Name
			}
		} else if waf.Relation == "any" {
			if matched > 0 {
				return true, waf.Name
			}
		}
	}
	return false, "Not Detected"
}

func removeRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
