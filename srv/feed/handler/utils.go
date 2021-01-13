package handler

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/srv/feed/proto/feed"
)

func (e *Feed) AddJPushReg(ctx context.Context, req *feed.JPushReg, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.ReadFocus request")
	id := req.GetId()
	emptyErr := fmt.Errorf("empty JPush_reg_id")
	if id == "" {
		rsp.Success = false
		rsp.Msg = emptyErr.Error()
		return emptyErr
	}
	err := feedService.AddJPushReg(id)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
	} else {
		rsp.Success = true
	}
	return err
}

func (e *Feed) Log(ctx context.Context, req *feed.PlainReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.Log request")

	err := feedService.Log(req.Msg)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
	} else {
		rsp.Success = true
	}
	return err
}
