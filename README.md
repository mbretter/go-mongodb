[![](https://github.com/mbretter/go-mongodb/actions/workflows/test.yml/badge.svg)](https://github.com/mbretter/go-mongodb/actions/workflows/test.yml)
[![](https://goreportcard.com/badge/mbretter/go-mongodb)](https://goreportcard.com/report/mbretter/go-mongodb "Go Report Card")
[![codecov](https://codecov.io/gh/mbretter/go-mongodb/graph/badge.svg?token=YMBMKY7W9X)](https://codecov.io/gh/mbretter/go-mongodb)
[![GoDoc](https://godoc.org/github.com/mbretter/go-mongodb?status.svg)](https://pkg.go.dev/github.com/mbretter/go-mongodb)

This package wraps the go mongo-driver by providing a so-called "Connector", this makes the mongodb connection testable/mockable.
The original driver is not really testable, it is hard/impossible to mock the package.

Usually in go the interfaces are defined at the consumer side, but in this case an interface is provided to keep the codebase small. 

The provided connector interface can easily be mocked using [mockery](https://github.com/vektra/mockery).

Additionally, this package provides some datatypes, like UUID, ObjectId, NullString, nullable numbers and a datatype for 
storing binary data.

This is version 2 of the package, which uses mongo-driver v2.

## Install

```
go get github.com/mbretter/go-mongodb/v2
```

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

err := connector.FindOne(bson.D{{"_id", id}}).Decode(&ret)
```

The functions are exactly the same as those of the mongo-driver, e.g. Find, FindOne, Count, UpdateOne, ...

### Sequences

Besided the wrapped functions of the mongo-driver, a function for fetching sequence numbers was implemented, it returns 
the next number by using `FindOneAndUpdate()`, with the upsert option.

```go
nextNumber, err := connector.GetNextSeq("Users")
```
The sequence numbers are stored into a "Sequences" collection, the _id is the provided name ("Users" in this case) and the 
current number is stored into the "Current" field. If no name was provided, the name of the current collection is used.
You can optionally provide the name of the collection where the sequences are stored.

## Datatypes

Besides the BSON conversion, all datatypes are supporting JSON encoding/decoding, by implementing the marshal/unmarshal functions.

### ObjectId

The types package adds an `ObjectId` replacement of the mongo-drivers `primitive.ObjectId`. 
The original ObjectId has two disadvantages:
* an empty ObjectId is stored as "000000000000000000000000", instead of null, which is kind of weird.
* every conversion between a string and an ObjectId has to be done using `ObjectIDFromHex()`, which adds a lot of extra code.

The ObjectId provided by the types package derives from string, so conversions can be easily done using a simple type cast, 
but you have to make sure, that the string contains a valid ObjectId in HEX format, otherwise you will get an error, when marshalling it to BSON.

```go
hexOid := "61791c74138d41367e52d832"

objectId := types.ObjectId(hexOid)
```

### UUID

The UUID derives from string for easy conversion, it's BSON represenation is `primitive.Binary` with the subtype of `bson.TypeBinaryUUID`.
This means it is stored as native UUID into the database. An empty UUID is treated as null when converting to BSON.

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

## Example

This is a real life example, demonstrating how the connector could be integrated.
It defines the database access layer (UserDb) and connects this to a model (UserModel) by using an interface (UserDbInterface), 
defined at the consumer side. This makes the codebase easily unit testable.

It also defines the provider functions, as you would use it with [wire](https://github.com/google/wire/).

### Definitions

```go
package user

import (
    "errors"
    "github.com/mbretter/go-mongodb/v2"
    "github.com/mbretter/go-mongodb/v2/types"
    "github.com/stretchr/testify/assert"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "log"
    "os"
    "testing"
)

// User represents a user entity with an ID, username, and personal details such as firstname and lastname.
type User struct {
    Id        types.ObjectId `bson:"_id"`
    Username  string         `bson:"username,omitempty"`
    Firstname string         `bson:"firstname,omitempty"`
    Lastname  string         `bson:"lastname,omitempty"`
}

// UserDbInterface defines the methods required to interact with the User entity in the database.
// Insert adds a new User to the database and returns an error if the operation fails.
// Update modifies an existing User in the database and returns an UpdateResult and an error if the operation fails.
// Read retrieves a User by their ObjectId and returns the User and an error if the operation fails.
// Delete removes a User by their ObjectId and returns a DeleteResult and an error if the operation fails.
type UserDbInterface interface {
    Insert(*User) error
    Update(user *User) (*mongo.UpdateResult, error)
    Read(id types.ObjectId) (*User, error)
    Delete(id types.ObjectId) (*mongo.DeleteResult, error)
}

// UserModel provides methods for managing user data by integrating with the UserDbInterface.
type UserModel struct {
    db UserDbInterface
}

// ProvideModel initializes a User instance with the provided UserDbInterface.
func ProvideModel(db UserDbInterface) UserModel {
    return UserModel{db: db}
}

func (m UserModel) Create(user User) (User, error) {
    user.Id = types.NewObjectId()
    
    err := m.db.Insert(&user)
    if err != nil {
        return user, err
    }
    
    return user, nil
}

func (m UserModel) ReadById(id types.ObjectId) (*User, error) {
    return m.db.Read(id)
}

func (m UserModel) Update(user User) error {
    _, err := m.db.Update(&user)
    return err
}

func (m UserModel) DeleteById(id types.ObjectId) error {
    _, err := m.db.Delete(id)
    
    return err
}

// UserDb provides methods to interact with the user collection in the database using a Connector interface.
type UserDb struct {
    conn mongodb.Connector
}

func (d *UserDb) Insert(user *User) error {
    _, err := d.conn.InsertOne(user)
    if err != nil {
        return err
    }
    
    return nil
}

func (d *UserDb) Update(user *User) (*mongo.UpdateResult, error) {
    res, err := d.conn.UpdateById(user.Id, bson.D{{"$set", user}})
    
    return res, err
}

func (d *UserDb) Read(id types.ObjectId) (*User, error) {
    var ret User
    
    err := d.conn.FindOne(bson.D{{"_id", id}}).Decode(&ret)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, nil
        }
        return nil, err
    }
    
    return &ret, nil
}

func (d *UserDb) Delete(id types.ObjectId) (*mongo.DeleteResult, error) {
    return d.conn.DeleteOne(bson.D{{"_id", id}})
}

func ProviderUserDb(conn mongodb.Connector) *UserDb {
    return &UserDb{
        conn: conn.WithCollection("user"),
    }
}
```

### Main

This is a code fragment wich shows how to put all the components together.

- build the connector and connect to the database
- build the database layer (UserDb)
- build the model (UserModel)
- do some CRUD operations

```go
func main() {
    conn, err := mongodb.NewConnector(mongodb.NewParams{
        Uri:      os.Getenv("MONGODB_URI"),
        Database: os.Getenv("MONGODB_DB"),
    })
    
    if err != nil {
        log.Fatalf("failed to connect to db: %v\n", err)
    }
    
    userDb := ProviderUserDb(conn)
    userModel := ProvideModel(userDb)
    user := User{
        Username:  "foo@bar.com",
        Firstname: "John",
        Lastname:  "Doe",
    }
    
    user, err = userModel.Create(user)
    if err != nil {
        log.Fatalf("failed to create user: %v", err)
    }
    
    log.Printf("created user: %v", user)
    
    existingUser, err := userModel.ReadById(user.Id)
    if err != nil {
        log.Fatalf("failed to read user: %v", err)
    }
    
    if existingUser == nil {
        log.Fatalf("user not found")
    }
    
    updateUser := User{
        Id:        existingUser.Id,
        Firstname: "Jane",
    }
    
    err = userModel.Update(updateUser)
    if err != nil {
        log.Fatalf("failed to update user: %v", err)
    }
    
    err = userModel.DeleteById(user.Id)
    if err != nil {
        log.Fatalf("failed to delete user: %v", err)
    }
}
```

### Tests

This integration test uses a mock of the Connector, the mock can easily be auto-generated by mockery.

If you would like to unit test the UserModel, you have to mock the UserDbInterface.

```go
func TestCreate(t *testing.T) {
    newUserId := "66cc9ca8c042f7a732b7fc2a"
    types.SetObjectIdGenerator(func() string { return newUserId })
    
    user := User{
        Id:        types.ObjectId(newUserId),
        Username:  "foo@bar.com",
        Firstname: "John",
        Lastname:  "Doe",
    }
    
    tests := []struct {
        name string
        err  error
    }{
        {
            "Success",
            nil,
        },
        {
            "Failed",
            errors.New("some database error occurred"),
        },
    }
    
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            conn := NewConnectorMock(t)
            conn.EXPECT().WithCollection("user").Return(conn)
    
            userDb := ProviderUserDb(conn)
            userModel := ProvideModel(userDb)
    
            res := mongo.InsertOneResult{}
            conn.EXPECT().InsertOne(&user).Return(&res, test.err)
    
            user, err := userModel.Create(user)
            if test.err == nil {
                assert.Nil(t, err)
                assert.Equal(t, user.Id, types.ObjectId(newUserId))
                assert.Equal(t, user.Username, "foo@bar.com")
                assert.Equal(t, user.Firstname, "John")
                assert.Equal(t, user.Lastname, "Doe")
            } else {
                assert.NotNil(t, err)
            }
        })
    }
}

```