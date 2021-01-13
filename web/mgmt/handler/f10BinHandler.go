package handler

import (
	"github.com/micro/go-micro/v2/util/log"
	"github.com/sidazhang123/f10-go/web/mgmt/auth_crypto"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetBin(w http.ResponseWriter, r *http.Request) {
	log.Info("GetBin called")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	// "/f10a/2020-01-01.a"
	binPath := auth_crypto.Decrypt(string(bodyBytes))
	log.Info("GetBin " + binPath)
	// provide resource
	fi, err := os.Stat(binPath)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	// get the size
	size := fi.Size()
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	filepath := strings.Split(binPath, "/")
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath[len(filepath)-1])
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, binPath)
}
