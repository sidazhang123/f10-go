package db

import (
	"context"
	"fmt"

	"github.com/sidazhang123/f10-go/basic/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type db struct {
	Mongodb Mongodb `json:"mongodb"`
}
type Mongodb struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	AuthDB   string `json:"auth_db"`
}

func initMongodb() {
	log.Info("[InitMongodb] Initializing Mongodb...")
	c := config.C()
	cfg := &db{}
	err := c.App("db", cfg)
	if err != nil {
		log.Error(fmt.Sprintf("[InitMongodb] Failed to get config\n%s", err))
		return
	}
	log.Info(fmt.Sprintf("%+v", cfg))
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.Mongodb.Username, cfg.Mongodb.Password, cfg.Mongodb.Host, cfg.Mongodb.Port, cfg.Mongodb.AuthDB))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error(fmt.Sprintf("[InitMongodb] Failed to connect\n%s", err))
		return
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Error(fmt.Sprintf("[InitMongodb] Ping failed\n%s", err))
		return
	}
	mongoClient = client
	log.Info("[InitMongodb] initialized successfully")
}

func CloseDb() error {
	err := mongoClient.Disconnect(context.TODO())

	if err != nil {
		return err
	}
	return nil

}
