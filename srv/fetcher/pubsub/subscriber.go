package pubsub

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/srv/scheduler/proto/scheduler"
)

type Sub struct{}

func (s *Sub) ToFetch(ctx context.Context, event *scheduler.Evt) error {
	if event.Tag == common.FetchStart {
		if available {
			available = false
			sentTime := event.SentTime
			msg := fmt.Sprintf("Fetch broadcast received. sentTime %s; msg %s", sentTime, event.Msg)
			log.Info(msg)
			SendLog(msg, common.LogInfolvl)

			err := fetcherService.GetByMQ()
			if len(err) > 0 {
				for _, e := range err {
					log.Error(e.Error())
					SendLog(e.Error(), common.LogErrorlvl)
				}
				log.Info(common.FetchComp + " with errors")
				SendLog(common.FetchComp+" with errors", common.LogCallbacklvl)
				available = true
				return nil
			}
			log.Info(common.FetchComp)
			SendLog(common.FetchComp, common.LogCallbacklvl)
			available = true
		}
	}
	return nil
}
