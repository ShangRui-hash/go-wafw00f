package retryhttp

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ShangRui-hash/go-wafw00f/settings"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

var client *http.Client

func Init() error {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = logrus.StandardLogger()
	tr := &http.Transport{
		//跳过证书验证
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//代理
	if settings.CurrentRunConf.Proxy != "" {
		proxy, err := url.Parse(settings.CurrentRunConf.Proxy)
		if err != nil {
			logrus.Error("url.Parse failed,err:", err)
			return err
		}
		tr.Proxy = http.ProxyURL(proxy)
	}
	retryClient.HTTPClient.Transport = tr
	client = retryClient.StandardClient()
	return nil
}

type Resp struct {
	StatusCode int
	URL        string
	Body       string
	Headers    map[string]string
}

func Get(u string) (resp *Resp, err error) {
	URL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	response, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	for k, v := range response.Header {
		headers[k] = strings.Join(v, "")
	}

	return &Resp{
		StatusCode: response.StatusCode,
		Body:       string(body),
		Headers:    headers,
		URL:        URL.Scheme + "://" + URL.Host,
	}, nil
}
