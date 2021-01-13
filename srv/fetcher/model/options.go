package model

type Opts struct {
	Addrs              string `json:"addrs"`
	Timeout            int    `json:"timeout"`
	MaxRetry           int    `json:"max_retry"`
	Worker             int    `json:"worker"`
	DbName             string `json:"db_name"`
	GetUpdatetimeRegex string `json:"get_updatetime_regex"`
}
