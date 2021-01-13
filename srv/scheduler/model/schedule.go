package model

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sidazhang123/f10-go/basic/common"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

type Service struct {
}

var (
	s *Service
	m sync.RWMutex
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
	s = &Service{}
	err, name, t := s.SetCron()
	if err != nil {
		log.Error("[setCron]Error - ", zap.Any("err", err))
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("[Init] Completed. Next Task [%s] starts at [%s]",
		name, t.Format(common.TimestampLayout)))
}

type SchedulerService interface {
	SetCron() (err error, name string, t time.Time)
	GetCodeName()
	DeleteOutdatedFocus()
	RemoveTask(funcName string) (e error, name string, t time.Time)
	ClearCron() (string, time.Time)
	Once(funcName string) (e error, name string, t time.Time)
	NextScheduledTask() (string, time.Time)
	SendAlarm(msg string) error
}
