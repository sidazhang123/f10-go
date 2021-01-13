package model

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

//"最新公告":{"msg":""},
//"最新报道":{"msg":""},
//"业绩预告":{"msg":""},
//"特别处理":{"msg":""},
//"最新提醒增发":{"0":""，"1":""},
//"股东户数变化":{"本期":"","股东数":"","上期":"", "增减":"","户数":"","幅度":""}
var ltnews = map[string]func(string) (error, string){
	"最新公告": handleMsg, "最新报道": handleMsg, "业绩预告": handleMsg,
	"特别处理": handleMsg, "最新提醒增发": handleSeq, "股东户数变化": handleMap}

func handleMsg(jsonStr string) (error, string) {
	type msg struct {
		Msg string `json:"msg" bson:"msg"`
	}
	var m msg
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return err, ""
	}
	return nil, m.Msg
}
func handleSeq(jsonStr string) (error, string) {
	var m map[string]string
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return err, ""
	}
	var vals []string
	for _, v := range m {
		vals = append(vals, v)
	}
	return nil, strings.Join(vals, "\n")
}
func handleMap(jsonStr string) (error, string) {
	var m map[string]string
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return err, ""
	}
	var vals []string
	for k, v := range m {
		vals = append(vals, k+v)
	}
	return nil, strings.Join(vals, ", ")
}
func LTNews(rulesDTOList []ReadRuleDTO) error {
	curDateStr := time.Now().UTC().Add(8 * time.Hour).Format(common.TimestampLayout[:10])
	focusCount := 0
	n := 2
	for _, v := range rulesDTOList {
		for _, code := range v.Cond1 {
			code = strings.TrimSpace(code)
			if len(code) > 0 {
				err, lastNRecords := GetLastNByField("latest_tips", "fetchtime", code, int64(n))
				if err != nil {
					return err
				}
				toEmit := map[string]Contain{}
				if len(lastNRecords) != n {
					continue
				}
				newRec, oldRec := lastNRecords[0], lastNRecords[1]
				for field, jsonStrNew := range newRec {
					// field must be selected
					if handler, ok := ltnews[field]; ok {
						//field in the last record
						if jsonStrOld, ok := oldRec[field]; ok {
							//values do not match
							if jsonStrNew != jsonStrOld {
								err := emit(handler, jsonStrNew, toEmit, field)
								if err != nil {
									return err
								}
							}
						} else {
							err := emit(handler, jsonStrNew, toEmit, field)
							if err != nil {
								return err
							}
						}
					}
				}

				if len(toEmit) > 0 {

					toInsert := map[string]interface{}{}
					toInsert["code"] = code
					toInsert["name"] = newRec["name"].(string)
					toInsert["rid"] = v.ID.Hex()
					toInsert["fetchtime"] = newRec["fetchtime"].(primitive.DateTime).Time().UTC().Format(common.TimestampLayout[:10])
					keys, err := json.Marshal(toEmit)
					if err != nil {
						return err
					}
					toInsert["keys"] = string(keys)
					toInsert["chan"] = v.Channel
					toInsert["gentime"] = curDateStr
					toInsert["del"] = 0
					toInsert["fav"] = 0
					toInsert["uid"] = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%s-%s", code, toInsert["keys"], v.Channel))))
					err = InsertOneFocus("focus", toInsert)
					if err != nil && !strings.Contains(err.Error(), "dup") {
						return err
					}
					if err == nil {
						focusCount += 1
					}
				}

			}
		}
	}

	log.Info("Gen Latest_tips_news Focus done")
	// jpush
	return makeJPush(fmt.Sprintf("news变化 %d 条", focusCount))
}

func emit(handler func(string) (error, string), jsonStrNew interface{}, toEmit map[string]Contain, field string) error {
	err, fieldStrToEmit := handler(jsonStrNew.(string))
	if err != nil {
		return err
	}
	toEmit[field] = Contain{
		Msg:     fieldStrToEmit,
		Contain: []string{},
	}
	return nil
}
