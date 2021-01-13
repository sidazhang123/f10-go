package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	f10-refined.LT/SA...:
	code:{
			fetchtime: ""
			field1:"{subField1:...}"     , ...
			...
			}
					||               with every code, put the field values together in a flattened way
					V               given the uniqueness of the jsonStr of the field (compat along with the "time" dimension)
	f10-acc.f10/f10补充...:
	code:{
		field1+subField1:[...],
		field1+subField2:[...],
	}
*/

func (s *Service) DoAll(code, start, end string, wg *sync.WaitGroup) (err []error) {
	var e error
	var lock sync.RWMutex
	for _, f := range []func(string, string, string) error{
		s.DoLT, s.DoFA, s.DoSA, s.DoOA,
	} {
		wg.Add(1)
		go func(f func(string, string, string) error) {

			e = f(code, start, end)
			if e != nil {
				lock.Lock()
				err = append(err, e)
				lock.Unlock()
			}
			wg.Done()
		}(f)

	}
	wg.Wait()
	return
}
func (s *Service) DoLT(code, start, end string) error {
	/* straightforwardly put json strings given a field together as an array
	if one is not in the array already
	*/

	collection := "latest_tips"
	var eList []string
	err, startT, endT := determineSpan(collection, start, end)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info(fmt.Sprintf("DoLT running, from %s to %s", startT.Format(common.TimestampLayout[:10]), endT.Format(common.TimestampLayout[:10])))

	err, cur := ReadMany(Params.RefinedDbName, collection, code, startT, endT)
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		var record map[string]interface{}
		err = cur.Decode(&record)
		if err != nil {
			return err
		}
		code := record["code"].(string)
		name := record["name"].(string)
		//for field -value
		for k, v := range record {
			if accCol, ok := AccFieldCollectionMap[k]; ok {
				err = AppendAndUpdateByField(accCol, code, name, k, rmCR(v.(string)))
				if err != nil {
					eList = append(eList, mark(err, 86, code, k).Error())
				}
			}
		}
	}

	err = ColToWinRepr[collection]()
	if err != nil {
		eList = append(eList, mark(err, 94).Error())
	}
	if len(eList) > 0 {
		return fmt.Errorf("DoLT===\n" + strings.Join(eList, "\n"))
	}
	return nil
}
func (s *Service) DoFA(code, start, end string) error {
	/*	fields of a records are fixed, accumulate new info by adding "date:value"
		pairs into the json string of the field
	*/

	collection := "financial_analysis"
	var eList []string
	err, startT, endT := determineSpan(collection, start, end)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info(fmt.Sprintf("DoFA running, from %s to %s", startT.Format(common.TimestampLayout[:10]), endT.Format(common.TimestampLayout[:10])))
	err, cur := ReadMany(Params.RefinedDbName, collection, code, startT, endT)
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		var record map[string]interface{}
		err = cur.Decode(&record)
		if err != nil {
			return err
		}
		code := record["code"].(string)
		name := record["name"].(string)
		//for field -value
		for k, v := range record {
			if _, ok := AccFieldCollectionMap[k]; ok {
				// new code -> insert; new field -> update; new date -> update;
				err = AddToJsonByFieldAndUpdate(Params.Win3Name, code, name, k, rmCR(v.(string)))
				if err != nil {
					eList = append(eList, mark(err, 132, code, k).Error())
				}
			}
		}
	}
	err = AddToJsonByFieldAndUpdate(Params.Win3Name, "", "", "", "")
	if err != nil {
		eList = append(eList, mark(err, 139).Error())
	}
	err = ColToWinRepr[collection]()
	if err != nil {
		eList = append(eList, mark(err, 143).Error())
	}

	if len(eList) > 0 {
		return fmt.Errorf("DoFA===\n" + strings.Join(eList, "\n"))
	}
	return nil
}
func (s *Service) DoOA(code, start, end string) error {

	dateRe := regexp.MustCompile(`[\d]{4}-[\d]{2}-[\d]{2}`)
	collection := "operational_analysis"
	var eList []string
	err, startT, endT := determineSpan(collection, start, end)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info(fmt.Sprintf("DoOA running, from %s to %s", startT.Format(common.TimestampLayout[:10]), endT.Format(common.TimestampLayout[:10])))
	err, cur := ReadMany(Params.RefinedDbName, collection, code, startT, endT)
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		var record map[string]interface{}
		err = cur.Decode(&record)
		if err != nil {
			return err
		}
		code := record["code"].(string)
		name := record["name"].(string)
		//for field -value
		for k, v := range record {
			if dateRe.MatchString(k) {
				err := AddFieldAndUpdate(Params.Win4Name, code, name, k, rmCR(v.(string)))
				if err != nil {
					eList = append(eList, mark(err, 179, code, k).Error())
				}
			}
		}
	}
	err = ColToWinRepr[collection]()
	if err != nil {
		eList = append(eList, mark(err, 186).Error())
	}
	if len(eList) > 0 {
		return fmt.Errorf("DoOA===\n" + strings.Join(eList, "\n"))
	}
	return nil
}
func (s *Service) DoSA(code, start, end string) error {
	var eList []string
	collection := "shareholder_analysis"
	extractNum := regexp.MustCompile(`流通占比表([\d]+)`)
	//股东变化表 f10b
	err, startT, endT := determineSpan(collection, start, end)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info(fmt.Sprintf("DoSA running, from %s to %s", startT.Format(common.TimestampLayout[:10]), endT.Format(common.TimestampLayout[:10])))
	err, cur := ReadMany(Params.RefinedDbName, collection, code, startT, endT)
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		var record map[string]interface{}
		err = cur.Decode(&record)
		if err != nil {
			eList = append(eList, mark(err, 212).Error())
		}
		code := record["code"].(string)
		name := record["name"].(string)
		//fetchtime:=record["fetchtime"].(primitive.DateTime).Time().Format(common.TimestampLayout[:10])

		if _, ok := record["流通占比表0"]; !ok {
			//fetchtime:=record["fetchtime"].(primitive.DateTime).Time().Format(common.TimestampLayout[:10])
			//eList = append(eList, fmt.Sprintf("code %s fetchtime %s: acc record present but no 流通占比表0 in the original", code, fetchtime))
			continue
		}
		// deal with the shareholder info
		for k, v := range record {
			// only 控股比例
			if accCol, ok := AccFieldCollectionMap[k]; ok {
				err = AppendAndUpdateByField(accCol, code, name, k, rmCR(v.(string)))
				if err != nil {
					eList = append(eList, mark(err, 229, code, k).Error())
				}
			}
		}
		// agg the tables:: {code,name,updatetime,"tables":[oldtab -> newtab]} does not need additional repr()
		err, m := FindOne(Params.AccumulatedDbName, Params.Win5Name, "code", code)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// if not present, insert it while reverting the table numbers, continue
				max := -1
				for k := range record {
					r := extractNum.FindStringSubmatch(k)
					if len(r) == 2 {
						curN, _ := strconv.Atoi(r[1])
						if curN > max {
							max = curN
						}
					}
				}
				// insert when there is one table at least
				if max > -1 {
					// update requires a record that meets the filter in place in advance
					err = InsertOne(Params.AccumulatedDbName, Params.Win5Name, map[string]interface{}{"code": code})
					if err != nil {
						eList = append(eList, mark(err, 253, code).Error())
					}
					// revert the table numbers
					for i := max; i > -1; i-- {
						k := fmt.Sprintf("流通占比表%d", i)
						err = AppendAndUpdateByField(Params.Win5Name, code, name, "tables", record[k].(string))
						if err != nil {
							eList = append(eList, mark(err, 260, code).Error())
						}
					}
				}
				continue
			} else {
				eList = append(eList, err.Error())
			}
		}

		// from now on, there has to be the record of the code with no less than one table
		var tab map[string]interface{}
		err = json.Unmarshal([]byte(record["流通占比表0"].(string)), &tab)
		if err != nil {
			eList = append(eList, mark(err, 274).Error())
		}
		tNew, err := time.Parse(common.TimestampLayout[2:10], tab["date"].(string))
		if err != nil {
			eList = append(eList, mark(err, 278).Error())
		}
		// get the date from f10-acc as well.  [oldtab -> newtab]
		tabsRec := m["tables"].(primitive.A)
		tabs := make([]string, len(tabsRec))
		for i, t := range tabsRec {
			tabs[i] = t.(string)
		}

		err = json.Unmarshal([]byte(tabs[len(tabs)-1]), &tab)
		if err != nil {
			eList = append(eList, mark(err, 289).Error())
		}
		tOld, err := time.Parse(common.TimestampLayout[2:10], tab["date"].(string))
		if tNew.After(tOld) {
			// if more recent, append the table
			err = AppendAndUpdateByField(Params.Win5Name, code, name, "tables", record["流通占比表0"].(string))
			if err != nil {
				eList = append(eList, mark(err, 296).Error())
			}
		}
	}
	err = ColToWinRepr[collection]()
	if err != nil {
		eList = append(eList, mark(err, 302).Error())
	}
	if len(eList) > 0 {
		return fmt.Errorf("DoSA===\n" + strings.Join(eList, "\n"))
	}
	return nil
}

func determineSpan(collection, start, end string) (error, time.Time, time.Time) {
	var startT, endT time.Time
	var err error

	if len(start) != 10 {
		err, startT = FindEndDate(collection, false)
		if err != nil {
			return err, startT, endT
		}
	} else {
		startT, err = time.Parse(common.TimestampLayout[:10], start)
		if err != nil {
			return err, startT, endT
		}
	}
	if len(end) != 10 {
		err, endT = FindEndDate(collection, true)
		if err != nil {
			return err, startT, endT.AddDate(0, 0, 1)
		}
	} else {
		endT, err = time.Parse(common.TimestampLayout[:10], end)
		if err != nil {
			return err, startT, endT
		}
	}
	if !endT.After(startT) {
		return fmt.Errorf("[determineTimespan] %s=>%s", startT.Format(common.TimestampLayout[:10]), endT.Format(common.TimestampLayout[:10])), time.Time{}, time.Time{}
	}
	return nil, startT, endT
}

func rmCR(s string) string {
	return strings.ReplaceAll(s, "\\r", "")
}
