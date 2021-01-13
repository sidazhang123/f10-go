package main

import (
	"fmt"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/config/source/grpc/v2"
	"github.com/sidazhang123/f10-go/basic"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/basic/config"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	"github.com/sidazhang123/f10-go/srv/scheduler/handler"
	"github.com/sidazhang123/f10-go/srv/scheduler/model"
	proto "github.com/sidazhang123/f10-go/srv/scheduler/proto/scheduler"
	"github.com/sidazhang123/f10-go/srv/scheduler/pubsub"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	cfg     = &cronCfg{}
	appName = "scheduler_srv"
	log     = z.GetLogger()
)

type cronCfg struct {
	common.AppCfg
}

func main() {

	initCfg()
	reg := etcd.NewRegistry(registryOptions)
	// New SchedulerService

	service := micro.NewService(
		micro.Name(cfg.Name),
		micro.Version(fmt.Sprintf("%f", cfg.Version)),
		micro.Address(cfg.Addr()),

		micro.Registry(reg),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
	)

	// Initialise service
	service.Init(
		micro.Action(func(context *cli.Context) error {
			model.Init()
			handler.Init()
			pubsub.Init(micro.NewEvent(common.ControlTopic, service.Client()))
			return nil
		}))

	// Register Handler
	_ = proto.RegisterSchedulerHandler(service.Server(), new(handler.Scheduler))

	//Register Subscriber
	_ = micro.RegisterSubscriber(common.LoggingTopic, service.Server(), new(pubsub.Sub))

	// Run service
	if err := service.Run(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
func initCfg() {
	source1 := grpc.NewSource(grpc.WithPath("env"), grpc.WithAddress(os.Getenv("CONFIG_GRPC_ADDR")))
	source2 := grpc.NewSource(grpc.WithPath("cron"), grpc.WithAddress(os.Getenv("CONFIG_GRPC_ADDR")))
	basic.Init(config.WithSource(source1), config.WithSource(source2), config.WithApp(appName))
	err := config.C().App(appName, cfg)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	log.Info("[initCfg] Conf Loaded...\n%+v", zap.Any("cfg", cfg))
}
func registryOptions(ops *registry.Options) {
	etcdConfig := &common.EtcdConfig{}
	err := config.C().App("etcd", etcdConfig)
	if err != nil {
		panic(err)
	}
	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdConfig.Host, etcdConfig.Port)}
	ops.Timeout = 5 * time.Second
}
