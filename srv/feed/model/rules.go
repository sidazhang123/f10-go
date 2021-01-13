package model

import (
	"context"
	"fmt"
	proto "github.com/sidazhang123/f10-go/srv/feed/proto/feed"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func (s Service) CreateRule(rules []*proto.Rule) (error, int) {
	var rules_to_insert []interface{}
	for _, r := range rules {
		var cond1, cond2 []string
		for _, c := range r.Cond1 {
			if len(strings.TrimSpace(c)) > 0 {
				cond1 = append(cond1, c)
			}
		}
		for _, e := range r.Cond2 {
			if len(strings.TrimSpace(e)) > 0 {
				cond2 = append(cond2, e)
			}
		}
		rules_to_insert = append(rules_to_insert, CreateRuleDTO{
			TarCol:  strings.TrimSpace(r.GetTarCol()),
			Channel: strings.TrimSpace(r.GetChannel()),
			Key:     strings.TrimSpace(r.GetKey()),
			Cond1:   cond1,
			Cond2:   cond2,
		})
	}
	err, insertedN := InsertMany(Params.RulesCollectionName, rules_to_insert)
	if err != nil {
		return err, -1
	}
	return nil, insertedN
}

// @Params a Rule with tarCol
func (s Service) ReadRule(rules []*proto.Rule) (error, []ReadRuleDTO) {
	if len(rules) != 1 {
		return fmt.Errorf("RuleReq malformed %+v", rules), nil
	}
	err, cur := FindRuleByTarCol(Params.RulesCollectionName, rules[0].TarCol)
	if err != nil {
		return fmt.Errorf("failed to FindAll from %s \n%s", Params.RulesCollectionName, err.Error()), nil
	}
	var rList []ReadRuleDTO
	for cur.Next(context.TODO()) {
		var r ReadRuleDTO
		err := cur.Decode(&r)

		if err != nil {
			return fmt.Errorf("failed to Decode cur to proto.Rule\n%s", err.Error()), nil
		}
		rList = append(rList, r)
	}
	return nil, rList
}

//Params array of rules with id
func (s Service) UpdateRule(rules []*proto.Rule) (error, int) {
	var rDTOList []ReadRuleDTO

	for _, r := range rules {
		var cond1, cond2 []string
		for _, c := range r.GetCond1() {
			if len(strings.TrimSpace(c)) > 0 {
				cond1 = append(cond1, c)
			}
		}
		for _, e := range r.GetCond2() {
			if len(strings.TrimSpace(e)) > 0 {
				cond2 = append(cond2, e)
			}
		}
		id, err := primitive.ObjectIDFromHex(r.Id)
		if err != nil {
			return err, -1
		}

		rDTOList = append(rDTOList, ReadRuleDTO{
			TarCol:  strings.TrimSpace(r.GetTarCol()),
			Channel: strings.TrimSpace(r.GetChannel()),
			Key:     strings.TrimSpace(r.GetKey()),
			Cond1:   cond1,
			Cond2:   cond2,
			ID:      id,
		})
	}
	return ReplaceMany(Params.RulesCollectionName, rDTOList)
}

//Params array of rule Ids
//identified by tarCol to remove corresponding focuses as well
func (s Service) DeleteRule(rules []*proto.Rule) (err error, num int) {
	var ids []string
	var rids []interface{}
	for _, r := range rules {
		if r.Id != "" {
			ids = append(ids, r.Id)
		}
		if r.TarCol == "applyToFocus" {
			rids = append(rids, r.Id)
		}
	}
	err, num = DeleteManyById(Params.RulesCollectionName, ids)
	e, n := DeleteManyByField(Params.FocusCollectionName, "rid", rids)
	if e != nil {
		err = fmt.Errorf("delR: %s\ndelF %s\n", err.Error(), e.Error())
	}
	num += n * 10000
	return
}
