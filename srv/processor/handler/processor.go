package handler

import (
	"context"
	"encoding/json"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	"github.com/sidazhang123/f10-go/srv/processor/model"
	"github.com/sidazhang123/f10-go/srv/processor/proto/processor"
	"go.uber.org/zap"
)

var (
	log              = z.GetLogger()
	processorService *model.Service
)

type Processor struct{}

func Init() { processorService, _ = model.GetService() }

/*
	@Params pluginName, testStr
	@Return resStr, errMsg
Given a test string, test if a .so works
*/
func (e *Processor) RegexTest(ctx context.Context, req *processor.RegexReq, rsp *processor.RegexRsp) error {
	log.Info("Received Processor.RegexTest request")
	resMap, err := processorService.RegexTest(req.TestStr, req.PluginPath)
	if err != nil {
		rsp.ErrMsg = err.Error()
		log.Error("[handler] RegexTest", zap.Any("err", err.Error()))
		return err
	}
	jsonString, err := json.Marshal(resMap)
	if err != nil {
		rsp.ErrMsg = err.Error()
		log.Error("[handler] RegexTest", zap.Any("err", err.Error()))
		return err
	}
	rsp.ResStr = string(jsonString)
	return nil
}

/*
	@Params PluginPath, sourceCode
	@Return ErrMsg, path
GetSourceCode first to get the current code of a plugin
*/
func (e *Processor) BuildSo(ctx context.Context, req *processor.BuildSoReq, rsp *processor.BuildSoRsp) error {
	log.Info("Received Processor.BuildSo request")
	soPath, err := processorService.BuildSo(req.SourceCode, req.PluginPath)
	if err != nil {
		rsp.ErrMsg = err.Error()
		log.Error("[handler] BuildSo", zap.Any("err", err.Error()))
		return err
	}
	rsp.Path = soPath
	return nil
}

/*
	@Params pluginName, path
	@Return errMsg, sourceCode, plugins(filename with path prefix in that directory)
*/
func (e *Processor) GetSourceCode(ctx context.Context, req *processor.GetSourceCodeReq, rsp *processor.GetSourceCodeRsp) error {
	log.Info("Received Processor.GetSourceCode request")
	src, err := processorService.GetPluginSrc(req.Path)
	if err != nil {
		rsp.ErrMsg = err.Error()
		log.Error("[handler] GetSourceCode", zap.Any("err", err.Error()))
		return err
	}
	rsp.SourceCode = src
	return nil
}

/*
	@Params date,flagName,params
	@Return success, stockList, error
*/
func (e *Processor) Process(ctx context.Context, req *processor.Request, rsp *processor.Response) error {
	log.Info("Received Processor.Process request")
	err := processorService.Refine(req)
	if err != nil {
		rsp.Success = false
		rsp.Error = &processor.Error{Detail: err.Error()}
		log.Error("[handler] Process", zap.Any("err", err.Error()))
		log.Info("Processor.Process done")
		return err
	}
	log.Info("Processor.Process done")
	rsp.Success = true
	return nil
}

func (e *Processor) GetPluginPath(ctx context.Context, req *processor.GetPluginPathReq, rsp *processor.GetPluginPathRsp) error {
	log.Info("Received Processor.GetPluginPath request")
	path, err := processorService.GetPluginPath()
	if err != nil {
		rsp.ErrMsg = err.Error()
		log.Error("[handler] GetPluginPath", zap.Any("err", err.Error()))
		return err
	}
	rsp.JoinedPath = path
	log.Info("path=" + path)
	return nil
}
