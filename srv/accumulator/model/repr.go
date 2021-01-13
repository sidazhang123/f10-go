package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"sort"
	"strings"
	"time"
)

var fieldFlattenMap = map[string]func(string, interface{}, map[string]interface{}) error{
	"公积":     _flatten.flattenConcatName,
	"未分":     _flatten.flattenConcatName,
	"质押":     _flatten.flattenDeposition,
	"发行前限售":  _flatten.flattenConcatName,
	"股改限售":   _flatten.flattenConcatName,
	"增发A股":   _flatten.flattenConcatName,
	"激励限售":   _flatten.flattenConcatName,
	"承诺到期":   _flatten.flattenConcatName,
	"最新公告":   _flatten.flattenRmvMsg,
	"最新报道":   _flatten.flattenRmvMsg,
	"业绩预告":   _flatten.flattenRmvMsg,
	"特别处理":   _flatten.flattenRmvMsg,
	"最新提醒增发": _flatten.flattenConcatSubfields,
	"股东户数变化": _flatten.flattenConcatName,
	"股东控股":   _flatten.flattenConcatKV,
}
var _getLatestN = getLatestN{}
var _flatten = flatten{}

/*	all records in f10-acc.winName are unique given their code
	this method makes the repr json string, ready to consume by the qt project, and insert it back to the record
	so the uniqueness of the code is kept in this collection
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	DoXX accumulates refined db records and puts in acc db and reprXX updates their repr fields
	win1:LT,SA
	win2:LT
	win3:FA
	win4:OA
	win5:SA
*/
var ColToWinRepr = map[string]func() error{
	"latest_tips": reprWin12, "financial_analysis": reprWin3, "operational_analysis": reprWin4, "shareholder_analysis": reprWin15,
}

var winToReprMap = map[string]map[string]string{}
var winToReprStr = map[string]string{}

func (s *Service) ReprAll() (err []error) {
	for _, f := range ColToWinRepr {
		e := f()
		if e != nil {
			err = append(err, e)
		}
	}
	return
}
func (s *Service) GetRepr(win string) string {
	return winToReprStr[win]
}
func makeReprMap(win, name, repr string) {
	winToReprMap[win][name] = repr
}
func makeReprStr(win string) {
	b, _ := json.Marshal(winToReprMap[win])
	winToReprStr[win] = string(b)
}
func initReprJson() error {
	init := true
	for _, v := range winToReprMap {
		if len(v) > 0 {
			init = false
		}
	}
	if init {
		for _, winName := range []string{Params.Win1Name, Params.Win2Name, Params.Win3Name, Params.Win4Name, Params.Win5Name} {
			winToReprMap[winName] = map[string]string{}
			err, cur := ReadMany(Params.AccumulatedDbName, winName, "", time.Time{}, time.Time{})
			if err != nil {
				return err
			}
			for cur.Next(context.TODO()) {
				var record map[string]interface{}
				err := cur.Decode(&record)
				if err != nil {
					return err
				}
				name, repr := record["name"].(string), record["repr"].(string)
				winToReprMap[winName][name] = repr
			}
			makeReprStr(winName)
		}
	}

	return nil
}
func reprWin15() error {
	err, cur := ReadMany(Params.AccumulatedDbName, Params.Win5Name, "", time.Time{}, time.Time{})
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		var rec map[string]interface{}
		err := cur.Decode(&rec)
		if err != nil {
			return err
		}
		code := rec["code"].(string)
		if v, ok := rec["tables"]; ok {
			jsonArray := v.(primitive.A)
			jsonStrArray := make([]string, len(jsonArray))
			for i, t := range jsonArray {
				// revert; new to old
				jsonStrArray[len(jsonArray)-1-i] = t.(string)
			}
			b, err := json.Marshal(map[string]interface{}{"tables": jsonStrArray,
				"code": code, "name": rec["name"].(string), "updatetime": rec["updatetime"].(string)})
			if err != nil {
				return err
			}
			err = UpdateOne(Params.AccumulatedDbName, Params.Win5Name, code, bson.M{"repr": string(b)})
			if err != nil {
				return err
			}
			makeReprMap(Params.Win5Name, rec["name"].(string), string(b))
		} else {
			return fmt.Errorf("[reprWin5] %s doesn't have 'tables'", code)
		}
	}
	//
	k := "股东控股"
	err, cur = ReadMany(Params.AccumulatedDbName, Params.Win1Name, "", time.Time{}, time.Time{})
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		var rec map[string]interface{}
		err := cur.Decode(&rec)
		if err != nil {
			return err
		}
		if sh, ok := rec[k]; ok {
			code := rec["code"].(string)
			if v, ok := rec["repr"]; ok {
				var repr map[string]interface{}
				err = json.Unmarshal([]byte(v.(string)), &repr)
				if err != nil {
					return err
				}
				err = fieldFlattenMap[k](k, sh, repr)
				if err != nil {
					return err
				}
				repr["name"] = rec["name"].(string)
				repr["updatetime"] = rec["updatetime"].(string)
				b, err := json.Marshal(repr)
				if err != nil {
					return err
				}
				err = UpdateOne(Params.AccumulatedDbName, Params.Win1Name, code, bson.M{"repr": string(b)})
				if err != nil {
					return err
				}
				makeReprMap(Params.Win1Name, rec["name"].(string), string(b))
			}
		}

	}
	makeReprStr(Params.Win1Name)
	makeReprStr(Params.Win5Name)
	return nil
}

func reprWin4() error {
	err, cur := ReadMany(Params.AccumulatedDbName, Params.Win4Name, "", time.Time{}, time.Time{})
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		var record map[string]interface{}
		err := cur.Decode(&record)
		if err != nil {
			return err
		}
		code := record["code"].(string)
		repr := map[string]interface{}{"code": code, "name": record["name"].(string), "updatetime": record["updatetime"].(string)}
		m := map[string]string{}
		for k, v := range record {
			if _getLatestN.matchExact(`[\d]{4}-[\d]{2}-[\d]{2}`)(k) {
				m[k] = v.(string)
			}
		}
		err, m = _getLatestN.getLatestN(m, Params.Win4Capacity, _getLatestN.pass)
		if err != nil {
			return err
		}
		for k, v := range m {
			repr[k] = v
		}
		b, err := json.Marshal(repr)
		if err != nil {
			return err
		}
		err = UpdateOne(Params.AccumulatedDbName, Params.Win4Name, code, bson.M{"repr": string(b)})
		if err != nil {
			return err
		}
		makeReprMap(Params.Win4Name, record["name"].(string), string(b))
	}
	makeReprStr(Params.Win4Name)
	return nil
}

func reprWin3() error {
	/*
		f10-acc
		code:{
			tabField1:{date:,...},
			tabField2:{date:,...}, // struct in the accumulation collection
			repr:"{code,name,updatetime,YTab:{tabField1:{date:,...},...},QTab:{tabField2:{date:,...},...}}"
		}
		The qt client gets "repr", unmarshal it, matches the grid
	*/
	err, cur := ReadMany(Params.AccumulatedDbName, Params.Win3Name, "", time.Time{}, time.Time{})
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		var record map[string]interface{}
		err := cur.Decode(&record)
		if err != nil {
			return err
		}
		code := record["code"].(string)
		repr := map[string]interface{}{"code": code, "name": record["name"].(string), "updatetime": record["updatetime"].(string)}
		// flatten the field value with varied methods
		YTab, QTab := map[string]interface{}{}, map[string]interface{}{}
		for field, v := range record {
			// if field in the seq, unmarshal the v,order by date, and pick up to five latest "-12-" as the YTab and five latest items as the QTab
			if _, ok := AccFieldCollectionMap[field]; ok {
				var dateMsg map[string]string
				err := json.Unmarshal([]byte(v.(string)), &dateMsg)
				if err != nil {
					return err
				}
				err, YearlyRow := _getLatestN.getLatestN(dateMsg, Params.Win3Capacity, _getLatestN.contains("-12-"))
				if err != nil {
					return err
				}
				err, QuaterlyRow := _getLatestN.getLatestN(dateMsg, Params.Win3Capacity, _getLatestN.pass)
				if err != nil {
					return err
				}
				YTab[field], QTab[field] = YearlyRow, QuaterlyRow
			}
		}
		yb, err := json.Marshal(YTab)
		if err != nil {
			return err
		}
		qb, err := json.Marshal(QTab)
		if err != nil {
			return err
		}
		repr["YTab"], repr["QTab"] = string(yb), string(qb)
		b, err := json.Marshal(repr)
		if err != nil {
			return err
		}
		err = UpdateOne(Params.AccumulatedDbName, Params.Win3Name, code, bson.M{"repr": string(b)})
		if err != nil {
			return err
		}
		makeReprMap(Params.Win3Name, record["name"].(string), string(b))
	}
	makeReprStr(Params.Win3Name)
	return nil
}

func reprWin12() error {
	/*
		f10-acc
		code:{
			field1+subField1:[...],
			field1+subField2:[...],
			repr:"{code,name,updatetime,field1+subField1:[...],field1+subField2:[...],...}"
		}
		The qt client gets "repr", unmarshal it, matches the grid
	*/
	for _, winName := range []string{Params.Win1Name, Params.Win2Name} {
		err, cur := ReadMany(Params.AccumulatedDbName, winName, "", time.Time{}, time.Time{})
		if err != nil {
			return err
		}
		for cur.Next(context.TODO()) {
			var record map[string]interface{}
			err := cur.Decode(&record)
			if err != nil {
				return err
			}
			code := record["code"].(string)
			repr := map[string]interface{}{"code": code, "name": record["name"].(string), "updatetime": record["updatetime"].(string)}
			// flatten the field value with varied methods
			for field, v := range record {
				if f, ok := fieldFlattenMap[field]; ok {
					err := f(field, v, repr)
					if err != nil {
						return err
					}
				}
			}
			b, err := json.Marshal(repr)
			if err != nil {
				return err
			}
			err = UpdateOne(Params.AccumulatedDbName, winName, code, bson.M{"repr": string(b)})
			if err != nil {
				return err
			}
			makeReprMap(winName, record["name"].(string), string(b))
		}
		makeReprStr(winName)
	}
	return nil
}

type getLatestN struct{}

func (g *getLatestN) contains(contains string) func(k string) bool {
	return func(k string) bool {
		return len(k) == 10 && strings.Contains(k, contains)
	}
}
func (g *getLatestN) matchExact(re string) func(k string) bool {
	return func(k string) bool {
		if !strings.HasPrefix(re, "^") {
			re = "^" + re
		}
		if !strings.HasSuffix(re, "$") {
			re += "$"
		}
		return regexp.MustCompile(re).MatchString(k)
	}
}
func (g *getLatestN) pass(k string) bool { return len(k) == 10 }
func (g *getLatestN) getLatestN(m map[string]string, n int, checker func(string) bool) (error, map[string]string) {
	res := map[string]string{}
	c := 0
	for k := range m {
		if checker(k) {
			c += 1
		}
	}
	if c <= n {
		for k, v := range m {
			if checker(k) {
				res[k] = v
			}
		}
		return nil, res
	}
	var kL []time.Time
	for k := range m {
		// a prerequisite is that k is valid in the date format
		if checker(k) {
			t, err := time.Parse(common.TimestampLayout[:10], k)
			if err != nil {
				return fmt.Errorf("[getLatestN]" + err.Error()), nil
			}
			kL = append(kL, t)
		}
	}
	sort.SliceStable(kL, func(i, j int) bool {
		return kL[i].After(kL[j])
	})

	for _, t := range kL[:n] {
		k := t.Format(common.TimestampLayout[:10])
		res[k] = m[k]
	}
	return nil, res
}

type flatten struct{}

//"L1+L2 : L2 value"
func (f *flatten) flattenConcatName(k string, v interface{}, m map[string]interface{}) error {
	items, err := jsonStrArray2MapArray(v)
	if err != nil {
		return err
	}

	for _, i := range items {
		for subField, val := range i {
			name := k + subField
			if _, ok := m[name]; ok {
				m[name] = append(m[name].([]string), val)
			} else {
				m[name] = []string{val}
			}
		}
	}
	return nil
}

// L1+股东: L2
// L1+占比: L2 value
func (f *flatten) flattenDeposition(k string, v interface{}, m map[string]interface{}) error {
	items, err := jsonStrArray2MapArray(v)
	if err != nil {
		return err
	}

	for _, i := range items {
		for entityName, val := range i {
			// 股东
			name := k + "股东"
			if _, ok := m[name]; ok {
				m[name] = append(m[name].([]string), entityName)
			} else {
				m[name] = []string{entityName}
			}
			// 占比
			name = k + "占比"
			if _, ok := m[name]; ok {
				m[name] = append(m[name].([]string), val)
			} else {
				m[name] = []string{val}
			}

		}
	}
	return nil
}

//"L1+L2 : L2 value"
var rmvSpace = regexp.MustCompile(`\s+`)

func (f *flatten) flattenConcatKV(k string, v interface{}, m map[string]interface{}) error {
	items, err := jsonStrArray2MapArray(v)
	if err != nil {
		return err
	}
	var s []string
	for _, i := range items {
		var tmp []string
		for subField, val := range i {
			tmp = append(tmp, rmvSpace.ReplaceAllString(subField, "")+" "+rmvSpace.ReplaceAllString(val, ""))
		}
		s = append(s, strings.Join(tmp, "\n"))
	}
	m[k] = s
	return nil
}

// convert []interface{} (underlying type is string) to []map[string]string
func jsonStrArray2MapArray(v interface{}) ([]map[string]string, error) {
	jsonArray := v.(primitive.A)
	jsonStrArray := make([]string, len(jsonArray))
	for i, t := range jsonArray {
		jsonStrArray[i] = t.(string)
	}

	var items []map[string]string
	for _, jsonItem := range jsonStrArray {
		var item map[string]string
		err := json.Unmarshal([]byte(jsonItem), &item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

//"L1 : join(L2 values)"
func (f *flatten) flattenConcatSubfields(k string, v interface{}, m map[string]interface{}) error {
	items, err := jsonStrArray2MapArray(v)
	if err != nil {
		return err
	}
	for _, i := range items {
		name := k
		toAppend := ""
		for _, val := range i {
			toAppend += val + "\n"
		}
		toAppend = strings.Trim(toAppend, "\n")
		if _, ok := m[name]; ok {
			m[name] = append(m[name].([]string), toAppend)
		} else {
			m[name] = []string{toAppend}
		}
	}
	return nil
}

//"L1 : L2 value"
func (f *flatten) flattenRmvMsg(k string, v interface{}, m map[string]interface{}) error {
	items, err := jsonStrArray2MapArray(v)
	if err != nil {
		return err
	}
	for _, i := range items {
		name := k
		val := i["msg"]
		if _, ok := m[name]; ok {
			m[name] = append(m[name].([]string), val)
		} else {
			m[name] = []string{val}
		}
	}
	return nil
}
