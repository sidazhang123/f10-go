package model

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sidazhang123/f10-go/basic/config"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	"sync"
)

type fetcherService interface {
	GetByMQ() error
	GetByInput()
}
type Service struct {
}

var (
	s    *Service
	m    sync.RWMutex
	opts = &Opts{}
	log  = z.GetLogger()
)

func GetService() (service *Service, err error) {
	if s == nil {
		err = errors.Errorf("[GetService] GetService is not initialized.")
		return
	}
	return s, nil
}

func Init() {
	m.Lock()
	defer m.Unlock()
	if s != nil {
		return
	}
	err := config.C().Path("params", opts)
	if err != nil {
		err = fmt.Errorf("[Init] Failed to get params\n" + err.Error())
		log.Error(err.Error())
		return
	}

	s = &Service{}
}
