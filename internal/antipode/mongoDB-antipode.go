package antipode

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	clientOptions *options.ClientOptions
	database      string
	collection    string
}

type MongoDBObject struct {
	key   string
	value AntiObj
}

func CreateMongoDB(host string, port string, database string, collection string) MongoDB {
	return MongoDB{options.Client().ApplyURI("mongodb://" + host + ":" + port),
		database,
		collection,
	}
}

// devo verificar primeiro se já existe essa key?
func (m MongoDB) write(ctx context.Context, key string, obj AntiObj) error {

	fmt.Println("write")

	client, err := mongo.Connect(ctx, m.clientOptions)
	if err != nil {
		fmt.Println(err)
		return err
	}

	mongoObj := MongoDBObject{key, obj}

	fmt.Println("middle")

	_, err = client.Database(m.database).Collection(m.collection).InsertOne(ctx, mongoObj)

	fmt.Println(err)

	return err
}

// posso assumir que não há mais do que um objeto com a mesma key?
func (m MongoDB) read(ctx context.Context, key string) (AntiObj, error) {

	client, err := mongo.Connect(ctx, m.clientOptions)
	if err != nil {
		return AntiObj{}, err
	}

	filter := bson.D{{"key", key}}

	cursor, err := client.Database(m.database).Collection(m.collection).Find(ctx, filter)
	if err != nil {
		return AntiObj{}, err
	}

	var results []MongoDBObject
	if err = cursor.All(ctx, &results); err != nil {
		return AntiObj{}, err
	}

	if len(results) != 1 {
		return AntiObj{}, fmt.Errorf("The key %s does not exist.", key)
	}

	return results[0].value, err
}

func (m MongoDB) barrier(ctx context.Context, lineage []WriteIdentifier, datastoreID string) error {

	client, err := mongo.Connect(ctx, m.clientOptions)
	if err != nil {
		return err
	}

	for _, writeIdentifier := range lineage {
		if writeIdentifier.Dtstid == datastoreID {
			for {
				filter := bson.D{{"key", writeIdentifier.Key}}

				cursor, err := client.Database(m.database).Collection(m.collection).Find(ctx, filter)
				if err != nil {
					return err
				}

				var results []MongoDBObject
				if err = cursor.All(ctx, &results); err != nil {
					return err
				}

				if len(results) == 0 {
					continue
				} else {
					if results[0].value.Version == writeIdentifier.Version { //the version replication process is already completed
						break
					} else { //the version replication process is not yet completed
						continue
					}
				}
			}
		}
	}

	return nil
}
