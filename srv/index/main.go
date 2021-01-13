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
	"github.com/sidazhang123/f10-go/srv/index/handler"
	"github.com/sidazhang123/f10-go/srv/index/model"
	proto "github.com/sidazhang123/f10-go/srv/index/proto/index"
	"github.com/sidazhang123/f10-go/srv/index/pubsub"
	"go.uber.org/zap"

	"os"
	"time"
)

var (
	log     = z.GetLogger()
	cfg     = &indexCfg{}
	appName = "index_srv"
)

type indexCfg struct {
	common.AppCfg
}

func main() {
	initCfg()
	reg := etcd.NewRegistry(registryOptions)
	// New Service
	service := micro.NewService(
		micro.Name(cfg.Name),
		micro.Version(fmt.Sprintf("%f", cfg.Version)),
		micro.Registry(reg),
		micro.Address(cfg.Addr()),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*10),
	)

	// Initialise service
	service.Init(
		micro.Action(func(ctx *cli.Context) error {
			model.Init()
			handler.Init()
			pubsub.Init(micro.NewEvent(common.LoggingTopic, service.Client()))
			return nil
		}))

	// Register Handler
	_ = proto.RegisterIndexHandler(service.Server(), new(handler.Index))

	//Register Subscriber
	err := micro.RegisterSubscriber(common.ControlTopic, service.Server(), new(pubsub.Sub))
	if err != nil {
		println(err.Error())
	}
	// Run service
	if err := service.Run(); err != nil {
		log.Error("[Main] Failed to start service.\n%v", zap.Any("err", err))
		os.Exit(1)
	}
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
func initCfg() {
	// ENV_ADDR
	source := grpc.NewSource(grpc.WithAddress(os.Getenv("CONFIG_GRPC_ADDR")), grpc.WithPath("env"))
	//source := grpc.NewSource(grpc.WithAddress("127.0.0.1:9600"), grpc.WithPath("env"))
	basic.Init(config.WithSource(source),
		config.WithApp(appName))
	err := config.C().App(appName, cfg)
	if err != nil {
		panic(err)
	}
	log.Info("[InitCfg] Conf Loaded...\n%+v", zap.Any("cfg", cfg))

}
