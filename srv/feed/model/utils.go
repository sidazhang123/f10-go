package model

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func (s Service) Log(msg string) error {
	return InsertOneFocus("log", map[string]interface{}{"err": msg, "timestamp": time.Now().UTC().Add(8 * time.Hour).Format(common.TimestampLayout)})
}

func (s Service) AddJPushReg(id string) error {
	return InsertJPushID("jpush_reg", id)
}
func (s Service) FindLatestFetchTime(collection string) (error, string) {
	return FindLatestFetchTime(collection)
}

func makeJPush(msg string) error {
	log.Info("makepush called")
	err, cur := FindAll("jpush_reg")
	if err != nil {
		log.Info("failed to visit jpush_reg, err=" + err.Error())
		return err
	}
	id := ""
	for cur.Next(context.TODO()) {
		var reg bson.M
		err := cur.Decode(&reg)
		if err != nil {
			log.Info("failed to decode reg_id, err=" + err.Error())
			return err
		}
		id += reg["reg_id"].(string) + "\",\""
	}
	id = strings.TrimSuffix(id, "\",\"")

	content := fmt.Sprintf("{\"platform\":\"all\",\"audience\":{\"registration_id\":[\"%s\"]},\"notification\" : {\"alert\" : \"%s\"}}", id, msg)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.jpush.cn/v3/push", bytes.NewBuffer([]byte(content)))
	if err != nil {
		log.Info("failed to make jpush req, err=" + err.Error())
		return err
	}
	req.SetBasicAuth(Params.JPush0, Params.JPush1)

	rsp, err := client.Do(req)
	if err != nil {
		log.Info("failed to push, err=" + err.Error())
		return err
	}
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Info("failed to read jpush rsp, err=" + err.Error())
		return err
	}
	if strings.Contains(string(b), "error") {
		log.Info("failed to push, err=" + string(b))
		return fmt.Errorf("failed to push, err=" + string(b))
	}

	return nil
}
