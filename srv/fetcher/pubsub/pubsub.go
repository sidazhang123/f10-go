package pubsub

import (
	"github.com/micro/go-micro/v2"
	"github.com/sidazhang123/f10-go/plugins/zap"
	fetcher "github.com/sidazhang123/f10-go/srv/fetcher/model"
)

var (
	fetcherService *fetcher.Service
	pub            micro.Publisher
	log            = zap.GetLogger()
	available      = true
)

func Init(publisher micro.Publisher) {
	fetcherService, _ = fetcher.GetService()
	pub = publisher

}
