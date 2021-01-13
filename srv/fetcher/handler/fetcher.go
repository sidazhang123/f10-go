package handler

import (
	"context"
	"fmt"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	"github.com/sidazhang123/f10-go/srv/fetcher/model"
	"github.com/sidazhang123/f10-go/srv/fetcher/proto/fetcher"
)

type Fetcher struct{}

var (
	log          = z.GetLogger()
	fetchService *model.Service
)

func Init() { fetchService, _ = model.GetService() }

// Call is a single request handler called via client.Call or the generated client code
func (e *Fetcher) FetchRaw(ctx context.Context, req *fetcher.Request, rsp *fetcher.Response) error {
	log.Info("Received Fetcher.FetchRaw request")
	err := fetchService.GetByMQ()
	if err != nil {
		rsp.Success = false
		rsp.Error.Code = 500
		eStr := ""
		for _, e := range err {
			eStr += e.Error()

		}
		rsp.Error.Detail = eStr
		log.Error("[handler] FetchRaw err\n" + eStr)
		return fmt.Errorf(eStr)
	}
	log.Info("FetchRaw done")
	rsp.Success = true
	return nil
}

func (e *Fetcher) QueryRaw(ctx context.Context, req *fetcher.Request, rsp *fetcher.Response) error {
	log.Info("Received Fetcher.QueryRaw request")

	return nil
}
