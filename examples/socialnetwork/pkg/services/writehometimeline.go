package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	sn_metrics "socialnetwork/pkg/metrics"
	"socialnetwork/pkg/model"
	"socialnetwork/pkg/storage"
	sn_trace "socialnetwork/pkg/trace"

	"github.com/TiagoMalhadas/xcweaver"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type WriteHomeTimelineService interface {
	WriteHomeTimeline(ctx context.Context, msg model.Message) error
}

type writeHomeTimelineServiceOptions struct {
	RabbitMQAddr string `toml:"rabbitmq_address"`
	MongoDBAddr  string `toml:"mongodb_address"`
	RedisAddr    string `toml:"redis_address"`
	RabbitMQPort int    `toml:"rabbitmq_port"`
	MongoDBPort  int    `toml:"mongodb_port"`
	RedisPort    int    `toml:"redis_port"`
	NumWorkers   int    `toml:"num_workers"`
	Region       string `toml:"region"`
}

type writeHomeTimelineService struct {
	xcweaver.Implements[WriteHomeTimelineService]
	xcweaver.WithConfig[writeHomeTimelineServiceOptions]
	socialGraphService      xcweaver.Ref[SocialGraphService]
	redisClient             *redis.Client
	rabbitClientWriteHomeTL xcweaver.Antipode
	mongoClientWriteHomeTL  xcweaver.Antipode
}

func (w *writeHomeTimelineService) Init(ctx context.Context) error {
	logger := w.Logger(ctx)

	w.redisClient = storage.RedisClient(w.Config().RedisAddr, w.Config().RedisPort)

	var wg sync.WaitGroup
	wg.Add(w.Config().NumWorkers)
	for i := 1; i <= w.Config().NumWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			err := w.workerThread(ctx, i)
			logger.Error("error in worker thread", "msg", err.Error())
		}(i)
	}

	logger.Info("write home timeline service running!", "region", w.Config().Region, "n_workers", w.Config().NumWorkers,
		"rabbitmq_addr", w.Config().RabbitMQAddr, "rabbitmq_port", w.Config().RabbitMQPort,
		"mongodb_addr", w.Config().MongoDBAddr, "mongodb_port", w.Config().MongoDBPort,
		"redis_addr", w.Config().RedisAddr, "redis_port", w.Config().RedisPort,
	)
	wg.Wait()
	return nil
}

func (w *writeHomeTimelineService) WriteHomeTimeline(ctx context.Context, msg model.Message) error {
	logger := w.Logger(ctx)

	span := trace.SpanFromContext(ctx)
	if trace.SpanContextFromContext(ctx).IsValid() {
		logger.Debug("valid span", "s", span.IsRecording(), "ctx", ctx.Value("TEST"))
	}

	regionLabel := sn_metrics.RegionLabel{Region: w.Config().Region}
	sn_metrics.QueueDurationMs.Get(regionLabel).Put(float64(time.Now().UnixMilli() - msg.NotificationSendTs))

	postIDStr := strconv.FormatInt(msg.PostID, 10)

	result, _, err := w.mongoClientWriteHomeTL.Read(ctx, "posts", postIDStr)
	if err != nil {
		logger.Error("error reading post from mongo", "msg", err.Error())
		return err
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			trace.SpanFromContext(ctx).SetAttributes(
				attribute.Bool("poststorage_consistent_read", false),
			)
			logger.Debug("inconsistency!")
			sn_metrics.Inconsistencies.Get(regionLabel).Inc()
			return nil
		} else {
			logger.Error("error reading post from mongodb", "msg", err.Error())
			return err
		}
	}

	var post model.Post
	err = json.Unmarshal([]byte(result), &post)
	if err != nil {
		errMsg := fmt.Sprintf("post_id: %s not found in mongodb", postIDStr)
		logger.Warn(errMsg)
		return fmt.Errorf(errMsg)
	}

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Bool("poststorage_consistent_read", true),
	)

	logger.Debug("found post! :)", "post_id", post.PostID, "text", post.Text)

	followersID, err := w.socialGraphService.Get().GetFollowers(ctx, msg.ReqID, msg.UserID)
	if err != nil {
		logger.Error("error getting followers from social graph service")
		return err
	}

	logger.Debug("got followers to write to their hometimeline", "num", len(followersID))
	uniqueIDs := make(map[int64]bool, 0)
	for _, followerID := range followersID {
		uniqueIDs[followerID] = true
	}
	for _, userMentionID := range msg.UserMentionIDs {
		uniqueIDs[userMentionID] = true
	}
	value := redis.Z{
		Member: msg.PostID,
		Score:  float64(msg.Timestamp),
	}
	_, err = w.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for id := range uniqueIDs {
			idStr := strconv.FormatInt(id, 10)
			err = w.redisClient.ZAddNX(ctx, idStr, value).Err()
			if err != nil {
				return err
			}
		}
		return nil
	})
	logger.Debug("leaving write home timeline")
	return nil
}

// onReceivedWorker adds the post to all the post's subscribed users (followers, mentioned users, etc)
func (w *writeHomeTimelineService) onReceivedWorker(ctx context.Context, workerid int, body []byte) error {
	logger := w.Logger(ctx)

	var msg model.Message
	err := json.Unmarshal(body, &msg)
	if err != nil {
		logger.Error("error parsing json message", "workerid", workerid, "msg", err.Error())
		return err
	}
	regionLabel := sn_metrics.RegionLabel{Region: w.Config().Region}
	sn_metrics.ReceivedNotifications.Get(regionLabel).Add(1)
	logger.Debug("received rabbitmq message", "workerid", workerid, "post_id", msg.PostID, "msg", msg)

	spanContext, err := sn_trace.ParseSpanContext(msg.SpanContext)
	if err != nil {
		logger.Error("error parsing span context", "workerid", workerid, "msg", err.Error())
		return err
	}

	ctx = trace.ContextWithRemoteSpanContext(ctx, spanContext)

	return w.WriteHomeTimeline(ctx, msg)
}

func (w *writeHomeTimelineService) workerThread(ctx context.Context, workerid int) error {
	logger := w.Logger(ctx)

	var forever chan struct{}
	go func() error {
		for {
			logger.Debug("before read", "region", w.Config().Region)
			msg, lineage, err := w.rabbitClientWriteHomeTL.Read(ctx, "write-home-timeline", "")
			logger.Debug("after read", "region", w.Config().Region)
			if err != nil {
				logger.Error("error consuming queue", "workerid", workerid, "msg", err.Error())
				return err
			}

			logger.Debug("write home timeline service | lineage and message successfully read!", "region", w.Config().Region, "lineage", lineage, "message", msg)

			ctx, err = xcweaver.Transfer(ctx, lineage)

			if err != nil {
				logger.Error("error transfering the lineage to context", "workerid", workerid, "msg", err.Error())
				return err
			}

			err = w.mongoClientWriteHomeTL.Barrier(ctx)
			if err != nil {
				logger.Error("error on barrier", "workerid", workerid, "msg", err.Error())
				return err
			}

			err = w.onReceivedWorker(ctx, workerid, []byte(msg))
			if err != nil {
				logger.Warn("error in worker thread", "msg", err.Error())
			}
		}
	}()
	<-forever
	return nil
}
