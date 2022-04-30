package mgdb

import (
	// "errors"
	"fmt"
	"strings"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/operator"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/midoks/vez/internal/tools"
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

type ContentBid struct {
	Url        string    `bson:"url"`
	Source     string    `bson:"source"`
	User       string    `bson:"user"`
	Id         string    `bson:"id"`
	Title      string    `bson:"title"`
	Html       string    `bson:"html"`
	Length     int       `bson:"length"`
	Updatetime time.Time `bson:"updatetime" json:"updatetime"`
	Createtime time.Time `bson:"createtime" json:"createtime"`

	MgID string `bson:"_id"`
}

func ContentAdd(data Content) (result *qmgo.InsertOneResult, err error) {

	data.Length = len(data.Html)

	one := Content{}
	err = cliContent.Find(ctx, M{"source": data.Source, "id": data.Id}).One(&one)

	if err != nil {
		return ContentOriginAdd(data)
	}

	oneData := M{"$set": M{
		"title":      one.Title,
		"html":       one.Html,
		"updatetime": time.Now(),
	}}

	err = cliContent.UpdateOne(ctx, M{"source": data.Source, "id": data.Id}, oneData)
	if err != nil {
		return nil, fmt.Errorf("content update error: %v", err)
	}
	return nil, nil
}

func ContentOriginAdd(data Content) (result *qmgo.InsertOneResult, err error) {

	data.Length = len(data.Html)
	data.Updatetime = time.Now()
	data.Createtime = time.Now()

	result, err = collection.InsertOne(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("add error: %T", err)
	}
	return result, nil
}

func ContentOriginFindOne(source, id string) (result ContentBid, err error) {
	one := ContentBid{}
	err = cliContent.Find(ctx, M{"source": source, "id": id}).One(&one)
	return one, err
}

func ContentOriginFindNewsestOne(source string) (result ContentBid, err error) {
	one := ContentBid{}
	err = cliContent.Find(ctx, M{"source": source}).Sort("-_id").One(&one)
	return one, err
}

func ContentOriginFind(limit ...int64) (result []ContentBid, err error) {
	var batch []ContentBid

	var bNum int64
	if len(limit) > 0 {
		bNum = limit[0]
	} else {
		bNum = 15
	}
	err = cliContent.Find(ctx, D{}).Sort("-_id").Limit(bNum).All(&batch)
	return batch, err
}

func ContentOriginFindSoso(id, sort, op, keyword string, limit ...int64) (result []ContentBid, err error) {
	var batch []ContentBid

	var bNum int64
	if len(limit) > 0 {
		bNum = limit[0]
	} else {
		bNum = 15
	}

	sortField := fmt.Sprintf("%s_id", sort)

	opt := M{}
	if !strings.EqualFold(id, "") {
		mgId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return batch, err
		}
		opt["_id"] = M{op: mgId}
	}

	if !strings.EqualFold(keyword, "") {
		opt["title"] = M{operator.Regex: keyword, "$options": "im"}
	}

	err = cliContent.Find(ctx, opt).Sort(sortField).Limit(bNum).All(&batch)
	return batch, err
}

func ContentOriginFindId(id, sort string, limit ...int64) (result []ContentBid, err error) {
	var batch []ContentBid

	var bNum int64
	if len(limit) > 0 {
		bNum = limit[0]
	} else {
		bNum = 15
	}

	sortField := fmt.Sprintf("%s_id", sort)
	if strings.EqualFold(id, "") {
		err = cliContent.Find(ctx, D{}).Sort(sortField).Limit(bNum).All(&batch)

		fmt.Println(sortField, batch, err)
		return batch, err
	}

	mgId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return batch, err
	}

	opt := M{"_id": M{operator.Lt: mgId}}
	err = cliContent.Find(ctx, opt).Sort(sortField).Limit(bNum).All(&batch)
	return batch, err
}

func ContentOriginFindIdGt(id, sort string, limit ...int64) (result []ContentBid, err error) {
	var batch []ContentBid

	var bNum int64
	if len(limit) > 0 {
		bNum = limit[0]
	} else {
		bNum = 15
	}

	sortField := fmt.Sprintf("%s_id", sort)
	if strings.EqualFold(id, "") {
		err = cliContent.Find(ctx, D{}).Sort(sortField).Limit(bNum).All(&batch)
		return batch, err
	}

	_idObj, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return batch, err
	}

	opt := M{"_id": M{operator.Gt: _idObj}}
	err = cliContent.Find(ctx, opt).Sort(sortField).Limit(bNum).All(&batch)
	return batch, err
}

func ContentNewsest() ([]ContentBid, error) {
	var batch []ContentBid

	err = cliContent.Find(ctx, D{}).Sort("-createtime").Limit(5).All(&batch)
	return batch, err
}

func ContentRand() (result ContentBid, err error) {
	one := ContentBid{}

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

func ContentRandSource(source string) (result ContentBid, err error) {
	one := ContentBid{}

	randStage := D{
		{
			operator.Sample,
			D{
				{"size", 1},
			},
		},
	}

	matchStage := D{{"$match", []E{{"source", D{{operator.Eq, source}}}}}}

	err = cliContent.Aggregate(ctx, qmgo.Pipeline{matchStage, randStage}).One(&one)
	return one, err
}

func ContentOneByOne(source string) (result ContentBid, err error) {
	one := ContentBid{}
	opt := M{"source": M{operator.Eq: source}}
	filePath := "/tmp/vez.txt"

	if tools.IsExist(filePath) {
		c, err := tools.ReadFile(filePath)
		if err != nil {
			return one, err
		}

		if strings.EqualFold(strings.TrimSpace(c), "") {

			goto ContentOneByOneGoto
		}
		_idObj, err := primitive.ObjectIDFromHex(c)
		if err != nil {
			return one, err
		}

		opt["_id"] = M{operator.Gt: _idObj}
	}

ContentOneByOneGoto:

	one, err = ContentFindSourceLimit(source, 3)
	if err != nil {
		return one, fmt.Errorf("mongodb find error: %v", err)
	}

	err = tools.WriteFile(filePath, string(one.MgID))
	if err != nil {
		return one, fmt.Errorf("write file error: %v", err)
	}
	return one, nil
}

func ContentFindSourceLimit(source string, limit ...int64) (result ContentBid, err error) {
	one := ContentBid{}
	opt := M{"source": M{operator.Eq: source}}

	var bLimit int64
	if len(limit) > 0 {
		bLimit = limit[0]
	} else {
		bLimit = 3
	}

	var i int64
	for i = 0; i < bLimit; i++ {
		err = cliContent.Find(ctx, opt).Skip(i).Sort("+_id").Limit(1).One(&one)
		if !strings.EqualFold(strings.TrimSpace(one.MgID), "") {
			return one, err
		}
	}
	return one, err
}

func ContentRandNum(num int64) ([]ContentBid, error) {
	var batch []ContentBid

	randStage := D{
		{
			operator.Sample,
			D{
				{
					"size",
					num,
				},
			},
		},
	}

	err = cliContent.Aggregate(ctx, qmgo.Pipeline{randStage}).All(&batch)
	return batch, err
}
