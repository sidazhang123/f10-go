package model

import (
	"fmt"
	"reflect"
	"strings"
)

type Opts struct {
	RefinedDbName     string `json:"refined_db_name"`
	AccumulatedDbName string `json:"accumulated_db_name"`
	//LTFields          string `json:"LT_fields"`
	//FAFields          string `json:"FA_fields"`
	//SAFields          string `json:"SA_fields"`
	Win1Name     string `json:"win_1_name"`
	Win1Seq      string `json:"win_1_seq"`
	Win1Capacity int    `json:"win_1_capacity"`
	Win2Name     string `json:"win_2_name"`
	Win2Seq      string `json:"win_2_seq"`
	Win2Capacity int    `json:"win_2_capacity"`
	Win3Name     string `json:"win_3_name"`
	Win3Seq      string `json:"win_3_seq"`
	Win3Capacity int    `json:"win_3_capacity"`
	Win4Name     string `json:"win_4_name"`
	Win4Seq      string `json:"win_4_seq"`
	Win4Capacity int    `json:"win_4_capacity"`
	Win5Name     string `json:"win_5_name"`
	Win5Seq      string `json:"win_5_seq"`
	Win5Capacity int    `json:"win_5_capacity"`
}

func mark(err interface{}, line int, appendix ...string) error {
	if reflect.TypeOf(err).String() == "string" {
		return fmt.Errorf("[%s mark@%d] %s\n", strings.Join(appendix, ","), line, err)
	} else {
		_, ok := reflect.TypeOf(err).MethodByName("Error")
		if ok {
			return fmt.Errorf("[%s mark@%d] %s\n", strings.Join(appendix, ","), line, err.(error).Error())
		}
	}
	return fmt.Errorf("[MARK ERR] wrong type passing in")
}
