package handler

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/srv/feed/proto/feed"
	"strconv"
)

func (e *Feed) CreateRule(ctx context.Context, req *feed.RuleReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.CreateRule request")
	err, insertedN := feedService.CreateRule(req.Rules)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()

	} else {
		rsp.Success = true
		rsp.Msg = "inserted " + strconv.Itoa(insertedN)
	}
	return err
}

func (e *Feed) ReadRule(ctx context.Context, req *feed.RuleReq, rsp *feed.RuleRsp) error {
	log.Info("Received Feed.ReadRule request")
	err, rulesDTO := feedService.ReadRule(req.Rules)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()

	} else {
		rsp.Success = true
		var rules []*feed.Rule
		for _, r := range rulesDTO {
			rules = append(rules, &feed.Rule{
				TarCol:  r.TarCol,
				Channel: r.Channel,
				Key:     r.Key,
				Cond1:   r.Cond1,
				Cond2:   r.Cond2,
				Id:      r.ID.Hex(),
			})
		}
		rsp.Rules = rules

	}

	return err
}

func (e *Feed) UpdateRule(ctx context.Context, req *feed.RuleReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.UpdateRule request")
	err, modifiedN := feedService.UpdateRule(req.Rules)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()

	} else {
		rsp.Success = true
		rsp.Msg = "modified " + strconv.Itoa(modifiedN)
	}
	return err
}

func (e *Feed) DeleteRule(ctx context.Context, req *feed.RuleReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.DeleteRule request")
	err, deletedN := feedService.DeleteRule(req.Rules)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()

	} else {
		rsp.Success = true
		delRNum := deletedN % 10000
		delFNum := deletedN / 10000
		rsp.Msg = fmt.Sprintf("删%d规则", delRNum)
		if delFNum != 0 {
			rsp.Msg += fmt.Sprintf("\n及小于%d推送", delFNum)
		}
	}
	return err
}
