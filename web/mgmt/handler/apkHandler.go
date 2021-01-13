package handler

import (
	"crypto/sha256"
	"fmt"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/sidazhang123/f10-go/web/mgmt/auth_crypto"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var localVersion string
var sha string
var apkPath = "ota_update.apk"

func Monitor() chan error {
	err := make(chan error)
	go func() {
		var lock sync.RWMutex
		verRe := regexp.MustCompile(`"tag_name":"([^"]+)"`)
		urlRe := regexp.MustCompile(`"browser_download_url":"([^"]+)"`)
		for {
			// get info from remote
			rsp, e := http.Get("https://gitee.com/api/v5/repos/sidazhang123/f10-flutter-release/releases/latest")
			if e != nil {
				err <- e
				continue
			}
			body, e := ioutil.ReadAll(rsp.Body)
			if e != nil {
				err <- e
				continue
			}
			b := string(body)
			ver := verRe.FindStringSubmatch(b)
			if len(ver) < 2 {
				err <- fmt.Errorf(fmt.Sprintf("failed to get version number; %+v", ver))
				time.Sleep(5 * time.Minute)
				continue
			}
			version := ver[1]
			u := urlRe.FindStringSubmatch(b)
			if len(u) < 2 {
				err <- fmt.Errorf(fmt.Sprintf("failed to get asset path; %+v", u))
				time.Sleep(5 * time.Minute)
				continue
			}
			url := u[1]
			// cmp with local
			if version != localVersion {
				//download
				e := downloadFile(url, "ota_update.apk.download")
				if e != nil {
					err <- e
					continue
				}
				e, s := filesha256("ota_update.apk.download")
				if e != nil {
					err <- e
					continue
				}
				lock.Lock()
				sha = s
				localVersion = version
				e = os.Rename("ota_update.apk.download", apkPath)
				if e != nil {
					err <- e
					continue
				}
				lock.Unlock()
				fmt.Println(version, sha)
			}
			time.Sleep(15 * time.Minute)
		}
	}()
	return err
}
func GetApkInfo(w http.ResponseWriter, r *http.Request) {
	// read info from local
	fmt.Fprintf(w, "%s:%s", localVersion, sha)
}

func DownloadApk(w http.ResponseWriter, r *http.Request) {
	// provide resource
	fi, err := os.Stat(apkPath)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, auth_crypto.Encrypt(err.Error()), 500)
		return
	}
	// get the size
	size := fi.Size()
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	w.Header().Set("Content-Disposition", "attachment; filename=ota_update.apk")
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, apkPath)
}

func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
func filesha256(path string) (error, string) {
	f, err := os.Open(path)
	if err != nil {
		return err, ""
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err, ""
	}

	return nil, fmt.Sprintf("%x", h.Sum(nil))
}
