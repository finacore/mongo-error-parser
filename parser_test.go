// The meparser package provide a way to parse the Mongo response errors. Their
// usage is really simple and consiste just in invocate the function Parse passing
// as parameter the error received by the mongodb operation. For more information
// see the Parse documentation bellow.
package meparser

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://guest:guest@localhost:42069"

func create_unique_index(c *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	name, err := c.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil || name == "" {
		log.Fatalf("Fail to create index: %v", err)
	}
}

func get_client() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic("Unable to get mongo client")
	}

	return client
}

func validatorSchema(cli *mongo.Database, collection string) {

	var jsonSchema = bson.M{
		"bsonType":             "object",
		"required":             []string{"name", "lastname"},
		"additionalProperties": false,
		"properties": bson.M{
			"lastname": bson.M{
				"bsonType":    "string",
				"description": "must be a string and is required",
			},
			"name": bson.M{
				"bsonType":    "string",
				"maxLength":   4,
				"description": "must be a string and is required",
			},
		},
	}

	var validator = bson.M{
		"$jsonSchema": jsonSchema,
	}

	val := options.CreateCollection().SetValidator(validator)

	if err := cli.CreateCollection(context.TODO(), collection, val); err != nil {
		panic("unable to create collection")
	}
}

func prepare() *mongo.Collection {
	client := get_client()
	cli := client.Database("meparser")

	collName := uuid.NewString()
	validatorSchema(cli, collName)

	c := cli.Collection(collName)
	create_unique_index(c)

	return c
}

func generate_index_error(c *mongo.Collection) error {
	type md struct {
		Name     string
		Lastname string
	}

	mdi := md{
		Name:     "john",
		Lastname: "Smith",
	}

	//force insertion of 2 register to create the error
	_, errs := c.InsertOne(context.TODO(), mdi)
	if errs == nil {
		_, errs = c.InsertOne(context.TODO(), mdi)
	}

	return errs
}

func generate_scheme_error(c *mongo.Collection) error {
	type md struct {
		Name     string
		Lastname string
	}

	mdi := md{
		Name:     "babbets",
		Lastname: "Smith",
	}

	//force insertion of 2 register to create the error
	_, errs := c.InsertOne(context.TODO(), mdi)

	return errs
}

func Test_toMongoException(t *testing.T) {
	c := prepare()
	err := generate_index_error(c)

	ex := toMongoException(err)

	assert.NotNil(t, ex)
}

func Test_toErrorInterface(t *testing.T) {
	c := prepare()
	err := generate_scheme_error(c)

	ex := toMongoException(err)
	ei := toErrorInterface(ex)

	assert.NotNil(t, ei)
}
