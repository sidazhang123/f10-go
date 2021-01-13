package pubsub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getAccessToken() (error, string) {
	key := Params.AppKey
	secret := Params.AppSecret
	resp, _ := http.Get(fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", key, secret))
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var res map[string]interface{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		return err, ""
	}

	if v, ok := res["access_token"]; ok {
		if v != nil {
			return nil, v.(string)
		}

	}
	return fmt.Errorf(fmt.Sprintf("%+v", res)), ""
}

func SendAlarm(msg string) error {
	agentId := string(Params.AgentId)
	deptId := string(Params.DeptId)
	jsonStr := []byte(fmt.Sprintf(`{
	"agent_id":%s,
	"msg":{
		"msgtype":"text",
		"text":{
			"content":"%s"
		}
	},
	"dept_id_list":"%s"
	}`, agentId, msg, deptId))
	err, accessToken := getAccessToken()
	if err != nil {
		return err
	}
	url := "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=" + accessToken
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var rsp map[string]interface{}
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return err
	}
	if v, ok := rsp["errmsg"]; ok {
		if v != nil && v.(string) == "ok" {
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("errmsg: %+v", v))
	}
	return fmt.Errorf("err ding rsp:" + string(body))
}
