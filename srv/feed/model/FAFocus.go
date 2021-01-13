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
	"strconv"
	"strings"
	"time"
)

func FAFocus(rules []ReadRuleDTO, date string) error {
	extractNum := regexp.MustCompile(`[-]?[\d.]+`)
	npKey := "净利润(万元)"
	npPerKey := "净利润增长率(%)"

	//determine the latest date
	var err error
	if len(date) != 10 {
		err, date = s.FindLatestFetchTime("financial_analysis")
		if err != nil {
			return err
		}
	}

	// use as genTime because the latest date in db is not necessarily today's
	curDate := time.Now().UTC().Add(8 * time.Hour)
	curDateStr := curDate.Format(common.TimestampLayout[:10])
	skipDate := curDate.AddDate(0, 0, -Params.SkipDays)
	//read refined shareholder_analysis of the latest date
	log.Info(fmt.Sprintf("date: %s", date))
	err, cur := FindByFetchTime("financial_analysis", date)
	if err != nil {
		return err
	}
	focusCount := 0
	for cur.Next(context.Background()) {
		//get the record of the latest on code
		var refinedFAItem bson.M
		err = cur.Decode(&refinedFAItem)
		if err != nil {
			return err
		}
		//if refinedFAItem["code"]!="000791"{continue}//testingcode
		// must have the 2 fields
		if _, ok := refinedFAItem[npKey]; !ok {
			//fmt.Printf("dont have npKey: %+v\n",refinedFAItem)//testingcode
			continue
		}
		if _, ok := refinedFAItem[npPerKey]; !ok {
			//fmt.Printf("dont have npPerKey: %+v\n",refinedFAItem)//testingcode
			continue
		}
		// skip if fetchtime too old
		if skipDate.After(refinedFAItem["updatetime"].(primitive.DateTime).Time().UTC()) {
			//fmt.Printf("affected by skipDate: %+v\n",refinedFAItem)//testingcode
			continue
		}
		// unmarshal the net profit field
		var npD2NStr map[string]string
		err := json.Unmarshal([]byte(refinedFAItem[npKey].(string)), &npD2NStr)
		if err != nil {
			return err
		}
		// get the latest
		var dateKey string
		var maxDate time.Time
		for d := range npD2NStr {
			t, err := time.Parse(common.TimestampLayout[:10], d)
			if err != nil {
				continue
			}
			if t.After(maxDate) {
				maxDate = t
				dateKey = d
			}
		}

		var npperD2NStr map[string]string
		err = json.Unmarshal([]byte(refinedFAItem[npPerKey].(string)), &npperD2NStr)
		if err != nil {
			return err
		}
		// check and transform d2n type
		var np, npper float64
		if n, err := strconv.ParseFloat(extractNum.FindString(npD2NStr[dateKey]), 64); err == nil {
			np = n
		}
		if n, err := strconv.ParseFloat(extractNum.FindString(npperD2NStr[dateKey]), 64); err == nil {
			npper = n
		}
		//fmt.Printf("np,npper %.2f, %.2f\n",np,npper)//testingcode
		chanToFocus := map[string]Focus{}
		for _, rule := range rules {
			/*
				rule: cond1:["10000|0","3000:15","800|30"]
				np thresholds are guaranteed to be ordered
			*/
			for _, combStr := range rule.Cond1 {
				comb := strings.Split(combStr, "|")
				npTh, _ := strconv.ParseFloat(comb[0], 64)
				npperTh, _ := strconv.ParseFloat(comb[1], 64)
				//fmt.Printf("npTh,npperTh %.2f, %.2f\n",npTh,npperTh)//testingcode
				// from big to small thresholds, emit if match
				if np > npTh && npper > npperTh {
					chanToFocus[rule.Channel] = Focus{
						Code:      refinedFAItem["code"].(string),
						RuleId:    rule.ID.Hex(),
						Name:      refinedFAItem["name"].(string),
						Fetchtime: refinedFAItem["fetchtime"].(primitive.DateTime).Time().UTC().Format(common.TimestampLayout[:10]),
						Keys: map[string]Contain{"": {
							Msg:     fmt.Sprintf("净利%.2f万元，增长%.2f%% @%s月", np, npper, dateKey[len(dateKey)-5:len(dateKey)-3]),
							Contain: []string{"净利", "增长"},
						}},
					}
					break
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
	log.Info("Gen Financial_analysis Focus done")
	// jpush
	return makeJPush(fmt.Sprintf("财务分析 %d 条", focusCount))
}
