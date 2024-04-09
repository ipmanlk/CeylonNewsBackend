package nosqldb

import (
	"context"
	"os"

	"ipmanlk/cnapi/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const databaseName = "cn"
const collectionName = "news"

type Book struct {
	Title  string
	Author string
}

func InitDB() error {
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		uri = "mongodb://localhost:27000"
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	return nil
}

func InsertItem(item common.NewsItem) error {
	coll := client.Database(databaseName).Collection(collectionName)
	_, err := coll.InsertOne(context.TODO(), item)
	return err
}

func InsertItems(items []common.NewsItem) error {
	coll := client.Database(databaseName).Collection(collectionName)
	ordered := false
	documents := make([]interface{}, len(items))
	for i, item := range items {
		documents[i] = item
	}
	_, err := coll.InsertMany(context.TODO(), documents, &options.InsertManyOptions{Ordered: &ordered})
	return err
}

func GetItemByID(id string) (common.NewsItem, error) {
	coll := client.Database(databaseName).Collection(collectionName)
	var item common.NewsItem
	err := coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&item)
	return item, err
}

func SearchItems(langs []common.Lang, sources []string, query string, cursor string, pageSize int) (*common.PaginationResponse, error) {
	var items []common.NewsItem

	decodedCursor, err := common.DecodeCursorNoSql(cursor)
	if err != nil {
		return nil, err
	}

	filter := bson.M{}

	if query != "" {
		regex := bson.M{"$regex": query, "$options": "i"} 
		filter["$or"] = []bson.M{
			{"title": regex},
			{"content_text": regex},
		}
	}

	if len(langs) > 0 {
		filter["language"] = bson.M{"$in": langs}
	}

	if len(sources) > 0 {
		filter["source_name"] = bson.M{"$in": sources}
	}

	coll := client.Database(databaseName).Collection(collectionName)

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "_id", Value: -1}})
	opts.SetLimit(int64(pageSize))

	if decodedCursor.Direction == common.PaginationDirectionPrev {
		opts.SetHint(bson.D{{Key: "_id", Value: -1}})
		filter["_id"] = bson.M{"$lt": decodedCursor.ItemID}
	} else {
		filter["_id"] = bson.M{"$gt": decodedCursor.ItemID}
	}

	csr, err := coll.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer csr.Close(context.Background())

	for csr.Next(context.Background()) {
		var item common.NewsItem
		err := csr.Decode(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	var prevCursor, nextCursor string

	if len(items) > 0 {
		prevCursor = common.CreateCursorNoSql(items[len(items)-1].ID, common.PaginationDirectionPrev)
		nextCursor = common.CreateCursorNoSql(items[0].ID, common.PaginationDirectionNext)
	}

	response := common.PaginationResponse{
		Data:   items,
		Paging: common.PaginationPaging{Prev: prevCursor, Next: nextCursor},
	}

	return &response, nil
}
