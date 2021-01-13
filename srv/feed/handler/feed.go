package handler

import (
	z "github.com/sidazhang123/f10-go/plugins/zap"
	"github.com/sidazhang123/f10-go/srv/feed/model"
)

var (
	log         = z.GetLogger()
	feedService *model.Service
)

func Init() { feedService, _ = model.GetService() }

type Feed struct{}
