package main

import (
	"fmt"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/config/source/grpc/v2"
	"github.com/sidazhang123/f10-go/basic"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/basic/config"
	"github.com/sidazhang123/f10-go/web/mgmt/auth_crypto"
	_ "github.com/sidazhang123/f10-go/web/mgmt/auth_crypto"
	"github.com/sidazhang123/f10-go/web/mgmt/handler"
	"net/http"
	"os"
	"time"
)

var (
	cfg     = &mgmtCfg{}
	appName = "mgmt_web"
)

type mgmtCfg struct {
	common.AppCfg
}

func main() {
	initCfg()
	reg := etcd.NewRegistry(registryOptions)
	// create new web service
	service := web.NewService(
		web.Name(cfg.Name),
		web.Version(fmt.Sprintf("%f", cfg.Version)),
		web.Registry(reg),
		web.Address(cfg.Addr()),
		//web.Address(cfg.Addr()),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	authDebug := false
	weak := true
	// register html handler
	service.HandleFunc("/", auth_crypto.AuthWrapper(func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir("html")).ServeHTTP(w, r)
	}, authDebug, weak))
	webRoute := map[string]func(http.ResponseWriter, *http.Request){
		"/mgmt/load_options": handler.GetPluginPath,
		"/mgmt/select":       handler.GetSrc,
		"/mgmt/update":       handler.Update,
		"/mgmt/test":         handler.Test,
	}
	for r, f := range webRoute {
		service.HandleFunc(r, auth_crypto.AuthWrapper(f, authDebug, weak))
	}

	// register feed handler
	weak = false
	apiRoute := map[string]func(http.ResponseWriter, *http.Request){
		"/feed/focus/get":      handler.GetFocus,
		"/feed/focus/purge":    handler.PurgeFocus,
		"/feed/focus/delete":   handler.ToggleFocusDel,
		"/feed/focus/fav":      handler.ToggleFocusFav,
		"/feed/focus/generate": handler.GenFocus,
		"/feed/rules/create":   handler.CreateRules,
		"/feed/rules/get":      handler.ReadRules,
		"/feed/rules/update":   handler.UpdateRules,
		"/feed/rules/delete":   handler.DeleteRules,
		"/feed/jpush_reg":      handler.AddRegId,
		"/feed/focus/stat":     handler.GetFocusStat,
		"/feed/chan_o_day/set": handler.SetChanODay,
		"/feed/chan_o_day/get": handler.GetChanODay,
		"/feed/log":            handler.Log,
	}
	for r, f := range apiRoute {
		service.HandleFunc(r, auth_crypto.AuthWrapper(f, authDebug, weak))
	}
	// register apk handler
	weak = false
	ec := handler.Monitor()
	go func() {
		for {
			e := <-ec
			if e != nil {
				log.Error(e.Error())
			}
		}
	}()
	service.HandleFunc("/apk/info", auth_crypto.AuthWrapper(handler.GetApkInfo, authDebug, weak))
	service.HandleFunc("/apk/download", handler.DownloadApk)

	service.HandleFunc("/repr", auth_crypto.AuthWrapper(handler.GetRepr, authDebug, weak))
	service.HandleFunc("/f10bin", auth_crypto.AuthWrapper(handler.GetBin, authDebug, weak))
	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
func initCfg() {
	//source := grpc.NewSource(grpc.WithPath("env"), grpc.WithAddress(os.Getenv("CONFIG_GRPC_ADDR")))
	source := grpc.NewSource(grpc.WithPath("env"), grpc.WithAddress("127.0.0.1:9600"))
	basic.Init(config.WithSource(source), config.WithApp(appName))
	err := config.C().App(appName, cfg)

	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	auth_crypto.Init()

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
