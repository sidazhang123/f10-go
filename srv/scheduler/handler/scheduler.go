package handler

import (
	"context"
	"github.com/sidazhang123/f10-go/srv/scheduler/model"
	"go.uber.org/zap"

	z "github.com/sidazhang123/f10-go/plugins/zap"
	proto "github.com/sidazhang123/f10-go/srv/scheduler/proto/scheduler"
)

type Scheduler struct{}

var (
	log         = z.GetLogger()
	cronService *model.Service
)

func Init() { cronService, _ = model.GetService() }

func (e *Scheduler) Once(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Info("[handler] Received Once request")

	err, name, t := cronService.Once(req.GetFuncName())
	if err != nil {
		rsp.Success = false
		rsp.Error.Detail = err.Error()
		log.Warn("[handler] Once err", zap.Any("err", err))
		return err
	}
	rsp.Success = true
	rsp.Task = &proto.NextSchedule{FuncName: name, ScheduledTime: t.Unix()}

	return nil
}

func (e *Scheduler) NextScheduledTask(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Info("[handler] Received NextScheduledTask request")
	name, t := cronService.NextScheduledTask()
	rsp.Task = &proto.NextSchedule{FuncName: name, ScheduledTime: t.Unix()}
	return nil
}

func (e *Scheduler) UpdateSchedule(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Info("[handler] Received UpdateSchedule request")

	err, name, t := cronService.SetCron()
	if err != nil {
		rsp.Success = false
		rsp.Error.Detail = err.Error()
		log.Info("[handler] SetCron err", zap.Any("err", err))
		return err
	}
	rsp.Success = true
	rsp.Task = &proto.NextSchedule{FuncName: name, ScheduledTime: t.Unix()}
	return nil
}
func (e *Scheduler) ClearSchedule(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Info("[handler] Received ClearSchedule request")
	name, t := cronService.ClearCron()
	rsp.Success = true
	rsp.Task = &proto.NextSchedule{FuncName: name, ScheduledTime: t.Unix()}
	return nil
}
func (e *Scheduler) RemoveTask(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Info("[handler] Received RemoveTask request")
	err, name, t := cronService.RemoveTask(req.GetFuncName())
	if err != nil {
		rsp.Success = false
		rsp.Error = &proto.Error{Code: 500, Detail: err.Error()}
		log.Info("[handler] RemoveTask err", zap.Any("err", err))
		return err
	}
	rsp.Success = true
	rsp.Task = &proto.NextSchedule{FuncName: name, ScheduledTime: t.Unix()}
	return nil
}
func (e *Scheduler) DingAlarm(ctx context.Context, req *proto.Request, rsp *proto.Error) error {
	log.Info("[handler] Ding")
	err := cronService.SendAlarm(req.GetFuncName())
	if err != nil {
		rsp.Detail = err.Error()
	}
	return nil
}
