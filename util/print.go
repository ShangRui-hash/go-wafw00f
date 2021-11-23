package util

import (
	"bytes"
	"fmt"
)

func PrintWafResult(items []string) {
	if len(items) == 0 {
		fmt.Println("WAF Not Detected")
		return
	}
	temp := bytes.Buffer{}
	temp.WriteString("WAF : [")
	for i, v := range items {
		if i == len(items)-1 {
			temp.WriteString(v + "]")
		} else {
			temp.WriteString(v + ",")
		}
	}
	fmt.Println(temp.String())
}
