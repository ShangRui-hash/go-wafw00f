package main

import (
	"strings"

	"github.com/ShangRui-hash/go-wafw00f/input"
	"github.com/ShangRui-hash/go-wafw00f/log"
	"github.com/ShangRui-hash/go-wafw00f/util"
	"github.com/ShangRui-hash/go-wafw00f/waf"
)

func main() {
	util.PrintLogo()
	params := input.ParseInput()
	if strings.TrimSpace(params.Url) == "" {
		log.Error("Need Url")
		return
	}
	log.Default("url : " + params.Url)
	log.Default("Start go-wafw00f")
	waf.ResolveWafLib()
	res := waf.DetectWaf(params.Url)
	util.PrintWafResult(res)
}
