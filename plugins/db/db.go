package db

import (
	"fmt"
	"github.com/sidazhang123/f10-go/basic"

	z "github.com/sidazhang123/f10-go/plugins/zap"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

var (
	inited      = false
	mongoClient *mongo.Client
	m           sync.RWMutex
	log         = z.GetLogger()
)

func init() {
	basic.Register(initDB)
}
func initDB() {
	m.Lock()
	defer m.Unlock()
	if inited {

		log.Warn(fmt.Errorf("[Init] DB was initialized").Error())
		return
	}

	initMongodb()

	inited = true
}
func GetDB() *mongo.Client {
	if mongoClient == nil || inited == false {
		initMongodb()
	}
	return mongoClient
}
