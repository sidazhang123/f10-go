package model

import (
	"context"
	"encoding/json"
	"github.com/sidazhang123/f10-go/basic/common"
	proto "github.com/sidazhang123/f10-go/srv/feed/proto/feed"
	"sort"
	"strconv"
	"time"
)

//filter docs from f10-refined according to rules
func (s Service) GenerateFocus(date string) error {
	//read rules
	err, rulesDTO := s.ReadRule([]*proto.Rule{{TarCol: ""}})
	if err != nil {
		return err
	}
	colToRule := map[string][]ReadRuleDTO{}
	for _, rDTO := range rulesDTO {
		if _, ok := colToRule[rDTO.TarCol]; ok {
			colToRule[rDTO.TarCol] = append(colToRule[rDTO.TarCol], rDTO)
		} else {
			colToRule[rDTO.TarCol] = []ReadRuleDTO{rDTO}
		}
	}
	for tarCol, rules := range colToRule {
		switch tarCol {
		case "latest_tips":
			err = LTFocus(rules, date)
			if err != nil {
				return err
			}
			break
		case "latest_tips_news":
			err = LTNews(rules)
			if err != nil {
				return err
			}
			break
		case "shareholder_analysis":
			err = SAFocus(rules, date)
			if err != nil {
				return err
			}
			break
		case "financial_analysis":
			err = FAFocus(rules, date)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

// purge collection by default
func (s Service) PurgeFocus(date string) (error, int) {
	dateKey := ""
	if len(date) > 0 {
		dateKey = "gentime"
	}
	return PurgeByDate(Params.FocusCollectionName, dateKey, date)
}

// read all by default
func (s Service) ReadFocus(date string, chanId string, del int32, fav int32) (error, string) {
	dateKey := ""
	if len(date) > 0 {
		dateKey = "gentime"
	}

	var focusItems []FocusItem
	err := ReadFocus(Params.FocusCollectionName, &focusItems, dateKey, date, chanId, del, fav)
	if err != nil {
		return err, ""
	}
	sort.SliceStable(focusItems, func(i, j int) bool {
		fetchtimex, _ := time.Parse(common.TimestampLayout[:10], focusItems[i].Fetchtime)
		fetchtimey, _ := time.Parse(common.TimestampLayout[:10], focusItems[j].Fetchtime)
		return fetchtimex.After(fetchtimey)
	})
	sort.SliceStable(focusItems, func(i, j int) bool {
		fetchtimex := focusItems[i].Fetchtime
		fetchtimey := focusItems[j].Fetchtime
		codex, _ := strconv.ParseInt(focusItems[i].Code, 10, 64)
		codey, _ := strconv.ParseInt(focusItems[j].Code, 10, 64)
		if fetchtimex == fetchtimey {
			if codex > codey {
				return false
			}
			return true
		}
		return false
	})

	bytesRes, err := json.Marshal(focusItems)
	if err != nil {
		log.Info(err.Error())
		//panic(err)
		return err, ""
	}
	return nil, string(bytesRes)
}

func (s Service) ToggleFocusDel(objectId string, v int32) error {
	return UpdateOneField(Params.FocusCollectionName, objectId, "del", v)
}

func (s Service) ToggleFocusFav(objectId string, v int32) error {
	return UpdateOneField(Params.FocusCollectionName, objectId, "fav", v)
}

func (s Service) GetFocusStat() (error, map[string]int32) {
	return FocusStatAgg(Params.FocusCollectionName)
}

func (s Service) DeleteOutdatedFocus() (error, int) {
	//get chan-day mapping
	purgeNum := 0
	err, cur := FindAll(Params.OutdatedCollectionName)
	if err != nil {
		return err, 0
	}
	for cur.Next(context.Background()) {
		var chanDayMap DaysToDel
		err := cur.Decode(&chanDayMap)
		if err != nil {
			return err, 0
		}

		err, num := DeleteByChanAndDate(Params.FocusCollectionName, chanDayMap.Channel, time.Now().AddDate(0, 0, -chanDayMap.Day))
		if err != nil {
			return err, 0
		}
		purgeNum += num
	}
	return nil, purgeNum
}
func (s Service) DeleteRecoveryBin() (error, int) {
	return DeleteManyByField(Params.FocusCollectionName, "del", []interface{}{int32(1)})
}

func (s Service) GetChanODay() (error, map[string]int32) {
	// get all chan names from rules
	err, chanNames := GetDistinctValue(Params.RulesCollectionName, "channel")
	if err != nil {
		return err, nil
	}
	// get all channel:day from outdated_chan
	err, cur := FindAll(Params.OutdatedCollectionName)
	chanDayMap := map[string]int32{}
	for cur.Next(context.Background()) {
		var chanDay DaysToDel
		err := cur.Decode(&chanDay)
		if err != nil {
			return err, nil
		}
		chanDayMap[chanDay.Channel] = int32(chanDay.Day)
	}
	// if chan only in rules - init it
	for _, c := range chanNames {
		if _, ok := chanDayMap[c]; !ok {
			chanDayMap[c] = -1
		}
	}
	// if chan only in map - del it
	for c := range chanDayMap {
		if !contains(chanNames, c) {
			delete(chanDayMap, c)
		}
	}
	return nil, chanDayMap
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (s Service) SetChanODay(daysToDel []interface{}) (error, int) {
	err, modified := PurgeCollectionAndInsertMany(Params.OutdatedCollectionName, daysToDel)
	if err != nil {
		return err, 0
	}
	return nil, modified
}
