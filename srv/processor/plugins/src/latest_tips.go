package src

import (
	"encoding/json"
	"regexp"
	"strconv"
)

// the concrete type of the return type map[string]interface{} is map[string]string
// the interface{} is used for adding time.Time before sending to db
func Latest_tips(s string) (map[string]interface{}, []error) {
	filter := regexp.MustCompile(`[ ｜─├┤└┘/\n\r\t]+`)
	var res = map[string]interface{}{}
	var errList []error
	//([\d]{4}-[\d]{2}-[\d]{2})每股资本公积[：:]([^营]+)营业收入\(万元\)[：:]([\d\s.｜-]+)(?:同比增([^%]+%))?
	r := regexp.MustCompile(`([\d]{4}-[\d]{2}-[\d]{2})[^每]*每股资本公积[：:]([^营]+)营业收入\(万元\)[：:]([\d\s.｜-]+)(?:同比增([^%]+%))?`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[4]string{"更新日期", "每股公积", "营收万元", "同比增"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["公积"] = leafItem
	}

	//([\d]{4}-[\d]{2}-[\d]{2})每股未分利润[：:]([^净]+)净利润\(万元\)[：:]([\d\s.｜-]+)(?:同比增([^%]+%))?
	r = regexp.MustCompile(`([\d]{4}-[\d]{2}-[\d]{2})[^每]*每股未分利润[：:]([^净]+)净利润\(万元\)[：:]([\d\s.｜-]+)(?:同比增([^%]+%))?`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[4]string{"更新日期", "每股未分", "净利万元", "同比增"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["未分"] = leafItem
	}
	//【质押占比】
	rr := regexp.MustCompile(`【质押占比】[：:]([^\d]+[^%]+%)`).FindAllStringSubmatch(s, -1)
	if rr != nil {
		tmp := map[string]string{}
		for _, line := range rr {
			itemInLine := filter.ReplaceAllString(line[1], "")
			t := regexp.MustCompile(`控股股东(.+)累计质押占其持股比例([\d.]+%)`).FindStringSubmatch(itemInLine)
			if t != nil {
				tmp[t[1]] = t[2]
			}
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["质押"] = leafItem
	}
	//发行前股份限售流通\(([\d-]{10})\)[：:]([^股]+)万?股
	r = regexp.MustCompile(`发行前股份限售流通\(([\d-]{10})\)[：:]([^股]+)万?股`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[2]string{"日期", "股数"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["发行前限售"] = leafItem
	}
	//★股改限售流通\(([\d-]{10})\)[:：]{1}([^股]+?)万?股\(详见股本结构\)｜<br>
	r = regexp.MustCompile(`★股改限售流通\(([\d-]{10})\)[:：]([^股]+?)万?股\(详见股本结构\)`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[2]string{"日期", "股数"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["股改限售"] = leafItem
	}
	//增发A股法人配售上市\(([\d-]{10})\)[：:]([^股]+)万?股
	r = regexp.MustCompile(`增发A股法人配售上市\(([\d-]{10})\)[：:]([^股]+)万?股`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[2]string{"日期", "股数"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["增发A股"] = leafItem
	}
	//股权激励限售流通\(([\d-]{10})\)[：:]([^股]+)万?股
	r = regexp.MustCompile(`股权激励限售流通\(([\d-]{10})\)[：:]([^股]+)万?股`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[2]string{"日期", "股数"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["激励限售"] = leafItem
	}
	//无限售不减持承诺到期\(([\d-]{10})\)[：:]([^股]+)万?股
	r = regexp.MustCompile(`无限售不减持承诺到期\(([\d-]{10})\)[：:]([^股]+)万?股`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[2]string{"日期", "股数"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["承诺到期"] = leafItem
	}
	//★最新公告[:：]?(.*?)[★【]{1}
	r = regexp.MustCompile(`★最新公告[:：]?([^★【]*?)[★【]`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[1]string{"msg"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["最新公告"] = leafItem
	}
	//★最新报道[:：]?(.*?)[★【]{1}
	r = regexp.MustCompile(`★最新报道[:：]?([^★【]*?)[★【]`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[1]string{"msg"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["最新报道"] = leafItem
	}
	//【业绩预告】[:：]?(.*?)[【★]{1}
	r = regexp.MustCompile(`【业绩预告】[:：]?([^【★]*?)[【★]`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[1]string{"msg"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["业绩预告"] = leafItem
	}
	//★特别处理[:：]?(.*?)｜<br>
	r = regexp.MustCompile(`★特别处理[:：]?([^【★]*?)[【★]`).FindStringSubmatch(s)
	if r != nil {
		tmp := map[string]string{}
		for ind, i := range r[1:] {
			tmp[[1]string{"msg"}[ind]] = i
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["特别处理"] = leafItem
	}
	//增发】[:：]?(.*?)[【★]{1}
	rr = regexp.MustCompile(`增发】[:：]?([^【★]*?)[【★]`).FindAllStringSubmatch(s, -1)
	if rr != nil {
		tmp := map[string]string{}
		for ind, i := range rr {
			tmp[strconv.Itoa(ind)] = i[1]
		}
		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["最新提醒增发"] = leafItem
	}
	//★股东户数变化[:：]?(.*?)[【★]{1}
	//截止([\d]{4}-[\d]{2}-[\d]{2}),公司股东数([\d]+),比上期\(([\d]{4}-[\d]{2}-[\d]{2})\)([\D]{2})([\d]+)户,幅度([^%]+%)\(详见主力追踪\)  //增减与幅度符号偶尔不一致
	//截止([\d]{4}-[\d]{2}-[\d]{2}),公司股东数([\d]+),比上期\(([\d]{4}-[\d]{2}-[\d]{2})\)(不变),幅度[^%]+%\(详见主力追踪\)   //不变
	//截止([\d]{4}-[\d]{2}-[\d]{2}),公司股东数([\d]+),幅度-%\(详见主力追踪\)      //不变
	//"股东户数变化":{"本期":"","股东数":"","上期":"", "增减":"","户数":"","幅度":""}
	r = regexp.MustCompile(`★股东户数变化[:：]?([^【★]*?)[【★]`).FindStringSubmatch(s)

	if r != nil {
		content := filter.ReplaceAllString(r[1], "")
		tmp := map[string]string{}
		r1 := regexp.MustCompile(`截止([\d]{4}-[\d]{2}-[\d]{2}),公司股东数([\d]+),比上期\(([\d]{4}-[\d]{2}-[\d]{2})\)([\D]{2})([\d]+)户,幅度([^%]+%)\(详见主力追踪\)`).FindStringSubmatch(content)
		if r1 != nil {
			for ind, i := range r1[1:] {
				tmp[[6]string{"本期时间", "本期股东数", "上期时间", "增减", "户数", "幅度"}[ind]] = i
			}
		} else {
			r1 = regexp.MustCompile(`截止([\d]{4}-[\d]{2}-[\d]{2}),公司股东数([\d]+),比上期\(([\d]{4}-[\d]{2}-[\d]{2})\)(不变),幅度[^%]+%\(详见主力追踪\)`).FindStringSubmatch(content)
			if r1 != nil {
				for ind, i := range r1[1:] {
					tmp[[4]string{"本期时间", "本期股东数", "上期时间", "增减"}[ind]] = i
				}
			} else {
				r1 = regexp.MustCompile(`截止([\d]{4}-[\d]{2}-[\d]{2}),公司股东数([\d]+),幅度-%\(详见主力追踪\)`).FindStringSubmatch(content)
				if r1 != nil {
					for ind, i := range r1[1:] {
						tmp[[2]string{"本期时间", "本期股东数"}[ind]] = i
					}
					tmp["增减"] = "不变"
				}
			}
		}

		leafItem, err := marshalLatestTips(tmp)
		if err != nil {
			errList = append(errList, err)
		}
		res["股东户数变化"] = leafItem
	}

	//
	//latest_tips = {
	//"公积":{"更新日期":"", "每股公积":"", "营收万元":"", "同比增":""},
	//"未分":{"更新日期":"", "每股未分":"", "净利万元":"", "同比增":""},
	//"质押":{"股东":"", "占比":""},
	//"发行前限售":{"日期":"", "股数":""}, //股数 可带 万
	//"股改限售":{"日期":"", "股数":""}, //股数 可带 万
	//"增发A股":{"日期":"", "股数":""}, //股数 可带 万
	//"激励限售":{"日期":"","股数":""}, //股数 可带 万
	//"承诺到期":{"日期":"","股数":""}, //股数 可带 万
	//"最新公告":{"msg":""},
	//"最新报道":{"msg":""},
	//"业绩预告":{"msg":""},
	//"特别处理":{"msg":""},
	//"最新提醒增发":{"0":""，"1":""},
	//"股东户数变化":{"本期":"","股东数":"","上期":"", "增减":"","户数":"","幅度":""}
	//}
	return res, nil
}

func marshalLatestTips(s map[string]string) (string, error) {
	for k, v := range s {
		v = regexp.MustCompile("[ ｜─├┤└┘/\n\r\t]+").ReplaceAllString(v, "")
		s[k] = v
	}
	leafItem, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(leafItem), nil
}
