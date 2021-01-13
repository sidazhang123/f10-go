package zap

import "go.uber.org/zap"

type Options struct {
	zap.Config
	LogFileDir    string `json:"logFileDir"`
	AppName       string `json:"appName"`
	DebugFileName string `json:"debugFileName"`
	WarnFileName  string `json:"warnFileName"`
	InfoFileName  string `json:"infoFileName"`
	ErrorFileName string `json:"errorFileName"`
	MaxSize       int    `json:"maxSize"`
	MaxAge        int    `json:"maxAge"`
	MaxBackups    int    `json:"maxBackups"`
}
