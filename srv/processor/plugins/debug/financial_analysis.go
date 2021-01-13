package debug

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var regex = []*regexp.Regexp{regexp.MustCompile(`(财务指标.*?)\n`),
	regexp.MustCompile(`(审计意见.*?\n.*?)\n`),
	regexp.MustCompile(`(净利润\(万元\).*?)\n`),
	regexp.MustCompile(`(净利润增长率\(%\).*?)\n`),
	regexp.MustCompile(`(营业总收入\(万元\).*?)\n`),
	regexp.MustCompile(`(营业总收入增长率\(%\).*?)\n`),
	regexp.MustCompile(`(加权净资产收益率\(%\).*?)\n`),
	regexp.MustCompile(`(资产负债比率\(%\).*?)\n`),
	regexp.MustCompile(`(净利润现金含量\(%\).*?)\n`),
}

func (d *Debug) Financial_analysis(s string) (map[string]interface{}, []error) {
	var err []error
	var res = map[string]interface{}{}
	var colIndex []string
	content := s
	r := regexp.MustCompile(`【主要财务指标】([\s\S]*?)【偿债能力指标】`).FindStringSubmatch(content)
	if r == nil || len(r) != 2 {
		err = append(err, fmt.Errorf("skeleton not found"))
		return nil, err
	}

	if strings.Count(r[1], "财务指标") == 1 {
		r = regexp.MustCompile(`【主要财务指标】[\s\S]*?(财务指标[\s\S]*?)【偿债能力指标】`).FindStringSubmatch(content)
		if r == nil || len(r) != 2 {

			err = append(err, fmt.Errorf("table not found fin1"))
			return nil, err
		}

		for _, reg := range regex {
			o := reg.FindStringSubmatch(r[1])
			if o != nil && len(o) == 2 {

				row := mergeRowsAndSplit(o[1])
				rowIndex := row[0]
				rowVals := row[1:]

				if rowIndex == "财务指标" {
					colIndex = rowVals
					continue
				}
				if _, hasKey := res[rowIndex]; !hasKey {
					res[rowIndex] = map[string]string{}

				}

				for i, v := range colIndex {
					if _, hasKey := res[rowIndex].(map[string]string)[v]; !hasKey {
						res[rowIndex].(map[string]string)[v] = rowVals[i]

					}
				}
			}
		}
	} else {
		r = regexp.MustCompile(`【主要财务指标】[\s\S]*?(财务指标[\s\S]*?)财务指标`).FindStringSubmatch(content)
		if r == nil || len(r) != 2 {
			err = append(err, fmt.Errorf("table not found fin2"))
			return nil, err
		}

		for _, reg := range regex {
			o := reg.FindStringSubmatch(r[1])
			if o != nil && len(o) == 2 {

				row := mergeRowsAndSplit(o[1])
				rowIndex := row[0]
				rowVals := row[1:]

				if rowIndex == "财务指标" {
					colIndex = rowVals
					continue
				}
				if _, hasKey := res[rowIndex]; !hasKey {
					res[rowIndex] = map[string]string{}

				}

				for i, v := range colIndex {
					if _, hasKey := res[rowIndex].(map[string]string)[v]; !hasKey {
						res[rowIndex].(map[string]string)[v] = rowVals[i]

					}
				}
			}

		}
		r = regexp.MustCompile(`【主要财务指标】[\s\S]*?财务指标[\s\S]*?(财务指标[\s\S]*?)【偿债能力指标】`).FindStringSubmatch(content)
		if r == nil || len(r) != 2 {
			err = append(err, fmt.Errorf("table not found fin3"))
			return nil, err
		}
		for _, reg := range regex {
			o := reg.FindStringSubmatch(r[1])
			if o != nil && len(o) == 2 {
				row := mergeRowsAndSplit(o[1])
				rowIndex := row[0]
				rowVals := row[1:]

				if rowIndex == "财务指标" {
					colIndex = rowVals
					continue
				}
				if _, hasKey := res[rowIndex]; !hasKey {
					res[rowIndex] = map[string]string{}
				}

				for i, v := range colIndex {
					if _, hasKey := res[rowIndex].(map[string]string)[v]; !hasKey {
						res[rowIndex].(map[string]string)[v] = rowVals[i]
					}
				}

			}
		}
	}
	/////////

	for k, v := range res {
		s, e := json.Marshal(v)
		if e != nil {
			err = append(err, e)
		}
		res[k] = string(s)
	}
	return res, err
}

func mergeRowsAndSplit(s string) []string {
	var res []string
	if !strings.Contains(s, "\n｜") {
		for _, c := range strings.Split(s, "｜") {
			c = strings.TrimSpace(c)
			if c != "" {
				res = append(res, c)
			}
		}
		return res
	}
	rows := strings.Split(s, "\n｜")
	for rid, row := range rows {

		if rid != 0 {
			for i, c := range strings.Split(row, "｜") {
				if c != "" {
					res[i] = res[i] + strings.TrimSpace(c)
				}
			}
		} else {
			for _, c := range strings.Split(row, "｜") {
				if c != "" {
					res = append(res, strings.TrimSpace(c))
				}
			}
		}

	}
	return res
}
