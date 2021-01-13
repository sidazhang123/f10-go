package model

import (
	"context"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/plugins/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

func InsertJPushID(collection string, id string) error {
	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	_, err := col.InsertOne(context.TODO(), bson.M{"reg_id": id})
	if err != nil && !strings.Contains(err.Error(), "duplicate key") {
		return err
	}
	return nil
}

func GetLastNByField(collection, field, code string, n int64) (error, []map[string]interface{}) {
	findOptions := options.Find()
	findOptions.SetSort(bson.M{field: -1})
	findOptions.SetLimit(n)
	proj := bson.M{"_id": 0, "code": 0}
	for _, i := range []string{"refinetime", "updatetime", "fetchtime"} {
		if field != i {
			proj[i] = 0
		}
	}
	findOptions.SetProjection(proj)
	col := db.GetDB().Database(Params.RefinedDbName).Collection(collection)
	cur, err := col.Find(context.TODO(), bson.M{"code": code}, findOptions)
	if err != nil {
		return err, nil
	}
	var res []map[string]interface{}
	for cur.Next(context.Background()) {
		var r map[string]interface{}
		err := cur.Decode(&r)
		if err != nil {
			return err, nil
		}
		res = append(res, r)
	}
	return nil, res
}

func InsertMany(collection string, toInsert []interface{}) (error, int) {
	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	insertResult, err := col.InsertMany(context.TODO(), toInsert)
	if err != nil {
		return err, -1
	}
	return nil, len(insertResult.InsertedIDs)
}

func GetDistinctValue(collection, field string) (err error, res []string) {
	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	var iL []interface{}
	iL, err = col.Distinct(context.TODO(), field, bson.M{})
	if err != nil {
		return
	}
	for _, i := range iL {
		res = append(res, i.(string))
	}
	return
}

func FindAll(collection string) (error, *mongo.Cursor) {

	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	cur, err := col.Find(context.TODO(), bson.M{})
	if err != nil {
		return err, nil
	}

	return nil, cur
}

func FindRuleByTarCol(collection, tarCol string) (error, *mongo.Cursor) {

	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	var filter bson.M
	if len(tarCol) == 0 {
		filter = bson.M{}
	} else {
		filter = bson.M{"tarCol": tarCol}
	}
	cur, err := col.Find(context.TODO(), filter)
	if err != nil {
		return err, nil
	}

	return nil, cur
}

func ReplaceMany(collection string, rules []ReadRuleDTO) (error, int) {

	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	count := 0
	for _, r := range rules {
		res, err := col.ReplaceOne(context.TODO(), bson.M{"_id": r.ID}, r)
		if err != nil {
			return err, -1
		}
		count += int(res.ModifiedCount)
	}

	return nil, count
}

func DeleteManyById(collection string, ids []string) (error, int) {

	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	count := 0
	for _, i := range ids {
		id, err := primitive.ObjectIDFromHex(i)
		if err != nil {
			return err, -1
		}
		res, err := col.DeleteMany(context.TODO(), bson.M{"_id": id})
		if err != nil {
			return err, -1
		}
		count += int(res.DeletedCount)
	}

	return nil, count
}
func DeleteManyByField(collection string, field string, values []interface{}) (error, int) {

	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	count := 0
	for _, v := range values {
		res, err := col.UpdateMany(context.TODO(), bson.M{field: v}, bson.M{"$set": bson.M{"del": 2}})
		if err != nil {
			return err, -1
		}
		count += int(res.ModifiedCount)
	}

	return nil, count
}

func FindByFetchTime(collection string, fetchdate string) (error, *mongo.Cursor) {
	col := db.GetDB().Database(Params.RefinedDbName).Collection(collection)
	fetchtime, err := time.Parse(common.TimestampLayout[:10], fetchdate)
	if err != nil {
		return err, nil
	}
	cur, err := col.Find(context.TODO(), bson.M{"fetchtime": bson.M{
		"$gte": fetchtime,
		"$lt":  fetchtime.AddDate(0, 0, 1),
	}})
	if err != nil {
		return err, nil
	}

	return nil, cur
}

func FindLatestFetchTime(collection string) (error, string) {
	col := db.GetDB().Database(Params.RefinedDbName).Collection(collection)
	opts := options.Find()
	opts.SetSort(bson.D{{"fetchtime", -1}})
	opts.SetProjection(bson.D{{"fetchtime", 1}})
	opts.SetLimit(1)
	cur, err := col.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return err, ""
	}
	type One struct {
		Fetchtime time.Time
		ID        primitive.ObjectID `json:"_id" bson:"_id"`
	}
	var refinedOne One
	for cur.Next(context.Background()) {
		err := cur.Decode(&refinedOne)
		if err != nil {
			return err, ""
		}
		break
	}
	return nil, refinedOne.Fetchtime.Format(common.TimestampLayout[:10])
}

func InsertOneFocus(collection string, res map[string]interface{}) error {
	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	_, err := col.InsertOne(context.Background(), res)
	return err
}
func UpdateOneField(collection string, objectId string, field string, v interface{}) error {
	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	objID, err := primitive.ObjectIDFromHex(objectId)
	if err != nil {
		return err
	}
	_, err = col.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": bson.M{field: v}})
	return err
}
func PurgeCollectionAndInsertMany(collection string, daysToDel []interface{}) (error, int) {
	var err error
	var modified int
	client := db.GetDB()
	col := client.Database(Params.FeedDbName).Collection(collection)
	var session mongo.Session
	if session, err = client.StartSession(); err != nil {
		return err, 0
	}
	if err = session.StartTransaction(); err != nil {
		return err, 0
	}
	if err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		_, err = col.DeleteMany(context.Background(), bson.M{})
		if err != nil {
			return err
		}
		insertResult, err := col.InsertMany(context.TODO(), daysToDel)
		if err != nil {
			return err
		}
		modified = len(insertResult.InsertedIDs)
		if err = session.CommitTransaction(sc); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err, 0
	}
	session.EndSession(context.Background())

	return nil, modified
}
func PurgeByDate(collection string, dateKey string, date string) (error, int) {
	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	var res *mongo.DeleteResult
	var err error
	if len(dateKey) > 0 {
		_, err := time.Parse(common.TimestampLayout[:10], date)
		if err != nil {
			return err, -1
		}
		res, err = col.DeleteMany(context.TODO(), bson.M{dateKey: date})
		if err != nil {
			return err, -1
		}
	} else {
		res, err = col.DeleteMany(context.TODO(), bson.M{})
		if err != nil {
			return err, -1
		}
	}
	return nil, int(res.DeletedCount)
}

func ReadFocus(collection string, focusItems *[]FocusItem, dateKey string, date string, chanId string, del int32, fav int32) error {
	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	filter := bson.M{}
	if del == 0 || del == 1 {
		filter["del"] = del
	}
	if fav == 0 || fav == 1 {
		filter["fav"] = fav
	}
	if chanId != "" {
		filter["chan"] = chanId
	}
	if len(dateKey) > 0 {
		_, err := time.Parse(common.TimestampLayout[:10], date)
		if err != nil {
			return err
		}
		filter[dateKey] = date
	}
	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	if err = cur.All(context.Background(), focusItems); err != nil {
		return err
	}
	return nil
}

func DeleteByChanAndDate(collection string, channel string, date time.Time) (error, int) {
	eStr := ""
	dCount := 0
	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	cur, err := col.Find(context.TODO(), bson.M{"chan": channel})
	if err != nil {
		return err, 0
	}
	for cur.Next(context.TODO()) {
		var f map[string]interface{}
		err := cur.Decode(&f)
		if err != nil {
			eStr += err.Error()
			continue
		}
		fetchtime, err := time.Parse(common.TimestampLayout[:10], f["fetchtime"].(string))
		if err != nil {
			eStr += err.Error()
			continue
		}

		if fetchtime.Before(date) {
			dr, err := col.UpdateOne(context.TODO(), bson.M{"_id": f["_id"]}, bson.M{"$set": bson.M{"del": 2}})
			if err != nil {
				eStr += err.Error()
				continue
			}
			dCount += int(dr.ModifiedCount)
		}
	}
	var e error
	if len(eStr) > 0 {
		e = fmt.Errorf(eStr)
	} else {
		e = nil
	}
	return e, dCount
}

func FocusStatAgg(collection string) (error, map[string]int32) {
	col := db.GetDB().Database(Params.FeedDbName).Collection(collection)
	if n, _ := col.CountDocuments(context.TODO(), bson.M{}); n == 0 {
		return nil, map[string]int32{}
	}
	notDel := bson.D{{"$match", bson.D{{"del", 0}}}}
	groupChan := bson.D{{"$group", bson.D{
		{"_id", bson.D{{"$toLower", "$chan"}}},
		{"count", bson.D{{"$sum", 1}}}}}}
	groupCount := bson.D{{"$group", bson.D{
		{"_id", nil},
		{"counts", bson.D{{"$push",
			bson.D{{"k", "$_id"}, {"v", "$count"}}}}}}}}
	replaceRoot := bson.D{{"$replaceRoot", bson.D{{
		"newRoot", bson.D{{"$arrayToObject", "$counts"}}}}}}
	cur, err := col.Aggregate(context.TODO(), mongo.Pipeline{notDel, groupChan, groupCount, replaceRoot})
	if err != nil {
		return err, nil
	}

	var stat []map[string]int32
	if err = cur.All(context.TODO(), &stat); err != nil {
		return err, nil
	}
	return nil, stat[0]
}
