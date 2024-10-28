package mongodb_test

import (
	"errors"
	"github.com/mbretter/go-mongodb"
	"github.com/mbretter/go-mongodb/types"
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

// Example demonstrates the process of creating, reading, updating, and deleting a user in a MongoDB database.
func Example() {
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

// Example_Test tests the scenario for user creation, ensuring correct insertion and data consistency.
// The Connector is mocked using the auto-generated mock by mockery.
// The ObjectId generator function is stubbed, to get reproducable results.
func Example_TestCreate(t *testing.T) {
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
