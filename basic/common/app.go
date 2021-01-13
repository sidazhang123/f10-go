package common

import "strconv"

type AppCfg struct {
	Name    string  `json:"name"`
	Version float32 `json:"version"`
	Address string  `json:"addr"`
	Port    int     `json:"port"`
}

func (a *AppCfg) Addr() string {
	return a.Address + ":" + strconv.Itoa(a.Port)
}
