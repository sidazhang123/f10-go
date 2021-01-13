package model

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"strings"
	"time"
)

func LTFocus(rulesDTOList []ReadRuleDTO, date string) error {
	type msg struct {
		Msg string `json:"msg" bson:"msg"`
	}
	dateRe := regexp.MustCompile(`[\d]{4}-[\d]{2}-[\d]{2}`)
	//determine the latest date
	var err error
	if len(date) != 10 {
		err, date = s.FindLatestFetchTime("latest_tips")
		if err != nil {
			return err
		}
	}

	// use as genTime because the latest date in db is not necessarily today's
	curDate := time.Now().UTC().Add(8 * time.Hour)
	curDateStr := curDate.Format(common.TimestampLayout[:10])
	skipDate := curDate.AddDate(0, 0, -Params.SkipDays)
	//read refined latest_tips of the latest date
	err, cur := FindByFetchTime("latest_tips", date)
	if err != nil {
		return err
	}
	// deconsolidate the key if contains ",/，" separator
	var rules []ReadRuleDTO
	for _, v := range rulesDTOList {
		if strings.Contains(v.Key, ",") || strings.Contains(v.Key, "，") {
			keys := func(c string) []string {
				a := strings.Split(c, ",")
				b := strings.Split(c, "，")
				if len(a) > len(b) {
					return a
				} else {
					return b
				}
			}(v.Key)
			for _, k := range keys {
				if len(strings.TrimSpace(k)) > 0 {
					v.Key = k
					rules = append(rules, v)
				}
			}
		} else {
			rules = append(rules, v)
		}
	}
	focusCount := 0
	for cur.Next(context.Background()) {
		//{"情况一":{"Code":"000001","Name":"平安银行","Fetchtime":"2020-01-01","Keys":{"未分":{"Str":"blablabla","Contain":["第三季度","1月至9月"]}}}}
		chanToFocus := map[string]Focus{}
		//get the record of the latest on code
		var refinedLTItem bson.M
		err = cur.Decode(&refinedLTItem)
		if err != nil {
			return err
		}
		// skip if fetchtime too old
		if skipDate.After(refinedLTItem["updatetime"].(primitive.DateTime).Time().UTC()) {
			continue
		}
		for _, rule := range rules {
			//skip a rule when the latest_tips record does not has the rule's key
			if _, ok := refinedLTItem[rule.Key]; !ok {
				continue
			}
			var msg msg
			itemMsg := refinedLTItem[rule.Key].(string)
			err := json.Unmarshal([]byte(itemMsg), &msg)
			if err != nil {
				return err
			}

			/*select by the rules:
			ignore when any keyword of Cond1(serves as contains) does not match, and
			when any keyword of Cond2(serves as excludes) matches
			*/
			notSatisfied := false

			for _, contains := range rule.Cond1 {
				if !strings.Contains(msg.Msg, contains) {
					notSatisfied = true
					break
				}
				dates := dateRe.FindAllStringSubmatch(msg.Msg, -1)
				if len(dates) > 0 {
					for _, d := range dates {
						dateT, e := time.Parse(common.TimestampLayout[:10], d[0])
						if e != nil {
							continue
						}
						if skipDate.After(dateT) {
							notSatisfied = true
							break
						}
					}
				}
			}
			if !notSatisfied {
				for _, excludes := range rule.Cond2 {
					if strings.Contains(msg.Msg, excludes) {
						notSatisfied = true
						break
					}
				}
			}

			if !notSatisfied {
				if _, ok := chanToFocus[rule.Channel]; !ok {
					chanToFocus[rule.Channel] = Focus{
						Code:      refinedLTItem["code"].(string),
						RuleId:    rule.ID.Hex(),
						Name:      refinedLTItem["name"].(string),
						Fetchtime: refinedLTItem["fetchtime"].(primitive.DateTime).Time().UTC().Format(common.TimestampLayout[:10]),
						Keys: map[string]Contain{rule.Key: {
							Msg:     msg.Msg,
							Contain: rule.Cond1,
						}},
					}
					// in case that >1 rules match one code while sending to the same channel
				} else {
					if _, ok := chanToFocus[rule.Channel].Keys[rule.Key]; !ok {
						chanToFocus[rule.Channel].Keys[rule.Key] = Contain{
							Msg:     msg.Msg,
							Contain: rule.Cond1,
						}
					} else {
						curMsg := chanToFocus[rule.Channel].Keys[rule.Key].Msg
						curContain := chanToFocus[rule.Channel].Keys[rule.Key].Contain
						curContain = append(curContain, rule.Cond1...)
						chanToFocus[rule.Channel].Keys[rule.Key] = Contain{
							Msg:     curMsg,
							Contain: unique(curContain),
						}
					}
				}
			}
		}
		for c, f := range chanToFocus {
			toInsert := map[string]interface{}{}
			toInsert["code"] = f.Code
			toInsert["name"] = f.Name
			toInsert["rid"] = f.RuleId
			toInsert["fetchtime"] = f.Fetchtime
			keys, err := json.Marshal(f.Keys)
			if err != nil {
				return err
			}
			toInsert["keys"] = string(keys)
			toInsert["chan"] = c
			toInsert["gentime"] = curDateStr
			toInsert["del"] = 0
			toInsert["fav"] = 0
			toInsert["uid"] = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%s-%s", f.Code, toInsert["keys"], c))))
			err = InsertOneFocus("focus", toInsert)
			if err != nil && !strings.Contains(err.Error(), "dup") {
				return err
			}
			if err == nil {
				focusCount += 1
			}
		}
	}
	log.Info("Gen Latest_tips Focus done")
	// jpush
	return makeJPush(fmt.Sprintf("最新提醒 %d 条", focusCount))
}
func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
