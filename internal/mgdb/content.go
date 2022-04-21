package mgdb

import (
	"errors"
	"fmt"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/operator"
	"go.mongodb.org/mongo-driver/bson"
)

type Content struct {
	Url    string `bson:"url"`
	Source string `bson:"source"`
	User   string `bson:"user"`
	Id     string `bson:"id"`
	Title  string `bson:"title"`
	Html   string `bson:"html"`
	Length int    `bson:"length"`
}

func ContentAdd(data Content) (result *qmgo.InsertOneResult, err error) {
	if collection != nil {
		data.Length = len(data.Html)

		one := Content{}
		err = cliContent.Find(ctx, bson.M{"source": data.Source, "id": data.Id}).One(&one)

		if err == nil {
			err = cliContent.UpdateOne(ctx, bson.M{"source": data.Source, "id": data.Id}, data)
			return nil, err
		}

		return ContentOriginAdd(data)
	}

	return nil, errors.New("mongo disconnected!")
}

func ContentOriginAdd(data Content) (result *qmgo.InsertOneResult, err error) {
	if collection != nil {
		data.Length = len(data.Html)

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
	err = cliContent.Find(ctx, bson.M{"source": source, "id": id}).One(one)
	return one, err
}

func ContentRand() (result *Content, err error) {
	one := &Content{}

	randStage := bson.D{
		{
			operator.Sample,
			bson.D{
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

func Debug() {

	fmt.Println("ddd")
}
