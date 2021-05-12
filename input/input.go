package input

import (
	"flag"
	"os"
	"strings"
)

type Input struct {
	Url string
}

func ParseInput() Input {
	var result Input
	var url string
	var help bool
	const prefix = "http://"
	flag.StringVar(&url, "u", "", "Input Target Url")
	flag.BoolVar(&help, "h", false, "Help Information")
	flag.Parse()
	if help == true {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if strings.TrimSpace(url) == "" {
		return result
	}
	if !strings.HasPrefix(url, prefix) {
		url = prefix + url
	}
	result.Url = url
	return result
}
