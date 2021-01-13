package pubsub

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/srv/scheduler/proto/scheduler"
	"time"
)

func SendControl(msg, tag string) {
	ev := &scheduler.Evt{
		Tag:      tag,
		SentTime: time.Now().Format(common.TimestampLayout),
		Msg:      msg,
	}

	// publish an event
	if err := pub.Publish(context.Background(), ev); err != nil {
		log.Error(fmt.Sprintf("error publishing: %v", err))
	}
	log.Info(fmt.Sprintf("sent %+v", ev))

}
