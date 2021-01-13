package model

import (
	"context"
	"github.com/sidazhang123/f10-go/plugins/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Raw struct {
	Code       string
	Name       string
	Flag       string
	FlagName   string
	Body       string
	FetchTime  time.Time
	UpdateTime time.Time
	Uid        string
}

func DeleteByTime(collection string, updateTime time.Time) error {
	year, month, day := updateTime.Date()
	utc := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	col := db.GetDB().Database(Params.DbName).Collection(collection)
	_, err := col.DeleteMany(context.TODO(), bson.M{
		"refinetime": bson.M{
			"$gte": utc,
			"$lt":  utc.AddDate(0, 0, 1),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
func InsertOne(collection string, toInsert map[string]interface{}) error {
	col := db.GetDB().Database(Params.DbName).Collection(collection)

	_, err := col.InsertOne(context.TODO(), toInsert)
	if err != nil {
		return err
	}
	//log2.Info(fmt.Sprintf("Inserted a single document: %v", insertResult.InsertedID))
	return nil
}

// pass in Asia/Shanghai time
func FindByDate(collection string, date time.Time) (error, *mongo.Cursor) {
	year, month, day := date.Date()
	utc := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	col := db.GetDB().Database(Params.RawDbName).Collection(collection)
	cur, err := col.Find(context.TODO(), bson.M{
		"fetchtime": bson.M{
			"$gte": utc,
			"$lt":  utc.AddDate(0, 0, 1),
		},
	})
	if err != nil {
		return err, nil
	}

	return nil, cur
}

func FindByDateAndCode(collection string, date time.Time, code string) (error, []*Raw) {
	year, month, day := date.Date()
	utc := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	col := db.GetDB().Database(Params.RawDbName).Collection(collection)
	cur, err := col.Find(context.TODO(), bson.M{
		"fetchtime": bson.M{
			"$gte": utc,
			"$lt":  utc.AddDate(0, 0, 1),
		},
		"code": code,
	})
	if err != nil {
		return err, nil
	}
	var res []*Raw
	for cur.Next(context.TODO()) {
		var raw Raw
		err := cur.Decode(&raw)
		if err != nil {
			return err, nil
		}
		res = append(res, &raw)
	}
	return nil, res
}
