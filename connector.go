// Package mongodb wraps the go mongodb driver by providing a so-called "Connector", this makes the mongodb connection testable/mockable.
// The original driver is not really testable, it is hard/impossible to mock the package.
// Usually in go the interfaces are defined at the consumer side, but in this case an interface is provided to keep the codebase small.
//
// The provided connector interface can easily be mocked using mockery.
//
// Additionally, this package provides some datatypes, like UUID, ObjectId, NullString, nullable numbers, and a datatype for
// storing binary data.
package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

// StdConnector handles connections and interactions with the MongoDB client, database, and collections.
type StdConnector struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	context    context.Context
}

// Connector provides methods for database and collection operations.
type Connector interface {
	Database() *mongo.Database
	Collection(coll string, opts ...options.Lister[options.CollectionOptions]) *mongo.Collection
	NewGridfsBucket() (*mongo.GridFSBucket, error)
	WithContext(context.Context) Connector
	WithCollection(coll string, opts ...options.Lister[options.CollectionOptions]) Connector
	Find(filter interface{}, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
	FindOne(filter interface{}, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	FetchAll(cur *mongo.Cursor, results interface{}) error
	Decode(cur *mongo.Cursor, val interface{}) error
	Next(cur *mongo.Cursor) bool
	Count(filter interface{}, opts ...options.Lister[options.CountOptions]) (cnt int64, err error)
	Distinct(fieldName string, filter interface{}, opts ...options.Lister[options.DistinctOptions]) (res *mongo.DistinctResult, err error)
	FindOneAndDelete(filter interface{}, opts ...options.Lister[options.FindOneAndDeleteOptions]) *mongo.SingleResult
	FindOneAndReplace(filter interface{}, replacement interface{}, opts ...options.Lister[options.FindOneAndReplaceOptions]) *mongo.SingleResult
	FindOneAndUpdate(filter interface{}, update interface{}, opts ...options.Lister[options.FindOneAndUpdateOptions]) *mongo.SingleResult
	UpdateOne(filter interface{}, update interface{}, opts ...options.Lister[options.UpdateOneOptions]) (res *mongo.UpdateResult, err error)
	UpdateMany(filter interface{}, update interface{}, opts ...options.Lister[options.UpdateManyOptions]) (res *mongo.UpdateResult, err error)
	UpdateById(id interface{}, update interface{}, opts ...options.Lister[options.UpdateOneOptions]) (res *mongo.UpdateResult, err error)
	ReplaceOne(filter interface{}, update interface{}, opts ...options.Lister[options.ReplaceOptions]) (res *mongo.UpdateResult, err error)
	InsertOne(document interface{}, opts ...options.Lister[options.InsertOneOptions]) (res *mongo.InsertOneResult, err error)
	InsertMany(document []interface{}, opts ...options.Lister[options.InsertManyOptions]) (res *mongo.InsertManyResult, err error)
	DeleteOne(filter interface{}, opts ...options.Lister[options.DeleteOneOptions]) (res *mongo.DeleteResult, err error)
	DeleteMany(filter interface{}, opts ...options.Lister[options.DeleteManyOptions]) (res *mongo.DeleteResult, err error)
	Aggregate(pipeline interface{}, opts ...options.Lister[options.AggregateOptions]) (cur *mongo.Cursor, err error)
	Indexes() (*mongo.IndexView, error)
	CreateIndex(model mongo.IndexModel, opts ...options.Lister[options.CreateIndexesOptions]) (string, error)
	SearchIndexes() (*mongo.SearchIndexView, error)
	CreateSearchIndex(model mongo.SearchIndexModel, opts ...options.Lister[options.CreateSearchIndexesOptions]) (string, error)
	Drop() error
	Watch(pipeline interface{}, opts ...options.Lister[options.ChangeStreamOptions]) (stream *mongo.ChangeStream, err error)
	GetNextSeq(name string, opts ...string) (res int64, err error)
}

var ErrNoCollectionSet = errors.New("no collection set")

// NewParams holds the parameters required to establish a new connection to a database.
type NewParams struct {
	Uri      string
	Database string
}

// NewConnector establishes a new connection to the mongo database using the provided parameters.
// It returns a StdConnector
func NewConnector(params NewParams) (*StdConnector, error) {
	opts := options.Client().ApplyURI(params.Uri)
	opts.SetConnectTimeout(5 * time.Second)

	bsonOpts := &options.BSONOptions{
		NilSliceAsEmpty: true,
	}

	opts.SetBSONOptions(bsonOpts)

	client, err := mongo.Connect(opts)
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

// Database returns the current mongo.Database instance associated with the StdConnector.
func (conn *StdConnector) Database() *mongo.Database {
	return conn.database
}

// Collection returns a mongo.Collection object for the specified collection name with additional options if provided.
func (conn *StdConnector) Collection(coll string, opts ...options.Lister[options.CollectionOptions]) *mongo.Collection {
	return conn.database.Collection(coll, opts...)
}

// NewGridfsBucket creates a new GridFS bucket for the current database.
func (conn *StdConnector) NewGridfsBucket() (*mongo.GridFSBucket, error) {
	return conn.database.GridFSBucket(), nil
}

// WithContext returns a copy of the StdConnector with the specified context.
func (conn *StdConnector) WithContext(ctx context.Context) Connector {
	newConn := *conn
	newConn.context = ctx
	return &newConn
}

// WithCollection returns a copy of StdConnector with the specified collection and optional collection options.
func (conn *StdConnector) WithCollection(coll string, opts ...options.Lister[options.CollectionOptions]) Connector {
	newConn := *conn
	newConn.collection = conn.database.Collection(coll, opts...)
	return &newConn
}

// read

// Find executes a find query in the collection with the given filter and options.
// Returns a cursor to the found documents or an error if the collection is not set or if the query fails.
func (conn *StdConnector) Find(filter interface{}, opts ...options.Lister[options.FindOptions]) (cur *mongo.Cursor, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.Find(conn.context, filter, opts...)
}

// FindOne retrieves a single document from the collection based on the provided filter and options.
// Returns a SingleResult containing the document or an error if the collection is not set.
func (conn *StdConnector) FindOne(filter interface{}, opts ...options.Lister[options.FindOneOptions]) (res *mongo.SingleResult) {
	if conn.collection == nil {
		// enforce a SingleResult
		return mongo.NewSingleResultFromDocument(nil, ErrNoCollectionSet, nil)
	}

	return conn.collection.FindOne(conn.context, filter, opts...)
}

// Count returns the count of documents matching the given filter and options or an error if the collection is not set.
func (conn *StdConnector) Count(filter interface{}, opts ...options.Lister[options.CountOptions]) (cnt int64, err error) {
	if conn.collection == nil {
		return -1, ErrNoCollectionSet
	}

	return conn.collection.CountDocuments(conn.context, filter, opts...)
}

// Distinct retrieves distinct values for a specified field in the collection that matches the given filter and options.
// Returns a slice of distinct values or an error if the collection is not set or the operation fails.
func (conn *StdConnector) Distinct(fieldName string, filter interface{}, opts ...options.Lister[options.DistinctOptions]) (res *mongo.DistinctResult, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}
	return conn.collection.Distinct(conn.context, fieldName, filter, opts...), nil
}

// cursor

// Decode decodes the current document pointed to by the cursor into the provided value. Returns an error if decoding fails.
func (conn *StdConnector) Decode(cur *mongo.Cursor, val interface{}) error {
	return cur.Decode(val)
}

// Next progresses the given MongoDB cursor to the next document and returns true if a next document is available.
func (conn *StdConnector) Next(cur *mongo.Cursor) bool {
	return cur.Next(conn.context)
}

// FetchAll retrieves all the documents from the provided MongoDB cursor and stores them in the results interface.
func (conn *StdConnector) FetchAll(cur *mongo.Cursor, results interface{}) (err error) {
	return cur.All(conn.context, results)
}

// read combos

// FindOneAndDelete deletes a single document from the collection that matches the provided filter and returns the deleted document.
func (conn *StdConnector) FindOneAndDelete(filter interface{}, opts ...options.Lister[options.FindOneAndDeleteOptions]) (cur *mongo.SingleResult) {
	if conn.collection == nil {
		// enforce a SingleResult
		return mongo.NewSingleResultFromDocument(nil, ErrNoCollectionSet, nil)
	}

	return conn.collection.FindOneAndDelete(conn.context, filter, opts...)
}

// FindOneAndReplace replaces a single document in the collection matching the given filter with the provided replacement.
func (conn *StdConnector) FindOneAndReplace(filter interface{}, replacement interface{}, opts ...options.Lister[options.FindOneAndReplaceOptions]) (cur *mongo.SingleResult) {
	if conn.collection == nil {
		// enforce a SingleResult
		return mongo.NewSingleResultFromDocument(nil, ErrNoCollectionSet, nil)
	}

	return conn.collection.FindOneAndReplace(conn.context, filter, replacement, opts...)
}

// FindOneAndUpdate updates a single document matching the filter and returns the updated document.
func (conn *StdConnector) FindOneAndUpdate(filter interface{}, update interface{}, opts ...options.Lister[options.FindOneAndUpdateOptions]) (cur *mongo.SingleResult) {
	if conn.collection == nil {
		// enforce a SingleResult
		return mongo.NewSingleResultFromDocument(nil, ErrNoCollectionSet, nil)
	}

	return conn.collection.FindOneAndUpdate(conn.context, filter, update, opts...)
}

// update

// UpdateOne executes an update operation on a single document in the collection.
// Parameters are a filter to select the document, an update to apply, and optional update options.
// Returns the result of the update operation or an error if the operation fails.
func (conn *StdConnector) UpdateOne(filter interface{}, update interface{}, opts ...options.Lister[options.UpdateOneOptions]) (res *mongo.UpdateResult, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.UpdateOne(conn.context, filter, update, opts...)
}

// UpdateMany updates multiple documents in the collection based on the provided filter and update parameters.
// It returns a mongo.UpdateResult containing details about the operation or an error if one occurred.
func (conn *StdConnector) UpdateMany(filter interface{}, update interface{}, opts ...options.Lister[options.UpdateManyOptions]) (res *mongo.UpdateResult, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.UpdateMany(conn.context, filter, update, opts...)
}

// UpdateById updates a single document in the collection based on its ID.
func (conn *StdConnector) UpdateById(id interface{}, update interface{}, opts ...options.Lister[options.UpdateOneOptions]) (res *mongo.UpdateResult, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.UpdateByID(conn.context, id, update, opts...)
}

// ReplaceOne replaces a single document in the collection that matches the specified filter with the provided update.
func (conn *StdConnector) ReplaceOne(filter interface{}, update interface{}, opts ...options.Lister[options.ReplaceOptions]) (res *mongo.UpdateResult, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.ReplaceOne(conn.context, filter, update, opts...)
}

// insert

// InsertOne inserts a single document into the collection.
// It returns the result of the insertion and any error encountered.
// The method takes a document to be inserted and optional insertion options.
func (conn *StdConnector) InsertOne(document interface{}, opts ...options.Lister[options.InsertOneOptions]) (res *mongo.InsertOneResult, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.InsertOne(conn.context, document, opts...)
}

// InsertMany inserts multiple documents into the collection.
// It returns an InsertManyResult containing the IDs of the inserted documents or an error if the insertion fails.
func (conn *StdConnector) InsertMany(document []interface{}, opts ...options.Lister[options.InsertManyOptions]) (res *mongo.InsertManyResult, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.InsertMany(conn.context, document, opts...)
}

// delete

// DeleteOne deletes a single document from the collection that matches the provided filter.
func (conn *StdConnector) DeleteOne(filter interface{}, opts ...options.Lister[options.DeleteOneOptions]) (res *mongo.DeleteResult, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.DeleteOne(conn.context, filter, opts...)
}

// DeleteMany deletes multiple documents from the collection that match the provided filter.
// Parameters:
// - filter: A document describing the documents to be deleted.
// - opts: Optional. Additional options to apply to the delete operation.
// Returns:
// - res: A DeleteResult containing the outcome of the delete operation.
// - err: An error if the operation fails.
func (conn *StdConnector) DeleteMany(filter interface{}, opts ...options.Lister[options.DeleteManyOptions]) (res *mongo.DeleteResult, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.DeleteMany(conn.context, filter, opts...)
}

// aggregate

// Aggregate executes an aggregation framework pipeline on the collection.
// The 'pipeline' parameter specifies an array of aggregation stages.
// The 'opts' parameters specify optional settings for the aggregate operation.
// It returns a cursor that iterates over the resulting documents.
func (conn *StdConnector) Aggregate(pipeline interface{}, opts ...options.Lister[options.AggregateOptions]) (cur *mongo.Cursor, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.Aggregate(conn.context, pipeline, opts...)
}

// SearchIndexes retrieves the search index view associated with the current collection in the MongoDB connection.
// Returns an error if no collection is set in the connector.
func (conn *StdConnector) SearchIndexes() (*mongo.SearchIndexView, error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}
	vw := conn.collection.SearchIndexes()

	return &vw, nil
}

// CreateSearchIndex creates a search index in the current collection using the specified model and options.
// Returns the name of the created index or an error if the operation fails.
func (conn *StdConnector) CreateSearchIndex(model mongo.SearchIndexModel, opts ...options.Lister[options.CreateSearchIndexesOptions]) (string, error) {
	if conn.collection == nil {
		return "", ErrNoCollectionSet
	}
	return conn.collection.SearchIndexes().CreateOne(conn.context, model, opts...)
}

// Indexes retrieves the index view object for the associated collection in the connector. Returns an error if no collection is set.
func (conn *StdConnector) Indexes() (*mongo.IndexView, error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}
	vw := conn.collection.Indexes()

	return &vw, nil
}

// CreateIndex creates a new index on the specified collection using the provided IndexModel and optional settings.
func (conn *StdConnector) CreateIndex(model mongo.IndexModel, opts ...options.Lister[options.CreateIndexesOptions]) (string, error) {
	if conn.collection == nil {
		return "", ErrNoCollectionSet
	}
	return conn.collection.Indexes().CreateOne(conn.context, model, opts...)
}

// various

// Drop removes the current collection from the database and returns an error if unsuccessful.
func (conn *StdConnector) Drop() (err error) {
	if conn.collection == nil {
		return ErrNoCollectionSet
	}

	return conn.collection.Drop(conn.context)
}

// Watch starts a change stream against the collection of the StdConnector, based on the given pipeline and options.
// It returns a pointer to a mongo.ChangeStream for iterating the changes, or an error if the collection is not set.
func (conn *StdConnector) Watch(pipeline interface{}, opts ...options.Lister[options.ChangeStreamOptions]) (stream *mongo.ChangeStream, err error) {
	if conn.collection == nil {
		return nil, ErrNoCollectionSet
	}

	return conn.collection.Watch(conn.context, pipeline, opts...)
}

// GetNextSeq increments and retrieves the next sequence number for a given name within the specified collection.
func (conn *StdConnector) GetNextSeq(name string, opts ...string) (seq int64, err error) {
	if len(name) == 0 {
		if conn.collection == nil {
			return 0, ErrNoCollectionSet
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
