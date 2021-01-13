package config

import (
	"fmt"
	"github.com/micro/go-micro/v2/config"

	"github.com/micro/go-micro/v2/util/log"

	"sync"
)

var (
	m      sync.RWMutex
	inited bool
	c      = &configurator{}
)

type configurator struct {
	conf    config.Config
	appName string
}
type Configurator interface {
	App(name string, config interface{}) (err error)
	Path(path string, config interface{}) (err error)
	AppPath(app string, path string, config interface{}) (err error)
}

func (c *configurator) App(name string, config interface{}) (err error) {
	v := c.conf.Get(name)
	if v != nil {
		err = v.Scan(&config)
	} else {
		err = fmt.Errorf("[config-App] Given config does not exist")
	}
	return
}
func (c *configurator) Path(path string, config interface{}) (err error) {
	v := c.conf.Get(c.appName, path)
	if v != nil {
		err = v.Scan(&config)
	} else {
		err = fmt.Errorf("[config-App] Given config does not exist")
	}
	return
}
func (c *configurator) AppPath(app string, path string, config interface{}) (err error) {
	v := c.conf.Get(app, path)
	if v != nil {
		err = v.Scan(&config)
	} else {
		err = fmt.Errorf("[config-App] Given config does not exist")
	}
	return
}
func C() Configurator {
	return c
}

func (c *configurator) init(ops Options) (err error) {
	m.Lock()
	defer m.Unlock()
	if inited {
		log.Logf("[Init] Configurations were loaded...skipping")
		return
	}
	c.conf, err = config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	c.appName = ops.AppName
	err = c.conf.Load(ops.Sources...)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Logf("[Init] Listening to conf changes...")
		watcher, err := c.conf.Watch()
		if err != nil {
			log.Fatal(err)
		}
		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatal(err)
			}
			log.Logf("[Init]Conf Changes:\n%v", string(v.Bytes()))

		}
	}()

	inited = true
	log.Logf("[Init]Config inited ")
	return
}

func Init(opts ...Option) {
	ops := &Options{}
	for _, o := range opts {
		o(ops)
	}
	_ = c.init(*ops)
}
