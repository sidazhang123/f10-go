package handler

import (
	"context"
	"encoding/json"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/util/log"
	proto "github.com/sidazhang123/f10-go/srv/accumulator/proto/accumulator"
	"github.com/sidazhang123/f10-go/web/mgmt/auth_crypto"
	"io/ioutil"
	"net/http"
)

var accClient = proto.NewAccumulatorService("sidazhang123.f10.srv.accumulator", grpc.NewClient(grpc.MaxRecvMsgSize(1024*1024*50), grpc.MaxSendMsgSize(1024*1024*50)))

func GetRepr(w http.ResponseWriter, r *http.Request) {
	log.Info("GetRepr called")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	msg := auth_crypto.Decrypt(string(bodyBytes))

	var request proto.ReprReq
	if err := json.Unmarshal([]byte(msg), &request); err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}

	rsp, err := accClient.GetRepr(context.TODO(), &request)

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
