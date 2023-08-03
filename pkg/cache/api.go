package cache

import (
	"context"
)

func Del(ctx context.Context, db int, key string) error {
	rdb, err := RedisConnFactory(db)
	if err != nil {
		return err
	}
	return rdb.Del(ctx, key).Err()
}
