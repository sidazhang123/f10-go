package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
	"github.com/micro/go-micro/v2/util/log"
	proto "github.com/micro/go-plugins/config/source/grpc/proto"
	"google.golang.org/grpc"
	"net"
	"strings"
	"time"
)

var (
	apps = []string{"env", "cron"} // "srv","api","cron"
)

type Service struct {
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Logf("[Config-grpc] Recovered from\n%s", r)
		}
	}()

	err := loadAndWatchConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	service := grpc.NewServer()

	proto.RegisterSourceServer(service, new(Service))
	ts, err := net.Listen("tcp", ":9600")
	if err != nil {
		log.Fatal(err)
	}
	log.Logf("[Config-grpc] Config Server started.")
	err = service.Serve(ts)
	if err != nil {
		log.Fatal(err)
	}

}
func getConfig(appName string) *proto.ChangeSet {
	bytes := config.Get(appName).Bytes()
	//log.Logf("[getConfig] appName: %s", appName)
	//log.Logf("[getConfig] %+v", config.Get(appName).StringMap(map[string]string{}))
	return &proto.ChangeSet{
		Data:      bytes,
		Checksum:  fmt.Sprintf("%x", md5.Sum(bytes)),
		Format:    "yml",
		Source:    "file",
		Timestamp: time.Now().Unix(),
	}
}

func parsePath(path string) (appName string) {
	paths := strings.Split(path, "/")
	if paths[0] == "" && len(paths) > 1 {
		return paths[1]
	}
	return paths[0]
}

func (s Service) Read(ctx context.Context, req *proto.ReadRequest) (rsp *proto.ReadResponse, err error) {
	appName := parsePath(req.Path)
	g := getConfig(appName)
	rsp = &proto.ReadResponse{ChangeSet: g}
	return
}

func (s Service) Watch(req *proto.WatchRequest, ws proto.Source_WatchServer) (err error) {
	appName := parsePath(req.Path)
	rsp := &proto.WatchResponse{ChangeSet: getConfig(appName)}
	if err = ws.Send(rsp); err != nil {
		log.Logf("[Watch] Error Occurred.\n%s", err)
		return
	}
	return
}

func loadAndWatchConfigFile() (err error) {
	for _, app := range apps {
		if err = config.Load(file.NewSource(
			file.WithPath("./conf/" + app + ".yml"))); err != nil {
			log.Fatalf("[loadAndWatchConfigFile] Failed to load conf file:\n%s", err)
			return
		}
		log.Logf("[loadAndWatchConfigFile] %s loaded", app)
	}

	go func() {
		for {

			watcher, err := config.Watch()
			if err != nil {
				log.Fatalf("[loadAndWatchConfigFile] Can't start to watch on conf files:\n%s", err)
				return
			}

			v, err := watcher.Next()
			if err != nil {
				log.Fatalf("[loadAndWatchConfigFile] Failed to watch changes on conf files:\n%s", err)
				return
			}
			log.Logf("[loadAndWatchConfigFile] Conf changed: %v", string(v.Bytes()))
		}
	}()
	return
}
