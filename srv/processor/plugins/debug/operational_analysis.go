package debug

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"regexp"
	"strings"
)

// only need two columns of percentage numbers so it's unnecessary to care about a number which displays in multiple rows
func (d *Debug) Operational_analysis(s string) (map[string]interface{}, []error) {
	//find dates of tables
	r := regexp.MustCompile(`【2\.主营[^】]+】[^】]*?\n[^【]*?【截止日期】`).FindStringSubmatch(s)
	if r == nil {
		return nil, []error{fmt.Errorf("1 tables not found")}
	}
	var res [][]string

	//extract blocks, components of the tables, divided by "─────────────────────────────────────"
	r = regexp.MustCompile(`【2\.[^】]+】[^】]*?\n[^【]*?(【截止日期】[\s\S]+)【3\.`).FindStringSubmatch(s)
	if r == nil {
		r = regexp.MustCompile(`【2\.[^】]+】[^】]*?\n[^【]*?(【截止日期】[\s\S]+)`).FindStringSubmatch(s)
		if r == nil {
			return nil, []error{fmt.Errorf("2 tables not found")}
		}
	}
	s = r[1]

	blocks := strings.Split(s, "─────────────────────────────────────")

	for _, i := range blocks {
		// for blocks of content
		if !strings.Contains(i, "【截止日期】") {
			i = strings.TrimSpace(i)
			//split a block into strips which group concrete contents
			//a strip:
			//注射用复方二氯醋酸二异丙 2214.80 2126.02 95.99 20.57\n
			//胺(产品)

			stripList := splitStrips(i)

			for _, strip := range stripList {
				tags := map[string]string{"(行": "(行业)", "业)": "(行业)",
					"(产": "(产品)", "品)": "(产品)",
					"(地": "(地区)", "区)": "(地区)"}
				for tagK, tagV := range tags {
					//verify that a strip has a header
					if strings.Contains(strip, tagK) {
						//dispose unimportant digits and header chars in the second line
						strip = strings.Split(strip, "\n")[0]
						cellsOfLine := strings.Split(strip, " ")
						if strings.Contains(cellsOfLine[0], tagV) {
							cellsOfLine[0] = strings.Replace(cellsOfLine[0], tagV, "", -1)
						}
						if strings.Contains(cellsOfLine[0], tagK) {
							cellsOfLine[0] = strings.Replace(cellsOfLine[0], tagK, "", -1)
						}
						header := strings.Replace(cellsOfLine[0], tagK, "", -1) + tagV
						res = append(res, []string{header, cellsOfLine[3], cellsOfLine[4]})
						break
					}
				}
			}
		} else {
			r = regexp.MustCompile(`【截止日期】.*?([\d]{4}-[\d]{2}-[\d]{2})`).FindStringSubmatch(i)
			if r == nil {
				return nil, []error{fmt.Errorf("date failed to match the table content")}
			}

			res = append(res, []string{r[1]})
		}
	}
	//for _, v := range res {
	//	fmt.Printf("%+v\n", v)
	//}
	//fmt.Printf("resOuter %+v\n",res)
	// res [][]string {{"2019-01-01"},{"新材料(行业)", "10.51", "12.32"}, {"2019-05-01"}, ...}
	var toPrint [][]string
	var toPrintPlanB [][]string
	var finRes = map[string]interface{}{}
	dateCount := 0

	tableString := &strings.Builder{}
	var table *tablewriter.Table
	preDateHeader := ""
	for ind, i := range res {
		switch {
		case len(i) == 1:
			if preDateHeader == "" {
				preDateHeader = i[0]
			}
			dateCount += 1
			if dateCount > 1 {
				//fmt.Printf(">1 date %s ;; data %+v\n",preDateHeader,toPrint)
				tp := printTable(table, tableString, preDateHeader, toPrint)
				if strings.TrimSpace(tp) == "" {
					finRes[preDateHeader] = printTable(table, tableString, preDateHeader, toPrintPlanB)
				} else {
					finRes[preDateHeader] = tp
				}
				toPrint = [][]string{}
				toPrintPlanB = [][]string{}
				preDateHeader = i[0]
			}
		case !strings.HasPrefix(i[0], "合计") && strings.Contains(i[0], "(产品)"):
			toPrint = append(toPrint, []string{i[0], dash(i[1]), dash(i[2])})
			//fmt.Printf("case1 date %s ;; data %+v\n",dateHeader,toPrint)
		case !strings.HasPrefix(i[0], "合计") && strings.Contains(i[0], "(行业)"):
			toPrintPlanB = append(toPrintPlanB, []string{i[0], dash(i[1]), dash(i[2])})
			//fmt.Printf("case2 date %s ;; data %+v\n",dateHeader,toPrintPlanB)
		}
		if ind == len(res)-1 {
			tp := printTable(table, tableString, preDateHeader, toPrint)
			if strings.TrimSpace(tp) == "" {
				finRes[preDateHeader] = printTable(table, tableString, preDateHeader, toPrintPlanB)
			} else {
				finRes[preDateHeader] = tp
			}
		}
	}
	return finRes, nil
}

func dash(s string) string {
	if strings.TrimSpace(s) == "-" {
		return s
	} else {
		return s + " %"
	}
}
func printTable(table *tablewriter.Table, tableString *strings.Builder, dateHeader string, toPrint [][]string) string {
	//fmt.Printf("%s,,,%+v\n", dateHeader, toPrint)
	if len(toPrint) == 0 {
		return ""
	} else {
		for _, v := range toPrint {
			if len(v) == 0 {
				return ""
			}
		}
	}
	table = tablewriter.NewWriter(tableString)
	table.SetHeader([]string{dateHeader, "毛利", "占营收"})
	for _, row := range toPrint {
		table.Append(row)
	}
	table.Render()
	tp := tableString.String()

	tableString.Reset()

	return tp
}

//split a block into strips which group concrete contents
//a strip:
//注射用复方二氯醋酸二异丙 2214.80 2126.02 95.99 20.57\n
//胺(产品)
func splitStrips(block string) []string {
	block = regexp.MustCompile("[ ]+").ReplaceAllString(block, " ")
	lines := strings.Split(block, "\n")
	var res []string
	concreteLine := 0
	tmp := ""
	for _, line := range lines {
		if len(strings.Split(line, " ")) != 5 {
			tmp = tmp + line + "\n"
		} else {
			if tmp != "" && concreteLine >= 1 {
				res = append(res, strings.TrimSpace(tmp))
				tmp = line + "\n"
				concreteLine = 0

			} else {
				if concreteLine == 0 {
					tmp = line + "\n"
				}
			}
		}
		concreteLine += 1
	}
	if concreteLine == 1 {
		res = append(res, tmp)
	}
	return res
}
