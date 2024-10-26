[![](https://github.com/mbretter/go-mongodb/actions/workflows/test.yml/badge.svg)](https://github.com/mbretter/go-mongodb/actions/workflows/test.yml)
[![](https://goreportcard.com/badge/mbretter/go-mongodb)](https://goreportcard.com/report/mbretter/go-mongodb "Go Report Card")
[![codecov](https://codecov.io/gh/mbretter/go-mongodb/graph/badge.svg?token=YMBMKY7W9X)](https://codecov.io/gh/mbretter/go-mongodb)

This package wraps the go mongodb driver by providing a so-called "Connector", this makes the mongodb connection testable/mockable.
The original driver is not really testable, it is hard/impossible to mock the package.

Usually in go the interfaces are defined at the consumer, but in this case an interface is provided to keep the codebase small. 

The provided connector interface can easily be mocked using [mockery](https://github.com/vektra/mockery).

Additionaly this package provides some datatypes, like UUID, ObjectId, NullString, nullable numbers and a datatype for 
storing binary data.

## Connector

Constructing a new connector, by returning a StdConnector:

```go
connector, err := mongodb.NewConnector(mongodb.NewParams{
    Uri:      "mongodb://user:pass@127.0.0.1/mydb",
    Database: "mydb",
})
```

Setting the target collection, all operations are executed against this collection, a copy of the connector is returned 
leaving the original connector unmodified:
```go
connector = connector.WithCollection("Users")
```

Setting the context, by default context.TODO() is used, a new connector will be returned using the supplied context 
for consecutive calls:
```go
connector = connector.WithContext(context.Background())
```

Finding a row:
```go
var ret User

err := d.conn.FindOne(bson.D{{"_id", id}}).Decode(&ret)
```

The functions are exactly the same as those of the mongodb driver, e.g. Find, FindOne, Count, UpdateOne, ...

### Sequences

Besided the wrapped functions of the mongodb driver, a function for fetching sequence numbers was implemented, it returns 
the next number by using `FindOneAndUpdate()`, with the upsert option.

```go
nextNumber, err := connector.GetNextSeq("Users")
```
The sequence numbers are stored into a "Sequences" collection, the _id is the provided name ("Users" in this case) and the 
current number is stored into the "Current" field. If no name was provided, the name of the current collection is used.
You can optionally provide the name of the collection where the sequences are stored.

## Datatypes

All datatypes are supporting JSON encoding/decoding, by implementing the marshal/unmarshal functions.

### ObjectId

The types package adds an `ObjectId` replacement of the mongodb drivers `primitive.ObjectId`. 
The original ObjectId has two disadvantages:
* an empty ObjectId is stored as "000000000000000000000000", instead of null, which is kind of weird.
* every conversion between a string and an ObjectId has to be done using `ObjectIDFromHex()`, which adds a lot of extra code.

The ObjectId provided by the types package derives from string, so conversions can be easily done using a simple type cast, 
but you have to make sure, that the string is a valid ObjectId, otherwise you will get an error, when marshalling it to BSON.

```go
hexOid := "61791c74138d41367e52d832"

objectId := types.ObjectId(hexOid)
```

### UUID

The UUID derives from string for easily converting from strings, is it mashaled as `primitive.Binary` with the subtype of `bson.TypeBinaryUUID`.
This means it is store as native UUID into the database. An empty UUID is treated as null when converting to BSON.

Under the hood, it uses github.com/google/uuid for parsing/generating values. Invalid values will produce an error if 
converting to BSON.

```go
uuid := types.NewUuid()

uuidStr := "9f53f39d-62b6-43ac-a267-f25848739aeb"
uuid = types.UUID(uuidStr)
```

### Binary

The binary datatype stores any arbitrary value as binary, the binary subtype is `bson.TypeBinaryGeneric`. The JSON 
representation of the binary is base64.
It is very useful if you do not want to/can use GridFS, but keep in mind that the maximum BSON document size is 16MBytes. 

### NullString

The NullString datatype BSON-encodes empty strings to null and vice versa.

### NullNumbers

The various number datatypes are treated as BSON-null if their value is 0 oder 0.0 and vice versa.
