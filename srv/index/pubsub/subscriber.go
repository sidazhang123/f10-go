package pubsub

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"time"

	"github.com/sidazhang123/f10-go/srv/scheduler/proto/scheduler"
)

type Sub struct{}

func (s *Sub) ToIndex(ctx context.Context, event *scheduler.Evt) error {
	if event.Tag == common.IndexStart {
		if time.Now().After(lastRun.Time().Add(interval)) {
			lastRun.SetTime(time.Now())
			sentTime := event.SentTime
			msg := fmt.Sprintf("Index broadcast received. sentTime %s; msg %s;", sentTime, event.Msg)
			log.Info(msg)
			SendLog(msg, common.LogInfolvl)

			err, res := indexService.Fetch()
			if err != nil {
				log.Error(err.Error())
				SendLog("err occurred while getting codename\n"+err.Error(), common.LogErrorlvl)
			} else {
				m := fmt.Sprintf("%s %d stocks on the market today.", common.IndexComp, len(res))
				log.Info(m)
				SendLog(m, common.LogCallbacklvl)
			}
		} else {
			SendLog(fmt.Sprintf("Index IS RUNNING. sentTime %s;", event.SentTime), common.LogErrorlvl)
		}
	}
	return nil
}
