package mgdb

import (
	"errors"
	// "fmt"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/operator"
)

type Content struct {
	Url        string    `bson:"url"`
	Source     string    `bson:"source"`
	User       string    `bson:"user"`
	Id         string    `bson:"id"`
	Title      string    `bson:"title"`
	Html       string    `bson:"html"`
	Length     int       `bson:"length"`
	Updatetime time.Time `bson:"updatetime" json:"updatetime"`
	Createtime time.Time `bson:"createtime" json:"createtime"`
}

func ContentAdd(data Content) (result *qmgo.InsertOneResult, err error) {
	if collection != nil {
		data.Length = len(data.Html)

		one := Content{}
		err = cliContent.Find(ctx, M{"source": data.Source, "id": data.Id}).One(&one)

		if err == nil {
			data.Updatetime = time.Now()
			err = cliContent.UpdateOne(ctx, M{"source": data.Source, "id": data.Id}, data)
			return nil, err
		}

		return ContentOriginAdd(data)
	}

	return nil, errors.New("mongo disconnected!")
}

func ContentOriginAdd(data Content) (result *qmgo.InsertOneResult, err error) {
	if collection != nil {
		data.Length = len(data.Html)
		data.Updatetime = time.Now()
		data.Createtime = time.Now()

		result, err = collection.InsertOne(ctx, data)
		if err != nil {
			return nil, err
		}
		return result, err
	}

	return nil, errors.New("mongo disconnected!")
}

func ContentOriginFindOne(source, id string) (result *Content, err error) {
	one := &Content{}
	err = cliContent.Find(ctx, M{"source": source, "id": id}).One(one)
	return one, err
}

func ContentOriginFind() (result []Content, err error) {
	var batch []Content

	// matchStage := bson.D{{"$match", []bson.D{{"name", bson.D{{"$gt", 30}}}}}}
	// groupStage := bson.D{{"$group", bson.D{{"_id", nil}, {"total", bson.D{{"$sum", "$age"}}}}}}

	// err = cliContent.Aggregate(ctx, qmgo.Pipeline{matchStage, groupStage}).All(&batch)
	// fmt.Println(err, batch)

	err = cliContent.Find(ctx, D{}).Sort("-createtime").Limit(15).All(&batch)

	// fmt.Println(err, batch)
	return batch, err
}

func ContentRand() (result *Content, err error) {
	one := &Content{}

	randStage := D{
		{
			operator.Sample,
			D{
				{
					"size",
					1,
				},
			},
		},
	}

	err = cliContent.Aggregate(ctx, qmgo.Pipeline{randStage}).One(&one)
	return one, err
}
