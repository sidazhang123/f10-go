package pubsub

import (
	"github.com/micro/go-micro/v2"
	"github.com/sidazhang123/f10-go/plugins/zap"
	"github.com/sidazhang123/f10-go/srv/feed/model"
)

var (
	feedService *model.Service
	pub         micro.Event
	log         = zap.GetLogger()
)

func Init(publisher micro.Event) {
	feedService, _ = model.GetService()
	pub = publisher

}
