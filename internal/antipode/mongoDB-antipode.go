package antipode

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	clientOptions *options.ClientOptions
	database      string
	collection    string
}

type Document struct {
	Key   string  `bson:"key"`
	Value AntiObj `bson:"value"`
}

func CreateMongoDB(host string, port string, database string, collection string) MongoDB {
	return MongoDB{options.Client().ApplyURI("mongodb://" + host + ":" + port),
		database,
		collection,
	}
}

// devo verificar primeiro se já existe essa key? E caso exista fazer update?
// falta fechar a conecção
func (m MongoDB) write(ctx context.Context, key string, obj AntiObj) error {

	client, err := mongo.Connect(ctx, m.clientOptions)
	if err != nil {
		return err
	}

	// Disconnect from MongoDB when the function returns
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	mongoObj := Document{
		Key:   key,
		Value: obj,
	}

	_, err = client.Database(m.database).Collection(m.collection).InsertOne(ctx, mongoObj)

	return err
}

// posso assumir que não há mais do que um objeto com a mesma key?
func (m MongoDB) read(ctx context.Context, key string) (AntiObj, error) {

	client, err := mongo.Connect(ctx, m.clientOptions)
	if err != nil {
		return AntiObj{}, err
	}

	// Disconnect from MongoDB when the function returns
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	filter := bson.D{{"key", key}}

	var result Document
	err = client.Database(m.database).Collection(m.collection).FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return AntiObj{}, err
	}

	return result.Value, nil
}

func (m MongoDB) barrier(ctx context.Context, lineage []WriteIdentifier, datastoreID string) error {

	client, err := mongo.Connect(ctx, m.clientOptions)
	if err != nil {
		return err
	}

	// Disconnect from MongoDB when the function returns
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	for _, writeIdentifier := range lineage {
		if writeIdentifier.Dtstid == datastoreID {
			for {
				filter := bson.D{{"key", writeIdentifier.Key}}

				fmt.Println("key: ", writeIdentifier.Key)

				var result Document
				err = client.Database(m.database).Collection(m.collection).FindOne(context.Background(), filter).Decode(&result)

				if !errors.Is(err, mongo.ErrNoDocuments) && err != nil {
					return err
				} else if errors.Is(err, mongo.ErrNoDocuments) { //the version replication process is not yet completed
					fmt.Println("replication in progress")
					continue
				} else {
					if result.Value.Version == writeIdentifier.Version { //the version replication process is already completed
						fmt.Println("replication done: ", result.Value.Version)
						break
					} else { //the version replication process is not yet completed
						fmt.Println("replication of the new version in progress")
						continue
					}
				}
			}
		}
	}

	return nil
}
