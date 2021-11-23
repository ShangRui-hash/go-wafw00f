package waf

import (
	"regexp"
	"strings"

	"github.com/ShangRui-hash/go-wafw00f/model"
	"github.com/ShangRui-hash/go-wafw00f/retryhttp"

	"github.com/sirupsen/logrus"
)

func DetectWaf(url string, payloads []string) (results []string, err error) {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	urlRe := regexp.MustCompile(`http.*?\..*?/`)
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
		success, result := detect(resp.Headers, resp.Body)
		if success {
			results = append(results, result)
		}
	}
	results = removeRepeatedElement(results)
	return results, nil
}

//matchHeaders 匹配响应头
func matchHeader(ruleItem model.WafItems, headers map[string]string) bool {
	//对header 没有要求
	if ruleItem.ReHeaders == nil || len(ruleItem.ReHeaders) != 0 {
		return ruleItem.Condition == "all"
	}
	flag := 0
	for key, reValue := range ruleItem.ReHeaders {
		for innerK, innerV := range headers {
			if strings.EqualFold(innerK, key) {
				if regexp.MustCompile(reValue).MatchString(innerV) {
					logrus.Debugf("condition:%srule:%s,match content:%s\n", ruleItem.Condition, reValue, innerV)
					flag++
				}
			}
		}
	}
	totalRule := len(ruleItem.ReHeaders)
	if flag == totalRule && totalRule != 0 {
		return true
	}
	return false
}

//matchCookie 匹配cookie
func matchCookie(ruleItem model.WafItems, headers map[string]string) bool {
	//对cookie 没有要求
	if ruleItem.ReCookies == nil || len(ruleItem.ReCookies) == 0 {
		return ruleItem.Condition == "all"
	}
	flag := 0
	value, ok := headers["Cookie"]
	if ok {
		for _, rule := range ruleItem.ReCookies {
			if regexp.MustCompile(rule).MatchString(value) {
				logrus.Debugf("rule:%s,match content:%s\n", rule, value)
				flag++
			}
		}
	}
	totalRule := len(ruleItem.ReCookies)
	if flag == totalRule && totalRule != 0 {
		return true
	}
	return false
}

func matchContent(ruleItem model.WafItems, content string) bool {
	//对content 没有要求
	if ruleItem.ReContent == nil || len(ruleItem.ReContent) == 0 {
		return ruleItem.Condition == "all"
	}
	for _, v := range ruleItem.ReContent {
		if regexp.MustCompile(v).MatchString(content) {
			logrus.Debugf("rule:%s,match content:%s\n", v, content)
			return true
		}
	}
	return false
}

//detect 对响应报文的headers 和body进行匹配
func detect(headers map[string]string, body string) (bool, string) {
	for _, item := range Waf.Items {
		//1.正则匹配响应头
		isMatchHeader := matchHeader(item, headers)
		//2.正则匹配cookie
		isMatchCookie := matchCookie(item, headers)
		//3.正则匹配body
		isMatchContent := matchContent(item, body)
		if item.Condition == "all" {
			if isMatchHeader && isMatchCookie && isMatchContent {
				return true, item.Name
			}
		} else if item.Condition == "any" {
			if isMatchHeader || isMatchCookie || isMatchContent {
				return true, item.Name
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
