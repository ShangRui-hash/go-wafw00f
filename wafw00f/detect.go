package wafw00f

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/ShangRui-hash/go-wafw00f/retryhttp"
	"github.com/ShangRui-hash/go-wafw00f/settings"

	"github.com/sirupsen/logrus"
)

func (w *WafW00f) DetectWaf(urlCh chan string, wg *sync.WaitGroup) {
	respCh := make(chan *retryhttp.Resp, 10240)
	//消费者：消费响应报文
	for i := 0; i < settings.CurrentRunConf.Routine; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for resp := range respCh {
				success, wafName, reasons := w.detect(resp.Headers, resp.Body)
				if success {
					_, ok := w.Knowledge.Load(resp.URL)
					if !ok {
						output := strings.Join([]string{resp.URL, wafName, strings.Join(reasons, " && ")}, " || ")
						fmt.Println(output)
						w.Knowledge.Store(resp.URL, wafName)
						w.NotDetectedURL.Remove(resp.URL)
					}
				} else {
					w.NotDetectedURL.Add(resp.URL)
				}
			}
		}()
	}
	//生产者：生产响应报文
	var respWg sync.WaitGroup
	for i := 0; i < settings.CurrentRunConf.Routine; i++ {
		respWg.Add(1)
		go func() {
			defer respWg.Done()
			for URL := range urlCh {
				u, err := url.Parse(URL)
				if err != nil {
					continue
				}
				_, ok := w.Knowledge.Load(u.Scheme + "://" + u.Host)
				if ok {
					continue
				}
				resp, err := retryhttp.Get(URL)
				if err != nil {
					logrus.Warn("retryhttp.Get failed,err:", err)
					continue
				}
				respCh <- resp
			}
		}()
	}
	//负责等待所有生产者工作完成，关闭写通道
	go func() {
		respWg.Wait()
		close(respCh)
	}()
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
func (w *WafW00f) detect(headers map[string]string, body string) (bool, string, []string) {
	for _, waf := range w.Wafs {
		matched := 0
		var reasons []string
		for ruleReg, location := range waf.Rules {
			if location == "content" {
				if w.matchContent(body, ruleReg) {
					reasons = append(reasons, "规则:"+ruleReg+",位置:"+location)
					matched++
				}
			} else if strings.HasPrefix(location, "header::") {
				if w.matchHeader(headers, location, ruleReg) {
					reasons = append(reasons, "规则:"+ruleReg+",位置:"+location)
					matched++
				}
			}
		}

		if waf.Relation == "all" {
			if matched == len(waf.Rules) {
				return true, waf.Name, reasons
			}
		} else if waf.Relation == "any" {
			if matched > 0 {
				return true, waf.Name, reasons
			}
		}
	}
	return false, "Not Detected", nil
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
