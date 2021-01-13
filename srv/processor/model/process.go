package model

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sidazhang123/f10-go/basic/config"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	proto "github.com/sidazhang123/f10-go/srv/processor/proto/processor"
	"sync"
)

type Service struct {
	RegexFunc map[string]interface{} //flagName:regFunc
}

var (
	s      *Service
	m      sync.RWMutex
	Params = &Opts{}
	log    = z.GetLogger()
)

type processorService interface {
	GetPluginPath() (string, error)
	GetPluginSrc(string) (string, error)
	BuildSo(string, string) (string, error)
	RegexTest(string, string) (map[string]string, error)
	RegisterPlugin(string) error
	CallSoPlugin(string, string) (map[string]interface{}, []error)
	Refine(*proto.Request) error
	CallGoPlugin(string, string) (map[string]interface{}, []error)
}

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
	if s != nil && (Params != &Opts{}) && len(s.RegexFunc) > 0 {
		return
	}

	err := config.C().Path("params", Params)
	if err != nil {
		err = fmt.Errorf("[Init] Failed to get params\n" + err.Error())
		log.Error(err.Error())
		return
	}

	s = &Service{RegexFunc: map[string]interface{}{}}
	err = s.RegisterPlugin(Params.PluginLevel)
	if err != nil {
		err = fmt.Errorf("[Init] Failed to register plugins\n" + err.Error())
		log.Error(err.Error())
		return
	}
}
