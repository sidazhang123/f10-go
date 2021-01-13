package model

import (
	"context"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/plugins/db"
	"time"
)

type StockBody struct {
	Code       string
	Name       string
	Flag       string
	FlagName   string
	Body       string
	FetchTime  time.Time
	UpdateTime time.Time
	Uid        string
}

//func InsertMany(stocks []interface{}) (err error) {
//	_, err = db.GetDB().Database(Params.DbDatabase).Collection(Params.DbCollectionPrefix+common.FlagNameToCollSuffix[]).InsertMany(context.TODO(), stocks)
//	return
//}
func InsertOne(i interface{}, flagName string) error {
	_, err := db.GetDB().Database(opts.DbName).Collection(common.FlagNameToCollSuffix[flagName]).InsertOne(context.TODO(), i)
	if err != nil {
		return err
	}
	//log2.Info(fmt.Sprintf("Inserted a single document: %v", insertResult.InsertedID))
	return nil
}
