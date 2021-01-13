package model

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sidazhang123/f10-go/basic/config"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	"strings"
	"sync"
)

type Service struct {
}

var (
	s      *Service
	m      sync.RWMutex
	Params = &Opts{}
	log    = z.GetLogger()
	// field name to collection name in f10-acc
	AccFieldCollectionMap = make(map[string]string)
)

type accumulateService interface {
	DoAll(string, string, string, *sync.WaitGroup) []error
	DoSA(string, string, string) error
	DoFA(string, string, string) error
	DoOA(string, string, string) error
	DoLT(string, string, string) error
	ReprAll() []error
	GetRepr(string) string
}

func makeAFCMap() {
	for _, f := range strings.Split(Params.Win1Seq, ",") {
		if len(strings.TrimSpace(f)) > 0 {
			AccFieldCollectionMap[f] = Params.Win1Name
		}
	}
	for _, f := range strings.Split(Params.Win2Seq, ",") {
		if len(strings.TrimSpace(f)) > 0 {
			AccFieldCollectionMap[f] = Params.Win2Name
		}
	}
	for _, f := range strings.Split(Params.Win3Seq, ",") {
		if len(strings.TrimSpace(f)) > 0 {
			AccFieldCollectionMap[f] = Params.Win3Name
		}
	}
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
	if s != nil && (Params != &Opts{}) {
		return
	}

	err := config.C().Path("params", Params)
	if err != nil {
		err = fmt.Errorf("[Init] Failed to get params\n" + err.Error())
		log.Error(err.Error())
		return
	}
	makeAFCMap()
	s = &Service{}
	err = initReprJson()
	if err != nil {
		err = fmt.Errorf("[Init] Failed to initReprJson\n" + err.Error())
		log.Error(err.Error())
		return
	}
	log.Info("[initReprJson] complete")
}
