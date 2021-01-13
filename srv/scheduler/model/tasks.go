package model

import (
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/srv/scheduler/pubsub"
)

func (s *Service) GetCodeName() {
	pubsub.SendControl("", common.IndexStart)
}

func (s *Service) DeleteOutdatedFocus() {
	pubsub.SendControl("", common.DeleteOutdatedFocus)
}
