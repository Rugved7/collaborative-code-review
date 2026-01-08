package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Rugved7/collaborative-code-review/internal/common/config"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis initializes Redis client with connection pooling
func InitRedis(cfg *config.Config) (*redis.Client, error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		PoolSize:     10,
		MinIdleConns: 2,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// verify connections
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	log.Printf("[Cache] Redis connected successfully at %s", cfg.RedisAddr)
	return RedisClient, nil
}

// Close redis connection
func Close() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// HealthCheck verifies Redis connectivity
func HealthCheck() error {
	if RedisClient == nil {
		return fmt.Errorf("redis not initialized")
	}
	return RedisClient.Ping(ctx).Err()
}

// PUB-SUB HELPERS

// PublishMessage publishes a message to a Redis channel
func PublishMessage(channel string, message interface{}) error {
	return RedisClient.Publish(ctx, channel, message).Err()
}

// SubscribeToChannel subscribes to a Redis channel and returns PubSub instance
func SubscribeToChannel(channels ...string) *redis.PubSub {
	return RedisClient.Subscribe(ctx, channels...)
}

// STREAM HELPERS

// AddToStream adds a message to Redis Stream (for event sourcing)
func AddToStreams(stream string, values map[string]interface{}) (string, error) {
	id, err := RedisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}).Result()
	return id, err
}

// ReadFromStream reads messages from Redis Stream with consumer group
func ReadFromStream(group, consumer, stream string, count int64) ([]redis.XStream, error) {
	return RedisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    count,
		Block:    0, // block indefinetely
	}).Result()
}

// CreateConsumerGroup creates a consumer group for a stream
func CreateConsumerGroup(stream, group string) error {
	return RedisClient.XGroupCreateMkStream(ctx, stream, group, "0").Err() // Use MKSTREAM to create stream if it doesn't exist
}

// AckStreamMessage acknowledges a message in a stream
func AckStreamMessage(stream, group, messageID string) error {
	return RedisClient.XAck(ctx, stream, group, messageID).Err()
}

// SET OPERATIONS FOR PRESENCE TRACKING

// AddToSet adds a member to a Redis set with optional TTL
func AddToSet(key string, member interface{}, ttl time.Duration) error {
	pipe := RedisClient.Pipeline()
	pipe.SAdd(ctx, key, member)
	if ttl > 0 {
		pipe.Expire(ctx, key, ttl)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// RemoveFromSet removes a member from a Redis set
func RemoveFromSet(key string, member interface{}) error {
	return RedisClient.SRem(ctx, key, member).Err()
}

// GetSetMembers retrieves all members from a Redis set
func GetSetMembers(key string) ([]string, error) {
	return RedisClient.SMembers(ctx, key).Result()
}

// STRINGS OPERATIONS

// SetWithExpiry sets a key-value pair with expiration
func SetWithExpiry(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func Get(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// Delete deletes keys
func Delete(keys ...string) error {
	return RedisClient.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func Exists(keys ...string) (int64, error) {
	return RedisClient.Exists(ctx, keys...).Result()
}
