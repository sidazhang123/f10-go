package pubsub

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	proto "github.com/sidazhang123/f10-go/srv/processor/proto/processor"
	"github.com/sidazhang123/f10-go/srv/scheduler/proto/scheduler"
)

type Sub struct {
}

func (s *Sub) ToProcess(ctx context.Context, event *scheduler.Evt) error {
	if event.Tag == common.ProcessStart {
		if available {
			available = false
			sentTime := event.SentTime
			msg := fmt.Sprintf("Process broadcast received. sentTime %s; msg %s; available = %t", sentTime, event.Msg, available)
			log.Info(msg)
			SendLog(msg, common.LogInfolvl)

			err := processService.Refine(&proto.Request{})
			if err != nil {
				log.Error(err.Error())
				SendLog(err.Error(), common.LogErrorlvl)
			} else {
				log.Info(common.ProcessComp)
				SendLog(common.ProcessComp, common.LogCallbacklvl)
			}
			available = true
		}
	}
	return nil
}
