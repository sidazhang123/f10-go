package handler

import (
	"context"
	"encoding/json"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/sidazhang123/f10-go/srv/processor/proto/processor"
	"net/http"
)

var regexClient = processor.NewProcessorService("sidazhang123.f10.srv.processor", client.DefaultClient)

func GetPluginPath(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service

	rsp, err := regexClient.GetPluginPath(context.TODO(), &processor.GetPluginPathReq{})
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		// err!=nil when path=="nil"
		"path": rsp.JoinedPath,
		"err":  rsp.ErrMsg,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
func GetSrc(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]string
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service

	rsp, err := regexClient.GetSourceCode(context.TODO(), &processor.GetSourceCodeReq{Path: request["pluginPath"]})
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		// err!=nil when path=="nil"
		"sourceCode": rsp.SourceCode,
		"err":        rsp.ErrMsg,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

// update src and build so
func Update(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]string
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service

	rsp, err := regexClient.BuildSo(context.TODO(), &processor.BuildSoReq{PluginPath: request["pluginPath"], SourceCode: request["sourceCode"]})

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Path,
		"err": rsp.ErrMsg,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

// test str
func Test(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]string
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service

	rsp, err := regexClient.RegexTest(context.TODO(), &processor.RegexReq{PluginPath: request["pluginPath"], TestStr: request["testStr"]})

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.ResStr,
		"err": rsp.ErrMsg,
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
