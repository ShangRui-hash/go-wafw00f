package model

type WafItems struct {
	Name      string            `json:"name"`
	ReHeaders map[string]string `json:"headers"`
	ReContent []string          `json:"content"`
	ReCookies []string          `json:"cookie"`
	Condition string            `json:"condition"`
}

type Waf struct {
	Items []WafItems `json:"items"`
}
