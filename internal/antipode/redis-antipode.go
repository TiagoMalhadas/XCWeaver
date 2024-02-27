package antipode

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/TiagoMalhadas/xcweaver"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func createRedis(redis_host string, redis_port string, redis_password string) Redis {
	return Redis{redis.NewClient(&redis.Options{
		Addr:     redis_host + ":" + redis_port,
		Password: redis_password, // no password set
		DB:       0,              // use default DB
	})}
}

func (r Redis) write(ctx context.Context, key string, obj xcweaver.AntiObj) error {

	jsonAntiObj, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, key, jsonAntiObj, 0).Err()

	return err
}

func (r Redis) read(ctx context.Context, key string) (xcweaver.AntiObj, error) {

	jsonAntiObj, err := r.client.Get(ctx, key).Bytes()

	if err != nil {
		return xcweaver.AntiObj{}, err
	}

	var obj xcweaver.AntiObj
	err = json.Unmarshal(jsonAntiObj, &obj)
	if err != nil {
		return xcweaver.AntiObj{}, err
	}

	return obj, err
}

func (r Redis) barrier(ctx context.Context, lineage []xcweaver.WriteIdentifier, datastoreID string) error {

	for _, writeIdentifier := range lineage {
		if writeIdentifier.Dtstid == datastoreID {
			for {
				jsonAntiObj, err := r.client.Get(ctx, writeIdentifier.Key).Bytes()

				if !errors.Is(err, redis.Nil) && err != nil {
					return err
				} else if errors.Is(err, redis.Nil) { //the version replication process is not yet completed
					continue
				} else {
					var obj xcweaver.AntiObj
					err = json.Unmarshal(jsonAntiObj, &obj)
					if err != nil {
						return err
					} else if obj.Version == writeIdentifier.Version { //the version replication process is already completed
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
