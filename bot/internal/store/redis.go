package store

import (
	"context"
	"encoding/json"
	"iter"
	"log/slog"
	"strings"

	"github.com/redis/go-redis/v9"
)

const redisPrefix = "discord_invites:"

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(host, password string) *RedisStore {
	addr := host
	if !strings.Contains(addr, ":") {
		addr += ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &RedisStore{client: client}
}

func (s *RedisStore) Set(invite CachedInvite) error {
	ctx := context.Background()
	data, err := json.Marshal(invite.Data)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, redisPrefix+invite.Code, data, 0).Err()
}

func (s *RedisStore) Delete(code string) error {
	ctx := context.Background()
	return s.client.Del(ctx, redisPrefix+code).Err()
}

func (s *RedisStore) Iterator() iter.Seq2[string, CachedInviteData] {
	return func(yield func(string, CachedInviteData) bool) {
		ctx := context.Background()
		var cursor uint64
		for {
			keys, nextCursor, err := s.client.Scan(ctx, cursor, redisPrefix+"*", 100).Result()
			if err != nil {
				slog.Error("failed to scan redis keys", "error", err)
				return
			}

			for _, key := range keys {
				code := key[len(redisPrefix):]
				val, err := s.client.Get(ctx, key).Result()
				if err != nil {
					slog.Error("failed to get redis key", "key", key, "error", err)
					continue
				}

				var data CachedInviteData
				if err := json.Unmarshal([]byte(val), &data); err != nil {
					slog.Error("failed to parse redis value", "key", key, "error", err)
					continue
				}

				if !yield(code, data) {
					return
				}
			}

			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}
	}
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}
