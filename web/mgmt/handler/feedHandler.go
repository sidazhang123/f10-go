package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/util/log"
	proto "github.com/sidazhang123/f10-go/srv/feed/proto/feed"
	"github.com/sidazhang123/f10-go/web/mgmt/auth_crypto"
	"io/ioutil"
	"net/http"
)

var feedClient = proto.NewFeedService("sidazhang123.f10.srv.feed", client.DefaultClient)

//@Params date - by gendate
func GetFocus(w http.ResponseWriter, r *http.Request) {
	log.Info("GetFocus called")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	msg := auth_crypto.Decrypt(string(bodyBytes))
	var request proto.ManipulateFocusReq
	if err := json.Unmarshal([]byte(msg), &request); err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}

	rsp, err := feedClient.ReadFocus(context.TODO(), &request)

	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))

}

//@Params date - by gendate
//@Params push - $gte push
func PurgeFocus(w http.ResponseWriter, r *http.Request) {
	log.Info("PurgeFocus called")
	// call the backend service
	var request proto.ManipulateFocusReq

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}

	rsp, err := feedClient.PurgeFocus(context.TODO(), &request)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}

//@Params objectId, del
func ToggleFocusDel(w http.ResponseWriter, r *http.Request) {
	log.Info("ToggleFocusDel called")
	// call the backend service
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	msg := auth_crypto.Decrypt(string(bodyBytes))
	var request proto.ManipulateFocusReq
	if err := json.Unmarshal([]byte(msg), &request); err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}

	rsp, err := feedClient.ToggleFocusDel(context.TODO(), &request)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}

//@Params objectId, fav
func ToggleFocusFav(w http.ResponseWriter, r *http.Request) {
	log.Info("ToggleFocusFav called")
	// call the backend service
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	msg := auth_crypto.Decrypt(string(bodyBytes))
	var request proto.ManipulateFocusReq
	if err := json.Unmarshal([]byte(msg), &request); err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}

	rsp, err := feedClient.ToggleFocusFav(context.TODO(), &request)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}

//@Params date - by fetchdate
func GenFocus(w http.ResponseWriter, r *http.Request) {
	log.Info("GenFocus called")

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	msg := auth_crypto.Decrypt(string(bodyBytes))

	var request proto.ManipulateFocusReq
	if err := json.Unmarshal([]byte(msg), &request); err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rsp, err := feedClient.GenerateFocus(context.TODO(), &request)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}

//@Params array of the following:
//@Params to_cond - fill #target channel
//@Params key - fill keyword, 1 keyword vs 1 rule
//@Params contains - array of search words
func CreateRules(w http.ResponseWriter, r *http.Request) {
	log.Info("CreateRules called")
	err, request := decryptReq(r, w)
	log.Info(fmt.Sprintf("%+v", request))
	if err != nil {
		log.Info(err.Error())
		return
	}
	rsp, err := feedClient.CreateRule(context.TODO(), &request)

	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}

//@Params body == ""
func ReadRules(w http.ResponseWriter, r *http.Request) {
	log.Info("ReadRules called")
	err, request := decryptReq(r, w)
	log.Info(fmt.Sprintf("%+v", request))
	if err != nil {
		log.Info(err.Error())
		return
	}
	rsp, err := feedClient.ReadRule(context.TODO(), &request)

	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}

//@Params {"rules":[{},{},{}]}
func UpdateRules(w http.ResponseWriter, r *http.Request) {
	log.Info("ReadRules called")
	err, request := decryptReq(r, w)
	log.Info(fmt.Sprintf("%+v", request))
	if err != nil {
		log.Info(err.Error())
		return
	}
	rsp, err := feedClient.UpdateRule(context.TODO(), &request)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}

//@Params {"rules":[{},{},{}]}
func DeleteRules(w http.ResponseWriter, r *http.Request) {
	log.Info("DeleteRules called")
	err, request := decryptReq(r, w)
	log.Info(fmt.Sprintf("%+v", request))
	if err != nil {
		log.Info(err.Error())
		return
	}

	rsp, err := feedClient.DeleteRule(context.TODO(), &request)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}

func decryptReq(r *http.Request, w http.ResponseWriter) (error, proto.RuleReq) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return err, proto.RuleReq{}
	}
	msg := auth_crypto.Decrypt(string(bodyBytes))
	log.Info(msg)
	var request proto.RuleReq
	err = json.Unmarshal([]byte(msg), &request)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return err, proto.RuleReq{}
	}
	return nil, request
}
func AddRegId(w http.ResponseWriter, r *http.Request) {
	log.Info("AddRegId called")
	id := r.Header.Get("jpush_id")

	rsp, err := feedClient.AddJPushReg(context.TODO(), &proto.JPushReg{
		Id: id,
	})

	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}
func GetFocusStat(w http.ResponseWriter, r *http.Request) {
	log.Info("GetFocusStat called")
	rsp, err := feedClient.GetFocusStat(context.TODO(), &proto.RuleReq{})
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}

func Log(w http.ResponseWriter, r *http.Request) {
	log.Info("Log called")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	msg := auth_crypto.Decrypt(string(bodyBytes))
	rsp, err := feedClient.Log(context.TODO(), &proto.PlainReq{
		Msg: msg,
	})
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
	} else {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(auth_crypto.Encrypt(rsp.Msg)))
	}
}

func SetChanODay(w http.ResponseWriter, r *http.Request) {
	log.Info("SetChanODay called")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	msg := auth_crypto.Decrypt(string(bodyBytes))
	var request proto.Chans
	if err := json.Unmarshal([]byte(msg), &request); err != nil {
		fmt.Println(err.Error())
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	log.Info(fmt.Sprintf("%+v", request))
	rsp, err := feedClient.SetChanODay(context.TODO(), &request)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}

func GetChanODay(w http.ResponseWriter, r *http.Request) {
	log.Info("GetChanODay called")
	rsp, err := feedClient.GetChanODay(context.TODO(), &proto.PlainReq{})
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	rspb, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte(auth_crypto.Encrypt(string(rspb))))
}
