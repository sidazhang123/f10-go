package model

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type res struct {
	name      string
	ts        string
	table     string
	firstSeen string
	dateDiff  string
}

func (s Service) GenerateOperationalAnalysisDiffCSV(ts string) error {
	diff := map[string]res{}
	var (
		today    time.Time
		today_ts string
	)
	if ts == "" {
		today = time.Now().UTC().Add(8 * time.Hour)
		today_ts = today.Format(common.TimestampLayout[:10])
	} else {
		today_ts = ts
		today, _ = time.Parse(common.TimestampLayout[:10], ts)
	}

	//find today's records
	var newRecord []map[string]interface{}
	err, cur := FindByFetchTime("operational_analysis", today_ts)
	if err != nil {
		return err
	}
	//put in a map
	err = cur.All(context.TODO(), &newRecord)
	if err != nil {
		return err
	}
	newRecordCount := len(newRecord)
	//try to find yesterday's records. Try one day earlier if yesterday's count is 0
	var oldRecord map[string]map[string]interface{}
	for i := 1; i < 6; i++ {
		oldRecord = map[string]map[string]interface{}{}
		err, cur := FindByFetchTime("operational_analysis", today.AddDate(0, 0, -i).Format(common.TimestampLayout[:10]))
		if err != nil {
			log.Error(err.Error())
		}
		//put in a map
		for cur.Next(context.TODO()) {
			var o map[string]interface{}
			err = cur.Decode(&o)
			if err != nil {
				log.Error(err.Error())
			} else {
				oldRecord[o["code"].(string)] = o
			}
		}
		if int(0.8*float64(newRecordCount)) < len(oldRecord) && len(oldRecord) < int(1.2*float64(newRecordCount)) {
			break
		}
		if i == 5 {
			log.Warn("cannot find valid operational_analysis records in last 5 days")
			return nil
		}
	}
	//compare
	for _, newMap := range newRecord {
		code := newMap["code"].(string)
		err, newLatestDate := getLatestDateKey(newMap)
		if err != nil {
			log.Error(err.Error())
		}
		// have 2 matched stocks by code
		if oldMap, ok := oldRecord[code]; ok {
			err, oldLatestDate := getLatestDateKey(oldMap)
			if err != nil {
				log.Error(err.Error())
			}

			if newLatestDate.After(oldLatestDate) {
				chosenDate := newLatestDate.Format(common.TimestampLayout[:10])
				diff[code] = res{
					name:      newMap["name"].(string),
					ts:        chosenDate,
					table:     newMap[chosenDate].(string),
					firstSeen: "",
					dateDiff:  fmt.Sprintf("%s => %s", oldLatestDate.Format(common.TimestampLayout[:10]), newLatestDate.Format(common.TimestampLayout[:10])),
				}
			}
			// new record has a earlier date?
			if newLatestDate.Before(oldLatestDate) {
				log.Error(fmt.Sprintf("new record has a earlier date?\nc: %s nDate: %s oDate: %s", code, newLatestDate.Format(common.TimestampLayout), oldLatestDate.Format(common.TimestampLayout)))
			}
			// newly added stock code
		} else {
			chosenDate := newLatestDate.Format(common.TimestampLayout[:10])
			diff[code] = res{
				name:      newMap["name"].(string),
				ts:        chosenDate,
				table:     newMap[chosenDate].(string),
				firstSeen: "FirstSeen",
				dateDiff:  fmt.Sprintf("%s => %s", "null", newLatestDate.Format(common.TimestampLayout[:10])),
			}
		}
	}
	if len(diff) > 0 {
		fileName := fmt.Sprintf("经营分析变化_%s.csv", today_ts)
		generateCSV(diff, fileName)
		jpushmsg := fmt.Sprintf("经营分析有 %d 条变化 => http://mgmt9.pro-klick.xyz/ad/%s", len(diff), fileName)
		return makeJPush(jpushmsg)
	}
	return nil
}

func getLatestDateKey(m map[string]interface{}) (error, time.Time) {
	date, _ := time.Parse(common.TimestampLayout[:10], "1970-01-01")
	isFound := false
	for k := range m {
		matched, _ := regexp.MatchString(`[\d]{4}-[\d]{2}-[\d]{2}`, k)
		if matched {
			d, err := time.Parse(common.TimestampLayout[:10], k)
			if err != nil {
				continue
			}
			if d.After(date) {
				date = d
				isFound = true
			}
		}

	}
	if isFound {
		return nil, date
	} else {
		return fmt.Errorf("failed to get latest date key: %s", m["code"]), date
	}

}

func generateCSV(diff map[string]res, fileName string) {
	var data [][]string
	for code, res := range diff {
		var fin [][]string
		tmp := strings.Split(res.table, "\n")
		// weird bug requires further check. Due to unreleased 'diff'?
		if len(tmp) < 3 {
			log.Warn(code)
			log.Warn(res.table)
			continue
		}
		block := tmp[2:]
		for _, line := range block {
			if hasSubStr(line, []string{"(产", "品)", "(行", "业)"}) {
				cells := strings.Split(strings.Trim(line, "|"), "|")
				header := removeSubStr(cells[0], []string{"(产", "品)", "(行", "业)", "(补", "充)", " "})
				percentage := removeSubStr(cells[2], []string{"%", " "})
				if percentage == "-" {
					percentage = "-1"
				}
				fin = append(fin, []string{header, percentage})
			}
		}
		sort.SliceStable(fin, func(i, j int) bool {
			x, _ := strconv.ParseFloat(fin[i][1], 64)
			y, _ := strconv.ParseFloat(fin[j][1], 64)
			return x > y
		})
		s := ""
		for _, f := range fin {
			p := round(f[1])
			if p > 1 {
				s += fmt.Sprintf("%s%d，", f[0], p)
			}
		}
		s = strings.Trim(s, "，")
		data = append(data, []string{code, res.name, s, res.ts, res.firstSeen, res.dateDiff})
	}
	sort.SliceStable(data, func(i, j int) bool {
		x, _ := strconv.ParseInt(data[i][0], 10, 64)
		y, _ := strconv.ParseInt(data[j][0], 10, 64)
		return x < y
	})
	exportCSV(fileName, data)
}

func round(x string) int {
	if s, err := strconv.ParseFloat(x, 64); err == nil {
		return int(math.Floor(s + 0.5))
	}
	return 0
}

func hasSubStr(s string, subStrSet []string) bool {
	b := false
	for _, i := range subStrSet {
		if strings.LastIndex(s, i) > -1 {
			b = true
			break
		}
	}
	return b
}

func removeSubStr(s string, subStrSet []string) string {
	for _, i := range subStrSet {
		s = strings.Replace(s, i, "", -1)
	}
	return s
}

func exportCSV(filePath string, data [][]string) {
	filePath = Params.LocalFilePathPrefix + filePath
	fp, err := os.Create(filePath)
	if err != nil {
		log.Error(fmt.Sprintf("failed to create %s\n %v", filePath, err))
		return
	}
	defer fp.Close()
	_, _ = fp.WriteString("\xEF\xBB\xBF") // write utf8 BOM
	w := csv.NewWriter(fp)
	err = w.WriteAll(data)
	if err != nil {
		log.Error(fmt.Sprintf("failed writing data to csv\n %v", err))
		return
	}
	w.Flush()
	fp.Close()
}
