package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type StdConnector struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	context    context.Context
}

type Connector interface {
	Database() *mongo.Database
	Collection(coll string, opts ...*options.CollectionOptions) *mongo.Collection
	NewGridfsBucket() (*gridfs.Bucket, error)
	WithContext(context.Context) Connector
	WithCollection(coll string, opts ...*options.CollectionOptions) Connector
	Find(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	FetchAll(cur *mongo.Cursor, results interface{}) error
	Decode(cur *mongo.Cursor, val interface{}) error
	Next(cur *mongo.Cursor) bool
	Count(filter interface{}, opts ...*options.CountOptions) (cnt int64, err error)
	Distinct(fieldName string, filter interface{}, opts ...*options.DistinctOptions) (res []interface{}, err error)
	FindOneAndDelete(filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	FindOneAndReplace(filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult
	FindOneAndUpdate(filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error)
	UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error)
	UpdateById(id interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error)
	ReplaceOne(filter interface{}, update interface{}, opts ...*options.ReplaceOptions) (res *mongo.UpdateResult, err error)
	InsertOne(document interface{}, opts ...*options.InsertOneOptions) (res *mongo.InsertOneResult, err error)
	InsertMany(document []interface{}, opts ...*options.InsertManyOptions) (res *mongo.InsertManyResult, err error)
	DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (res *mongo.DeleteResult, err error)
	DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (res *mongo.DeleteResult, err error)
	Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) (cur *mongo.Cursor, err error)
	Drop() error
	Watch(pipeline interface{}, opts ...*options.ChangeStreamOptions) (stream *mongo.ChangeStream, err error)
	GetNextSeq(name string, opts ...string) (res int64, err error)
}

type NewParams struct {
	Uri      string
	Database string
}

func NewConnector(params NewParams) (*StdConnector, error) {
	opts := options.Client()
	opts.SetConnectTimeout(1 * time.Second)

	bsonOpts := &options.BSONOptions{
		NilSliceAsEmpty: true,
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(params.Uri).SetBSONOptions(bsonOpts), opts)
	if err != nil {
		return nil, err
	}

	conn := StdConnector{
		client:   client,
		database: client.Database(params.Database),
		context:  context.TODO(),
	}

	return &conn, nil
}

func (conn *StdConnector) Database() *mongo.Database {
	return conn.database
}

func (conn *StdConnector) Collection(coll string, opts ...*options.CollectionOptions) *mongo.Collection {
	return conn.database.Collection(coll, opts...)
}

func (conn *StdConnector) NewGridfsBucket() (*gridfs.Bucket, error) {
	return gridfs.NewBucket(conn.database)
}

func (conn *StdConnector) WithContext(ctx context.Context) Connector {
	newConn := *conn
	newConn.context = ctx
	return &newConn
}

func (conn *StdConnector) WithCollection(coll string, opts ...*options.CollectionOptions) Connector {
	newConn := *conn
	newConn.collection = conn.database.Collection(coll, opts...)
	return &newConn
}

// read

func (conn *StdConnector) Find(filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.Find(conn.context, filter, opts...)
}

func (conn *StdConnector) FindOne(filter interface{}, opts ...*options.FindOneOptions) (res *mongo.SingleResult) {
	if conn.collection == nil {
		// enforce a SingleResult
		return mongo.NewSingleResultFromDocument(nil, errors.New("no collection set"), nil)
	}

	return conn.collection.FindOne(conn.context, filter, opts...)
}

func (conn *StdConnector) Count(filter interface{}, opts ...*options.CountOptions) (cnt int64, err error) {
	if conn.collection == nil {
		return -1, errors.New("no collection set")
	}

	return conn.collection.CountDocuments(conn.context, filter, opts...)
}

func (conn *StdConnector) Distinct(fieldName string, filter interface{}, opts ...*options.DistinctOptions) (res []interface{}, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.Distinct(conn.context, fieldName, filter, opts...)
}

// cursor

func (conn *StdConnector) Decode(cur *mongo.Cursor, val interface{}) error {
	return cur.Decode(val)
}

func (conn *StdConnector) Next(cur *mongo.Cursor) bool {
	return cur.Next(conn.context)
}

func (conn *StdConnector) FetchAll(cur *mongo.Cursor, results interface{}) (err error) {
	return cur.All(conn.context, results)
}

// read combos

func (conn *StdConnector) FindOneAndDelete(filter interface{}, opts ...*options.FindOneAndDeleteOptions) (cur *mongo.SingleResult) {
	if conn.collection == nil {
		// enforce a SingleResult
		return mongo.NewSingleResultFromDocument(nil, errors.New("no collection set"), nil)
	}

	return conn.collection.FindOneAndDelete(conn.context, filter, opts...)
}

func (conn *StdConnector) FindOneAndReplace(filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) (cur *mongo.SingleResult) {
	if conn.collection == nil {
		// enforce a SingleResult
		return mongo.NewSingleResultFromDocument(nil, errors.New("no collection set"), nil)
	}

	return conn.collection.FindOneAndReplace(conn.context, filter, replacement, opts...)
}

func (conn *StdConnector) FindOneAndUpdate(filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (cur *mongo.SingleResult) {
	if conn.collection == nil {
		// enforce a SingleResult
		return mongo.NewSingleResultFromDocument(nil, errors.New("no collection set"), nil)
	}

	return conn.collection.FindOneAndUpdate(conn.context, filter, update, opts...)
}

// update

func (conn *StdConnector) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.UpdateOne(conn.context, filter, update, opts...)
}

func (conn *StdConnector) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.UpdateMany(conn.context, filter, update, opts...)
}

func (conn *StdConnector) UpdateById(id interface{}, update interface{}, opts ...*options.UpdateOptions) (res *mongo.UpdateResult, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.UpdateByID(conn.context, id, update, opts...)
}

func (conn *StdConnector) ReplaceOne(filter interface{}, update interface{}, opts ...*options.ReplaceOptions) (res *mongo.UpdateResult, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.ReplaceOne(conn.context, filter, update, opts...)
}

// insert

func (conn *StdConnector) InsertOne(document interface{}, opts ...*options.InsertOneOptions) (res *mongo.InsertOneResult, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.InsertOne(conn.context, document, opts...)
}

func (conn *StdConnector) InsertMany(document []interface{}, opts ...*options.InsertManyOptions) (res *mongo.InsertManyResult, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.InsertMany(conn.context, document, opts...)
}

// delete

func (conn *StdConnector) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (res *mongo.DeleteResult, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.DeleteOne(conn.context, filter, opts...)
}

func (conn *StdConnector) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (res *mongo.DeleteResult, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.DeleteMany(conn.context, filter, opts...)
}

// aggregate

func (conn *StdConnector) Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) (cur *mongo.Cursor, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.Aggregate(conn.context, pipeline, opts...)
}

// various

func (conn *StdConnector) Drop() (err error) {
	if conn.collection == nil {
		return errors.New("no collection set")
	}

	return conn.collection.Drop(conn.context)
}

func (conn *StdConnector) Watch(pipeline interface{}, opts ...*options.ChangeStreamOptions) (stream *mongo.ChangeStream, err error) {
	if conn.collection == nil {
		return nil, errors.New("no collection set")
	}

	return conn.collection.Watch(conn.context, pipeline, opts...)
}

func (conn *StdConnector) GetNextSeq(name string, opts ...string) (seq int64, err error) {
	if len(name) == 0 {
		if conn.collection == nil {
			return 0, errors.New("no collection set")
		}

		name = conn.collection.Name()
	}

	seqCollection := "Sequences"
	if len(opts) > 0 {
		seqCollection = opts[0]
	}

	res := conn.WithCollection(seqCollection).FindOneAndUpdate(
		bson.D{{"_id", name}},
		bson.D{{"$inc", bson.D{{"Current", 1}}}},
		options.FindOneAndUpdate().SetUpsert(true),
		options.FindOneAndUpdate().SetReturnDocument(options.After),
		options.FindOneAndUpdate().SetProjection(bson.D{{"Current", 1}}))

	if res == nil {
		return 0, nil
	}

	var data bson.M
	if err := res.Decode(&data); err != nil {
		return 0, err
	}

	switch v := data["Current"].(type) {
	case int32:
		return int64(int(v)), nil
	case int64:
		return v, nil
	case float64:
		return int64(int(v)), nil
	default:
		return 0, errors.New("unknown return type")
	}
}
