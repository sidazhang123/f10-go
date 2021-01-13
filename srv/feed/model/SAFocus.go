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

func SAFocus(rules []ReadRuleDTO, date string) error {
	spaceRe := regexp.MustCompile(`[\s]+`)
	type refinedSATable struct {
		Date        string `json:"date" bson:"date"`
		Percentage  string `json:"percentage" bson:"percentage"`
		TableString string `json:"tableString" bson:"tableString"`
		Trend       string `json:"trend" bson:"trend"`
		Variation   string `json:"variation" bson:"variation"`
	}

	fieldKey := "流通占比表0"
	//determine the latest date
	var err error
	if len(date) != 10 {
		err, date = s.FindLatestFetchTime("shareholder_analysis")
		if err != nil {
			return err
		}
	}

	// use as genTime because the latest date in db is not necessarily today's
	curDate := time.Now().UTC().Add(8 * time.Hour)
	curDateStr := curDate.Format(common.TimestampLayout[:10])
	skipDate := curDate.AddDate(0, 0, -Params.SkipDays)
	//read refined shareholder_analysis of the latest date
	err, cur := FindByFetchTime("shareholder_analysis", date)
	if err != nil {
		return err
	}
	focusCount := 0
	for cur.Next(context.Background()) {

		//get the record of the latest on code
		var refinedSAItem bson.M
		err = cur.Decode(&refinedSAItem)
		if err != nil {
			return err
		}
		if _, ok := refinedSAItem[fieldKey]; !ok {
			continue
		}
		// skip if fetchtime too old
		if skipDate.After(refinedSAItem["updatetime"].(primitive.DateTime).Time().UTC()) {
			continue
		}
		var refinedSATable refinedSATable
		err := json.Unmarshal([]byte(refinedSAItem[fieldKey].(string)), &refinedSATable)
		if err != nil {
			return err
		}
		tableString := refinedSATable.TableString
		tabUpdateTime := refinedSATable.Date
		var rows [][]string
		for i, r := range strings.Split(strings.Split(tableString, "───────────────────────────────────────")[0], "\n") {
			if i == 0 {
				continue
			}
			r = strings.TrimSpace(r)
			if len(r) > 0 {
				r = spaceRe.ReplaceAllString(r, " ")
				cells := strings.Split(r, " ")
				if len(cells) > 2 {
					rows = append(rows, cells)
				} else {
					if len(rows) == 0 {
						continue
					}
					// merge name segments of multilines
					rows[len(rows)-1][0] += cells[0]
				}

			}
		}
		// it's only reasonable to bind channel with state to achieve minimal #Focus in terms of one code
		chanToFocus := map[string]Focus{}
		for _, row := range rows {
			shareholder, holding, state := row[0], row[1], row[len(row)-1]
			for _, rule := range rules {
				////input malformed
				//if len(rule.Cond2) != 1 {
				//	return fmt.Errorf("length of rule.Cond2(state) should've been 1\n%+v", rule)
				//}
				//st := rule.Cond2[0]
				//if len(st) == 0 {
				//	continue
				//}
				/*
					check sh=李四 & st=↑, "张三李四公司 .... ↑5%\n李四集团 .... ↑20%"
					{ chan: msg{msg:"张三李四公司 ↑5%,李四集团 ↑20%", markRed:[]string{李四,↑}}
					}
				*/
				for _, st := range []string{"↑", "新进", "未变", "↓"} {
					for _, sh := range rule.Cond1 {
						if len(sh) == 0 {
							continue
						}
						if (strings.Contains(state, st)) && func(shareholder, sh, channel string) bool {
							if channel == "题材达人" {
								return strings.TrimSpace(shareholder) == strings.TrimSpace(sh)
							} else {
								return strings.Contains(shareholder, sh)
							}
						}(shareholder, sh, rule.Channel) {
							//emit
							if _, ok := chanToFocus[rule.Channel]; ok {
								if _, ok := chanToFocus[rule.Channel].Keys[st]; ok {
									cm := chanToFocus[rule.Channel].Keys[st]
									in := false
									for _, v := range cm.Contain {
										if sh == v {
											in = true
											break
										}
									}
									if !in {
										chanToFocus[rule.Channel].Keys[st] = Contain{
											Msg:     cm.Msg + fmt.Sprintf("%s-%s-%s\n", shareholder, holding, state),
											Contain: append(cm.Contain, sh),
										}
									} else {
										chanToFocus[rule.Channel].Keys[st] = Contain{
											Msg:     cm.Msg + fmt.Sprintf("%s-%s-%s\n", shareholder, holding, state),
											Contain: cm.Contain,
										}
									}
								} else {
									chanToFocus[rule.Channel].Keys[st] = Contain{
										Msg:     fmt.Sprintf("%s-%s-%s\n", shareholder, holding, state),
										Contain: []string{sh},
									}
								}

							} else {
								chanToFocus[rule.Channel] = Focus{
									Code:      refinedSAItem["code"].(string),
									Name:      refinedSAItem["name"].(string),
									RuleId:    rule.ID.Hex(),
									Fetchtime: refinedSAItem["fetchtime"].(primitive.DateTime).Time().UTC().Format(common.TimestampLayout[:10]),
									Keys: map[string]Contain{st: {
										Msg:     fmt.Sprintf("%s-%s-%s\n", shareholder, holding, state),
										Contain: []string{sh},
									}},
								}

							}
						}
					}
				}

			}
		}

		for c, f := range chanToFocus {
			toInsert := map[string]interface{}{}
			toInsert["code"] = f.Code
			toInsert["name"] = f.Name
			toInsert["tabupdatetime"] = tabUpdateTime
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
	log.Info("Gen Shareholder_analysis Focus done")
	// jpush
	return makeJPush(fmt.Sprintf("股东研究 %d 条", focusCount))
}
