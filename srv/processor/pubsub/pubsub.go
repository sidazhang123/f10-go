package pubsub

import (
	"github.com/micro/go-micro/v2"
	"github.com/sidazhang123/f10-go/plugins/zap"
	processor "github.com/sidazhang123/f10-go/srv/processor/model"
)

var (
	processService *processor.Service
	pub            micro.Publisher
	log            = zap.GetLogger()
	available      = true
)

func Init(publisher micro.Publisher) {
	processService, _ = processor.GetService()
	pub = publisher

}
