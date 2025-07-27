package mgdb

import (
	"context"
	"fmt"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/midoks/vez/internal/conf"
)

var (
	err        error
	ctx        context.Context
	client     *qmgo.Client
	db         *qmgo.Database
	collection *qmgo.Collection

	cliContent *qmgo.QmgoClient
)

type (
	// M is an alias of bson.M
	M = bson.M
	// A is an alias of bson.A
	A = bson.A
	// D is an alias of bson.D
	D = bson.D
	// E is an alias of bson.E
	E = bson.E
)

func Init() error {
	link := "mongodb://" + conf.Mongodb.Addr

	ctx = context.Background()

	// 配置连接池参数
	connectTimeoutMS := int64(10000) // 10秒连接超时
	maxPoolSize := uint64(100)       // 最大连接池大小
	minPoolSize := uint64(10)        // 最小连接池大小

	client, err = qmgo.NewClient(ctx, &qmgo.Config{
		Uri:              link,
		ConnectTimeoutMS: &connectTimeoutMS,
		MaxPoolSize:      &maxPoolSize,
		MinPoolSize:      &minPoolSize,
	})
	if err != nil {
		return fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	db = client.Database(conf.Mongodb.Db)
	collection = db.Collection("content")

	cliContent, err = qmgo.Open(ctx, &qmgo.Config{
		Uri:              link,
		Database:         conf.Mongodb.Db,
		Coll:             "content",
		ConnectTimeoutMS: &connectTimeoutMS,
		MaxPoolSize:      &maxPoolSize,
		MinPoolSize:      &minPoolSize,
	})
	if err != nil {
		return fmt.Errorf("failed to open MongoDB collection: %v", err)
	}

	// 创建复合索引
	err = cliContent.CreateIndexes(ctx, []options.IndexModel{
		{Key: []string{"source", "id"}},
		{Key: []string{"-createtime"}},
		{Key: []string{"title"}},
	})
	if err != nil {
		return fmt.Errorf("failed to create indexes: %v", err)
	}

	return nil
}
