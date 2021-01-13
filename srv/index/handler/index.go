package handler

import (
	"context"
	"fmt"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	"github.com/sidazhang123/f10-go/srv/index/model"
	"github.com/sidazhang123/f10-go/srv/index/proto/index"
	"sync"
	"time"
)

type Index struct{}

var (
	log          = z.GetLogger()
	fetchService *model.Service
	lastRun      time.Time
	lock         sync.RWMutex
)

const interval = 5 * time.Minute

func Init() {
	fetchService, _ = model.GetService()
	lastRun = time.Now().Add(-interval)
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Index) GetCodeName(ctx context.Context, req *index.Request, rsp *index.Response) error {
	log.Info("Received Index.GetCodeName request")
	lock.Lock()
	// can't be called twice in 5 min
	if !time.Now().After(lastRun.Add(interval)) {
		return fmt.Errorf("Index.GetCodeName is running!")
	}
	lastRun = time.Now()
	lock.Unlock()

	err, res := fetchService.Fetch()
	if err != nil {
		rsp.Success = false
		rsp.Error.Detail = err.Error()
		log.Error(fmt.Sprintf("[handler] GetCodeName err \n%s\n though %d codes fetched", err, len(res)))
		return err
	}
	log.Info(fmt.Sprintf("[handler]  %d codes fetched", len(res)))
	rsp.Success = true
	rsp.StockList = res

	return nil
}
