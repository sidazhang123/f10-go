package pubsub

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/srv/scheduler/proto/scheduler"
)

type Sub struct{}

func (s *Sub) ToSchedule(ctx context.Context, l *scheduler.Log) error {
	tag := l.Tag
	msg := l.Msg
	switch l.Level {
	case common.LogInfolvl:
		log.Info(fmt.Sprintf("[%s] %s", tag, msg))
	case common.LogErrorlvl:
		emsg := fmt.Sprintf("[%s] %s", tag, msg)
		log.Error(emsg)
		err := SendAlarm(emsg)
		if err != nil {
			log.Error(err.Error())
		}
	case common.LogCallbacklvl:
		log.Info(fmt.Sprintf("[%s] %s", tag, msg))
		SendControl("", common.CommandChain[tag])
		log.Info(fmt.Sprintf("SendControl [%s]", common.CommandChain[tag]))
	}
	return nil
}
