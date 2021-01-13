package handler

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/srv/feed/model"
	"github.com/sidazhang123/f10-go/srv/feed/proto/feed"
	"strconv"
)

func (e *Feed) GetFocusStat(ctx context.Context, req *feed.RuleReq, rsp *feed.Chans) error {
	log.Info("Received Feed.GetFocusStat request")
	err, stat := feedService.GetFocusStat()
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()

	} else {
		rsp.Success = true
		var chans []*feed.Chan
		for k, v := range stat {
			chans = append(chans, &feed.Chan{
				ChanName: k,
				NoMsg:    v,
			})
		}
		rsp.Chans = chans
	}
	return nil
}

func (e *Feed) GenerateFocus(ctx context.Context, req *feed.ManipulateFocusReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.GenerateFocus request")
	// append mode, don't purge
	//err, _ := feedService.PurgeFocus("", 0)
	//if err != nil {
	//	rsp.Success = false
	//	rsp.Str = err.Error()
	//	return err
	//}
	err := feedService.GenerateFocus("")
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
		return err
	}
	err = feedService.GenerateOperationalAnalysisDiffCSV("")
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
		return err
	}
	return nil
}

func (e *Feed) GenerateOperationalAnalysisDiffCSV(ctx context.Context, req *feed.PlainReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.GenerateOperationalAnalysisDiffCSV request")

	err := feedService.GenerateOperationalAnalysisDiffCSV(req.Msg)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
	} else {
		rsp.Success = true
	}
	log.Info(" Feed.GenerateOperationalAnalysisDiffCSV DONE")
	return err
}

func (e *Feed) PurgeFocus(ctx context.Context, req *feed.ManipulateFocusReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.PurgeFocus request")
	err, count := feedService.PurgeFocus(req.Date)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
	} else {
		rsp.Success = true
		rsp.Msg = "purged " + strconv.Itoa(count)
	}
	return err
}

func (e *Feed) ReadFocus(ctx context.Context, req *feed.ManipulateFocusReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.ReadFocus request")
	chanId := ""
	if req.Chan != nil {
		chanId = req.Chan.Id
	}

	err, json := feedService.ReadFocus(req.Date, chanId, req.Del, req.Fav)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
	} else {
		rsp.Success = true
		rsp.Msg = json
	}
	return err
}

func (e *Feed) ToggleFocusDel(ctx context.Context, req *feed.ManipulateFocusReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.ToggleFocusDel request")
	err := feedService.ToggleFocusDel(req.ObjectId, req.Del)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
	} else {
		rsp.Success = true
	}
	return err
}

func (e *Feed) ToggleFocusFav(ctx context.Context, req *feed.ManipulateFocusReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.ToggleFocusFav request")
	err := feedService.ToggleFocusFav(req.ObjectId, req.Fav)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
	} else {
		rsp.Success = true
	}
	return err
}

func (e *Feed) DeleteOutdatedFocus(ctx context.Context, req *feed.PlainReq, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.DeleteOutdatedFocus request")
	err, deleted := feedService.DeleteOutdatedFocus()
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
	} else {
		rsp.Success = false
		rsp.Msg = strconv.Itoa(deleted)
	}
	return nil
}
func (e *Feed) GetChanODay(ctx context.Context, req *feed.PlainReq, rsp *feed.Chans) error {
	log.Info("Received Feed.GetChanODay request")
	err, chanDayMap := feedService.GetChanODay()
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
	} else {
		rsp.Success = true
		var chans []*feed.Chan
		for k, v := range chanDayMap {
			chans = append(chans, &feed.Chan{
				ChanName: k,
				NoMsg:    v,
			})
		}
		rsp.Chans = chans
	}
	return nil
}
func (e *Feed) SetChanODay(ctx context.Context, req *feed.Chans, rsp *feed.PlainRsp) error {
	log.Info("Received Feed.SetChanODay request")
	var toInsert []interface{}
	for _, c := range req.GetChans() {
		toInsert = append(toInsert, model.DaysToDel{
			Channel: c.GetChanName(),
			Day:     int(c.NoMsg),
		})
	}
	err, modified := feedService.SetChanODay(toInsert)
	if err != nil {
		rsp.Success = false
		rsp.Msg = err.Error()
	} else {
		rsp.Success = true
		rsp.Msg = fmt.Sprintf("set %d", modified)
	}
	return nil
}
