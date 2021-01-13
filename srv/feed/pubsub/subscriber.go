package pubsub

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/srv/scheduler/proto/scheduler"
	"strconv"
	"time"
)

type Sub struct {
}

func (s *Sub) ToGenFeed(ctx context.Context, event *scheduler.Evt) error {
	if event.Tag == common.GenFeedStart {
		sentTime := event.SentTime
		msg := fmt.Sprintf("ToGenFeed broadcast received. sentTime %s; msg %s", sentTime, event.Msg)
		log.Info(msg)
		SendLog(msg, common.LogInfolvl)
		//append mode, don't purge
		//_, _ = feedService.PurgeFocus("", 0)
		err := feedService.GenerateFocus("")
		if err != nil {
			log.Error(err.Error())
			SendLog(err.Error(), common.LogErrorlvl)
		}
		err = feedService.GenerateOperationalAnalysisDiffCSV("")
		if err != nil {
			log.Error(err.Error())
			SendLog(err.Error(), common.LogErrorlvl)
		}
		log.Info(common.GenFeedComp)
		SendLog(common.GenFeedComp, common.LogCallbacklvl)

	}
	return nil
}

func (s *Sub) DeleteOutdatedFocus(ctx context.Context, event *scheduler.Evt) error {
	if event.Tag == common.DeleteOutdatedFocus {
		sentTime := event.SentTime
		msg := fmt.Sprintf("DeleteOutdatedFocus broadcast received. sentTime %s; msg %s", sentTime, event.Msg)
		log.Info(msg)
		SendLog(msg, common.LogInfolvl)
		// big wipe
		err, num := feedService.DeleteOutdatedFocus()
		if err != nil {
			log.Error(err.Error())
			SendLog(err.Error(), common.LogErrorlvl)
		}
		msg = fmt.Sprintf("DeleteOutdatedFocus %d deleted", num)
		log.Info(msg)
		SendLog(msg, common.LogInfolvl)
		// purge the recycle bin
		todayTs := time.Now().UTC().Add(8 * time.Hour).Format(common.TimestampLayout[:10])[8:10]
		t, _ := strconv.Atoi(todayTs)
		if t%2 == 0 {
			err, num = feedService.DeleteRecoveryBin()
			if err != nil {
				log.Error(err.Error())
				SendLog(err.Error(), common.LogErrorlvl)
			}
			msg = fmt.Sprintf("DeleteRecoveryBin %d deleted", num)
			log.Info(msg)
			SendLog(msg, common.LogInfolvl)
		}

	}
	return nil
}
