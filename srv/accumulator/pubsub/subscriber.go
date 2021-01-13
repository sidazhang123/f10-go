package pubsub

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/srv/scheduler/proto/scheduler"
	"sync"
	"time"
)

type Sub struct {
}

func (s *Sub) ToAccumulate(ctx context.Context, event *scheduler.Evt) error {
	if event.Tag == common.AccumulateStart {
		sentTime := event.SentTime
		msg := fmt.Sprintf("Accumulate broadcast received. sentTime %s; msg %s", sentTime, event.Msg)
		log.Info(msg)
		SendLog(msg, common.LogInfolvl)
		var wg sync.WaitGroup
		errList := accumulatorService.DoAll("", time.Now().Format(common.TimestampLayout[:10]),
			time.Now().AddDate(0, 0, 1).Format(common.TimestampLayout[:10]), &wg)
		if len(errList) > 0 {
			for _, err := range errList {
				log.Error(err.Error())
				SendLog(err.Error(), common.LogErrorlvl)
			}
		}

		log.Info(common.AccumulateComp)
		SendLog(common.AccumulateComp, common.LogCallbacklvl)

	}
	return nil
}
