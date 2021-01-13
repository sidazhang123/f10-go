package model

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	proto "github.com/sidazhang123/f10-go/srv/processor/proto/processor"
	"strings"
	"time"
)

func (s *Service) Refine(req *proto.Request) error {

	//default
	var srcTime time.Time
	var err error
	var flagNames []string
	if len(req.Date) > 0 {
		srcTime, err = time.Parse(common.TimestampLayout, req.Date+"T11:04:05.000Z")
		if err != nil {
			return err
		}
	} else {
		srcTime = time.Now().UTC().Add(8 * time.Hour)
	}
	if len(req.FlagName) > 0 {
		if _, ok := common.FlagNameToCollSuffix[req.FlagName]; ok {
			flagNames = []string{req.FlagName}
		} else {
			return fmt.Errorf("req.FlagName incorrect")
		}
	} else {
		flagNames = strings.Split(Params.Flags, ",")
	}

	for _, flagName := range flagNames {
		if flagName == "" {
			continue
		}
		engFlagName := common.FlagNameToCollSuffix[flagName]
		//get raw list of one flagName/collection
		err, cur := FindByDate(engFlagName, srcTime)
		if err != nil {
			return fmt.Errorf("failed to FindByDate from %s in raw DB\n%s", engFlagName, err.Error())
		}
		//process

		var refinedRes []map[string]interface{}
		for cur.Next(context.TODO()) {
			var raw Raw
			var res map[string]interface{}
			var errList []error
			err := cur.Decode(&raw)

			if err != nil {
				log.Error(fmt.Sprintf("failed to Decode cur to Raw\n%s", err.Error()))
			}
			//log2.Info("get Raw " + raw.Code)
			if Params.PluginLevel == "prod" {
				res, errList = s.CallSoPlugin(engFlagName, raw.Body)
				if errList != nil && len(errList) > 0 {
					eStr := ""
					for _, e := range errList {
						log.Error(e.Error())
						eStr += e.Error()
					}
					log.Error(fmt.Sprintf("error occurred and logged when CallSoPlugin in %s with %s\n%s\n", raw.Code, engFlagName, eStr))
				}
			} else {
				//println(raw.Code)
				res, errList = s.CallGoPlugin(engFlagName, raw.Body)
				if errList != nil && len(errList) > 0 {
					eStr := ""
					for _, e := range errList {
						//log.Error(e.Error())
						eStr += e.Error()
					}
					log.Error(fmt.Sprintf("error occurred and logged when CallGoPlugin in %s with %s\n%s\n", raw.Code, engFlagName, eStr))
				}
			}
			//check if processing has a result
			if len(res) == 0 {
				continue
			}
			res["code"] = raw.Code
			res["name"] = raw.Name
			res["fetchtime"] = raw.FetchTime
			res["updatetime"] = raw.UpdateTime
			res["refinetime"] = srcTime

			refinedRes = append(refinedRes, res)
		}
		err = cur.Close(context.TODO())
		if err != nil {
			log.Error(fmt.Sprintf("failed to close cur with %s\n%s", engFlagName, err.Error()))
		}
		//clear the collection in refined DB
		err = DeleteByTime(engFlagName, srcTime)
		if err != nil {
			log.Error(fmt.Sprintf("failed to purge %s\n%s", engFlagName, err.Error()))
		}

		//insert into collection
		for _, res := range refinedRes {
			err = InsertOne(engFlagName, res)
			if err != nil {
				log.Error(fmt.Sprintf("failed to insert %s in %s\n%s", res["code"], engFlagName, err.Error()))
			}
		}

	}
	log.Info("process done " + srcTime.String()[:10])
	return nil
}
