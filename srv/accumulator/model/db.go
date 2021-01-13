package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sidazhang123/f10-go/basic/common"
	"github.com/sidazhang123/f10-go/plugins/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var today = time.Now().Format(common.TimestampLayout[:10])

//
func ReadMany(dbName, collection, code string, start, end time.Time) (error, *mongo.Cursor) {
	col := db.GetDB().Database(dbName).Collection(collection)
	filter := bson.M{}
	if len(code) == 6 {
		filter["code"] = code
	}
	if start.IsZero() != end.IsZero() {
		return fmt.Errorf("[ReadMany] start & end time must be provided to match"), nil
	}
	if !start.IsZero() && !end.IsZero() && end.After(start) {
		filter["fetchtime"] = bson.M{"$gte": start, "$lt": end}
	}

	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		return err, nil
	}
	return nil, cur
}

func FindOne(dbName, collection, field, value string) (error, map[string]interface{}) {
	col := db.GetDB().Database(dbName).Collection(collection)
	r := col.FindOne(context.TODO(), bson.M{field: value})
	if r.Err() != nil {
		return r.Err(), nil
	}
	var res map[string]interface{}
	err := r.Decode(&res)
	if err != nil {
		return err, nil
	}
	return nil, res
}

func InsertOne(dbName, collection string, v map[string]interface{}) error {
	col := db.GetDB().Database(dbName).Collection(collection)
	_, err := col.InsertOne(context.TODO(), v)
	return err
}

func UpdateOne(dbName, collection, code string, v bson.M) error {
	col := db.GetDB().Database(dbName).Collection(collection)
	_, err := col.UpdateOne(context.TODO(), bson.M{"code": code}, bson.M{"$set": v})
	return err
}

func AddFieldAndUpdate(collection, code, name, k, v string) error {
	col := db.GetDB().Database(Params.AccumulatedDbName).Collection(collection)
	r := col.FindOne(context.TODO(), bson.M{"code": code})
	if r.Err() != nil {
		_, err := col.InsertOne(context.TODO(), bson.M{"code": code, "name": name, k: v, "updatetime": today})
		if err != nil {
			return err
		}
	} else {
		var record map[string]interface{}
		err := r.Decode(&record)
		if err != nil {
			return err
		}
		if _, ok := record[k]; !ok {
			_, err := col.UpdateOne(context.TODO(), bson.M{"code": code}, bson.M{"$set": bson.M{k: v, "updatetime": today, "name": name}})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
	if code not in the new location - f10-acc.{win_#_name}, insertOne with code,name,[field],updateTime
	else: add to the field array if not exist, and update the document's name and updateTime
*/
func AppendAndUpdateByField(collection, code, name, field, value string) error {

	col := db.GetDB().Database(Params.AccumulatedDbName).Collection(collection)
	opts := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(context.TODO(), bson.M{"code": code}, bson.M{"$addToSet": bson.M{field: value}, "$set": bson.M{"updatetime": today, "name": name}}, opts)
	return err
}

// only FA uses this
//                 collection     code      field    value(objId not string)
var existingRec = map[string]map[string]map[string]interface{}{}

func AddToJsonByFieldAndUpdate(collection, code, name, field, value string) error {

	if code == "" && name == "" && field == "" && value == "" {
		// a signal to store the map
		col := db.GetDB().Database(Params.AccumulatedDbName).Collection(collection)
		for code, v := range existingRec[collection] {
			_, err := col.UpdateOne(context.TODO(), bson.M{"code": code}, bson.M{"$set": v}, options.Update().SetUpsert(true))
			if err != nil {
				return err
			}
		}
		// clear map before the next run
		existingRec = map[string]map[string]map[string]interface{}{}
		return nil
	}
	if len(existingRec) == 0 {
		existingRec[collection] = map[string]map[string]interface{}{}
		col := db.GetDB().Database(Params.AccumulatedDbName).Collection(collection)
		cur, err := col.Find(context.TODO(), bson.M{})
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
			existingRec[collection][code] = rec
		}
	}

	if _, ok := existingRec[collection][code]; !ok {
		// new code
		existingRec[collection][code] = map[string]interface{}{"code": code, "name": name, field: value, "updatetime": today}
	} else {
		// new field
		if _, ok := existingRec[collection][code][field]; !ok {
			existingRec[collection][code][field] = value
			existingRec[collection][code]["updatetime"] = today
			existingRec[collection][code]["name"] = name
		} else {
			// new date in an existing field
			var newDateMsg, oldDateMsg map[string]string
			err := json.Unmarshal([]byte(value), &newDateMsg)
			if err != nil {
				return err
			}
			err = json.Unmarshal([]byte(existingRec[collection][code][field].(string)), &oldDateMsg)
			if err != nil {
				return err
			}
			changed := false
			for dateK, msgV := range newDateMsg {
				if _, ok := oldDateMsg[dateK]; !ok {
					oldDateMsg[dateK] = msgV
					changed = true
				}
			}
			if changed {
				b, err := json.Marshal(oldDateMsg)
				if err != nil {
					return err
				}
				existingRec[collection][code][field] = string(b)
				existingRec[collection][code]["updatetime"] = today
				existingRec[collection][code]["name"] = name
			}

		}

	}
	return nil
}

func FindEndDate(collection string, latest bool) (error, time.Time) {
	col := db.GetDB().Database(Params.RefinedDbName).Collection(collection)
	opts := options.Find()
	v := 1
	if latest {
		v = -1
	}
	opts.SetSort(bson.D{{"fetchtime", v}})
	opts.SetProjection(bson.D{{"fetchtime", 1}})
	opts.SetLimit(1)
	cur, err := col.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return err, time.Time{}
	}
	type One struct {
		Fetchtime time.Time
		ID        primitive.ObjectID `json:"_id" bson:"_id"`
	}
	var refinedOne One
	for cur.Next(context.Background()) {
		err := cur.Decode(&refinedOne)
		if err != nil {
			return err, time.Time{}
		}
		break
	}
	return nil, refinedOne.Fetchtime
}
