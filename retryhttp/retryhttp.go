package retryhttp

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

var client *http.Client

func Init() {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = logrus.StandardLogger()
	//跳过证书验证
	retryClient.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = retryClient.StandardClient()
}

type Resp struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

func Get(url string) (resp *Resp, err error) {
	response, err := client.Get(url)
	if err != nil {
		logrus.Error("client.Get failed,err:", err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.Error("ioutil.ReadAll failed,err:", err)
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
	}, nil
}
