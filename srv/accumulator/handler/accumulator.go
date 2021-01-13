package handler

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	"github.com/sidazhang123/f10-go/srv/accumulator/model"
	"github.com/sidazhang123/f10-go/srv/accumulator/proto/accumulator"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Accumulator struct{}

var (
	log                = z.GetLogger()
	accumulatorService *model.Service
	colToAggMethod     = map[string]string{"latest_tips": "DoLT", "financial_analysis": "DoFA", "operational_analysis": "DoOA", "shareholder_analysis": "DoSA"}
)

func Init() { accumulatorService, _ = model.GetService() }

// given specific column names or run all of them
func (e *Accumulator) Agg(_ context.Context, req *accumulator.Request, _ *accumulator.Response) error {
	log.Info(fmt.Sprintf("Received Accumulator.Agg request %s; %s", req.GetCollection(), time.Now().Format(common.TimestampLayout)))

	cols := strings.Split(req.GetCollection(), ",")
	runAll := true
	var wg sync.WaitGroup
	for _, col := range cols {
		if _, ok := colToAggMethod[col]; ok {
			runAll = false
			method := reflect.ValueOf(accumulatorService).MethodByName(colToAggMethod[col])
			if !method.IsValid() {
				return fmt.Errorf("[handler] method %s invalid in AccumulatorService", colToAggMethod[col])
			}
			wg.Add(1)
			go func() {
				params := []reflect.Value{reflect.ValueOf(req.GetCode()), reflect.ValueOf(req.GetStart()), reflect.ValueOf(req.GetEnd())}
				err := method.Call(params)[0].Interface()
				if err != nil {
					log.Error(fmt.Sprintf("[handler call %s] %s", colToAggMethod[col], err.(error).Error()))
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
	if runAll {
		errList := accumulatorService.DoAll(req.GetCode(), req.GetStart(), req.GetEnd(), &wg)
		if len(errList) > 0 {
			for _, err := range errList {
				log.Error(err.Error())
			}
		}
	}

	log.Info(fmt.Sprintf("Received Accumulator.Agg done %s; %s", req.GetCollection(), time.Now().Format(common.TimestampLayout)))
	return nil
}

func (e *Accumulator) ReprAll(_ context.Context, req *accumulator.Request, _ *accumulator.Response) error {
	log.Info("Received Accumulator.ReprAll request; " + time.Now().Format(common.TimestampLayout))

	err := accumulatorService.ReprAll()
	if err != nil && len(err) > 0 {
		eMsg := ""
		for _, e := range err {
			eMsg += e.Error() + "\n"
		}
		log.Error(eMsg)
	}
	log.Info("Accumulator.ReprAll done; " + time.Now().Format(common.TimestampLayout))
	return nil
}

func (e *Accumulator) GetRepr(_ context.Context, req *accumulator.ReprReq, rsp *accumulator.Response) error {
	log.Info("Received Accumulator.GetRepr request; " + time.Now().Format(common.TimestampLayout))

	s := accumulatorService.GetRepr(req.GetWin())
	rsp.Success = true
	rsp.Msg = s

	return nil
}
