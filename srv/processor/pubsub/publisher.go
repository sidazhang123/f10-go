package pubsub

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/srv/scheduler/proto/scheduler"
	"time"
)

func SendLog(msg string, lvl int32) {
	ev := &scheduler.Log{
		Tag:      common.ProcessStart,
		Level:    lvl,
		SentTime: time.Now().Format(common.TimestampLayout),
		Msg:      msg,
	}

	// publish an event
	if err := pub.Publish(context.Background(), ev); err != nil {
		log.Error(fmt.Sprintf("error publishing: %v", err))
	}

}
