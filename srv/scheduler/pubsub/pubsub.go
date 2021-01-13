package pubsub

import (
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/sidazhang123/f10-go/basic/config"
	"github.com/sidazhang123/f10-go/plugins/zap"
)

var (
	pub    micro.Event
	log    = zap.GetLogger()
	Params = &Opts{}
)

func Init(publisher micro.Publisher) {
	pub = publisher
	err := config.C().Path("params", Params)
	if err != nil {
		err = fmt.Errorf("[Init] Failed to get params\n" + err.Error())
		log.Error(err.Error())
		return
	}
	log.Info(fmt.Sprint("[Init] pubsub"))
	log.Info(fmt.Sprintf("%+v", Params))
}
