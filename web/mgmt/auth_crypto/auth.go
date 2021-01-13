package auth_crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	z "github.com/sidazhang123/f10-go/plugins/zap"
	"hash"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	sharedSecret []byte
	encKey       []byte
	mac          hash.Hash
	log          = z.GetLogger()
)

func Init() {
	b, e := ioutil.ReadFile("secret")
	if e != nil {
		log.Error("failed to read secret on server: " + e.Error())
		return
	}
	sharedSecret = b
	hasher := md5.New()
	hasher.Write(b)
	encKey = []byte(hex.EncodeToString(hasher.Sum(nil)))

}

func getSecret() (error, []byte) {
	if sharedSecret == nil || len(sharedSecret) == 0 {
		return fmt.Errorf("getSharedSecret: cannot get secret"), nil
	}
	return nil, sharedSecret
}

type innerFn func(w http.ResponseWriter, r *http.Request)

func AuthWrapper(fn innerFn, debug bool, weak bool) func(w http.ResponseWriter, r *http.Request) {
	if debug {
		return fn
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			e error
			b bool
		)

		if weak == false {
			e, b = VerifyIdentity(r)
		} else {
			e, b = VerifyWeakIdentifyfunc(r)
		}
		if e != nil {
			http.Error(w, e.Error(), 503)
		}
		if b == false {
			http.Error(w, "unauthorized access", 401)
		} else {
			fn(w, r)
		}
	}
}
func VerifyWeakIdentifyfunc(r *http.Request) (error, bool) {
	e, ss := getSecret()
	if e != nil {
		return e, false
	}
	sStr := string(ss)
	str := r.Header.Get("sno")
	if len(str) != len(sStr) {
		return nil, false
	}
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	if sStr != string(runes) {
		return nil, false
	}
	return nil, true
}

//sha256(sharedSecret+hmac(sharedSecret, nonce))
func VerifyIdentity(r *http.Request) (error, bool) {
	e, ss := getSecret()
	if e != nil {
		return e, false
	}
	str := r.Header.Get("sno")
	if str == "" {
		return nil, false
	}

	now := time.Now().UTC().Add(8 * time.Hour)
	nonce := []string{now.Format(common.TimestampLayout[:16]), now.Add(time.Minute).Format(common.TimestampLayout[:16]), now.Add(-1 * time.Minute).Format(common.TimestampLayout[:16])}
	for _, n := range nonce {
		mac := hmac.New(sha256.New, ss)
		mac.Write([]byte(n))
		if str == fmt.Sprintf("%x", sha256.Sum256(append(ss, mac.Sum(nil)...))) {
			return nil, true
		}
	}
	return nil, false
}

func genIdentityCode() {
	nonce := time.Now().UTC().Add(8 * time.Hour).Format(common.TimestampLayout[:16])[:15]
	mac := hmac.New(sha256.New, sharedSecret)
	mac.Write([]byte(nonce))
	fmt.Println(fmt.Sprintf("%x", sha256.Sum256(append(sharedSecret, mac.Sum(nil)...))))
}
