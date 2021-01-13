package pubsub

import (
	"github.com/micro/go-micro/v2"
	"github.com/sidazhang123/f10-go/plugins/zap"
	"github.com/sidazhang123/f10-go/srv/accumulator/model"
)

var (
	accumulatorService *model.Service
	pub                micro.Event
	log                = zap.GetLogger()
)

func Init(publisher micro.Event) {
	accumulatorService, _ = model.GetService()
	pub = publisher

}
