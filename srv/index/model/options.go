package model

type Opts struct {
	Addrs    string `json:"addrs"`
	Timeout  int    `json:"timeout"`
	MaxRetry int    `json:"max_retry"`
	FlagName string `json:"flag_name"`
}
