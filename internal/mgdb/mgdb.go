package mgdb

import (
	"context"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
)

var (
	err        error
	ctx        context.Context
	client     *qmgo.Client
	db         *qmgo.Database
	collection *qmgo.Collection

	cliContent *qmgo.QmgoClient
)

func Init() error {

	ctx = context.Background()
	client, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: "mongodb://127.0.0.1:27017"})
	if err != nil {
		return err
	}
	db = client.Database("vez")
	collection = db.Collection("content")

	cliContent, err = qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://127.0.0.1:27017", Database: "vez", Coll: "content"})
	if err != nil {
		return err
	}

	cliContent.CreateIndexes(ctx, []options.IndexModel{{Key: []string{"source", "id"}}})
	return err
}
