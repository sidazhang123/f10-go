package debug

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func (d *Debug) Shareholder_analysis(s string) (map[string]interface{}, []error) {

	filter := regexp.MustCompile(`[ ｜─├┤└┘/\n\r\t]+`)
	regex := []*regexp.Regexp{
		regexp.MustCompile(`\n【1\.控股股东与实际控制人】[^实]{10,}实际控制人[^｜]*｜([\s\S]*?)\(([^%]*?%)[^%]*?\)[^【]*?【`),
		regexp.MustCompile(`\n【1\.控股股东与实际控制人】[^第]{10,}控股股东[^｜]*｜([\s\S]*?)\(([^%]*?%)[^%]*?\)[^【]*?【`),
		regexp.MustCompile(`\n【1\.控股股东与实际控制人】[^第]{10,}第一大股东[^｜]*｜([\s\S]*?)\(([^%]*?%)[^%]*?\)[^【]*?【`)}

	var res = map[string]interface{}{}
	var shareholderMap = map[string]string{}
	var err []error
	for _, re := range regex {
		t := re.FindStringSubmatch(s)
		if t != nil {
			shareholder := filter.ReplaceAllString(t[1], "")
			share := filter.ReplaceAllString(t[2], "")
			shareholderMap[shareholder] = share
		}
	}
	if len(shareholderMap) == 0 {
		err = append(err, fmt.Errorf("shareholder info not found"))
		res["股东控股"] = ""
	} else {
		shareholderJson, e := marshalShareholder(shareholderMap)
		if e != nil {
			err = append(err, e)
			res["股东控股"] = ""
		} else {
			res["股东控股"] = shareholderJson
		}
	}

	t := regexp.MustCompile(`【3\.股东变化】[\s\S]+【3\.股东变化】([\s\S]+)【4\.`).FindStringSubmatch(s)
	if t == nil {
		t = regexp.MustCompile(`【3\.股东变化】[\s\S]+【3\.股东变化】([\s\S]+)$`).FindStringSubmatch(s)
		if t == nil {
			res["流通占比表"] = ""
			return res, err
		}
	}
	body := t[1]
	bodyList := strings.Split(body, "截至日期")
	tableCount := 0
	for _, i := range bodyList {
		if len(strings.TrimSpace(i)) == 0 || !strings.Contains(i, "累计占流通股比例") {
			continue
		}
		tableMap := map[string]string{}
		t = regexp.MustCompile(`[：:][\d]{2}([\d]{2}-[\d]{2}-[\d]{2})[\s\S]*?累计占流通股比例[：:]([\d.-]+%)[^%]*较上期变化[：:](.+)\n`).FindStringSubmatch(i)
		if len(t) != 4 {
			fmt.Println("~~~~~~~~~~~~~~~~~~")
			fmt.Println(i)
			fmt.Printf("+++++++++++++++++\n%+v\n", t)
			fmt.Println(len(t))
			fmt.Println("~~~~~~~~~~~~~~~~~~")
		}
		tableMap["date"] = t[1]
		tableMap["percentage"] = t[2]
		variation := strings.Split(strings.TrimSpace(t[3]), "股")
		// case  较上期变化:-
		if len(variation) == 2 {
			tableMap["variation"] = variation[0]
			tableMap["trend"] = variation[1]
		} else {
			tableMap["variation"] = variation[0]
			tableMap["trend"] = variation[0]
		}
		var current, withdraw string
		tableSeg := regexp.MustCompile(`─────────────────────────────────────([\s\S]+)─────────────────────────────────────`).FindStringSubmatch(i)
		if tableSeg == nil {
			tableSeg = regexp.MustCompile(`─────────────────────────────────────([\s\S]+)`).FindStringSubmatch(i)
			if tableSeg == nil {
				current, withdraw = "", ""
			} else {
				current, withdraw = tableSeg[1], ""
			}
		} else {
			tableList := strings.Split(tableSeg[1],
				"─────────────────────────────────────")
			current = strings.TrimSpace(tableList[0])
			if len(tableList) >= 2 {
				withdraw = strings.Trim(tableList[len(tableList)-1], "\n")
			} else {
				withdraw = ""
			}

		}

		tableMap["tableString"] = "股东名称 (单位:万股)            持股数 占流通股比(%)  股东性质    增减情况\n" + current + "\n───────────────────────────────────────\n" + withdraw
		tableJson, e := marshalShareholder(tableMap)
		if e != nil {
			err = append(err, e)
		}
		res[fmt.Sprintf("流通占比表%d", tableCount)] = tableJson
		tableCount += 1
	}
	return res, err
}

func marshalShareholder(s map[string]string) (string, error) {

	leafItem, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(leafItem), nil
}
