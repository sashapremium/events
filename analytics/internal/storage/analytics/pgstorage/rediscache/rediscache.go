package rediscache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	svc "github.com/sashapremium/events/analytics/internal/services/analyticsService"
)

type RedisCache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *RedisCache {
	return &RedisCache{rdb: rdb}
}

func totalsKey(postID uint64) string { return "analytics:totals:" + strconv.FormatUint(postID, 10) }
func deltaKey(postID uint64) string  { return "analytics:delta:" + strconv.FormatUint(postID, 10) }
func uniqKey(postID uint64) string   { return "analytics:uniq:" + strconv.FormatUint(postID, 10) }

func topKey(metric string) string { return "analytics:top:" + metric }
func dirtyKey() string            { return "analytics:dirty" }

func lastSyncedKey() string { return "analytics:last_synced_at" }

func (c *RedisCache) readTotalsHash(ctx context.Context, key string) (svc.PostTotals, bool, error) {
	m, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return svc.PostTotals{}, false, err
	}
	if len(m) == 0 {
		return svc.PostTotals{}, false, nil
	}

	parse := func(field string) (int64, error) {
		v, ok := m[field]
		if !ok || v == "" {
			return 0, nil
		}
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse %s=%q: %w", field, v, err)
		}
		return n, nil
	}

	views, err := parse("views")
	if err != nil {
		return svc.PostTotals{}, false, err
	}
	likes, err := parse("likes")
	if err != nil {
		return svc.PostTotals{}, false, err
	}
	comments, err := parse("comments")
	if err != nil {
		return svc.PostTotals{}, false, err
	}
	reposts, err := parse("reposts")
	if err != nil {
		return svc.PostTotals{}, false, err
	}
	unique, err := parse("unique_users")
	if err != nil {
		return svc.PostTotals{}, false, err
	}

	return svc.PostTotals{
		Views:       uint64(views),
		Likes:       uint64(likes),
		Comments:    uint64(comments),
		Reposts:     uint64(reposts),
		UniqueUsers: uint64(unique),
	}, true, nil
}

func (c *RedisCache) incrHashFields(ctx context.Context, key string, d svc.TotalsDelta) error {
	pipe := c.rdb.Pipeline()

	if d.Views != 0 {
		pipe.HIncrBy(ctx, key, "views", d.Views)
	}
	if d.Likes != 0 {
		pipe.HIncrBy(ctx, key, "likes", d.Likes)
	}
	if d.Comments != 0 {
		pipe.HIncrBy(ctx, key, "comments", d.Comments)
	}
	if d.Reposts != 0 {
		pipe.HIncrBy(ctx, key, "reposts", d.Reposts)
	}
	if d.UniqueUsers != 0 {
		pipe.HIncrBy(ctx, key, "unique_users", d.UniqueUsers)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (c *RedisCache) IncrTotals(ctx context.Context, postID uint64, delta svc.TotalsDelta) error {
	return c.incrHashFields(ctx, totalsKey(postID), delta)
}

func (c *RedisCache) GetTotals(ctx context.Context, postID uint64) (svc.PostTotals, bool, error) {
	return c.readTotalsHash(ctx, totalsKey(postID))
}

func (c *RedisCache) GetDelta(ctx context.Context, postID uint64) (svc.TotalsDelta, bool, error) {
	t, ok, err := c.readTotalsHash(ctx, deltaKey(postID))
	if err != nil {
		return svc.TotalsDelta{}, false, err
	}
	if !ok {
		return svc.TotalsDelta{}, false, nil
	}

	return svc.TotalsDelta{
		Views:       int64(t.Views),
		Likes:       int64(t.Likes),
		Comments:    int64(t.Comments),
		Reposts:     int64(t.Reposts),
		UniqueUsers: int64(t.UniqueUsers),
	}, true, nil
}

func (c *RedisCache) ResetDelta(ctx context.Context, postID uint64) error {
	return c.rdb.Del(ctx, deltaKey(postID)).Err()
}

func (c *RedisCache) AddUniqueUser(ctx context.Context, postID uint64, userHash string) (bool, error) {
	added, err := c.rdb.SAdd(ctx, uniqKey(postID), userHash).Result()
	if err != nil {
		return false, err
	}
	return added > 0, nil
}

func (c *RedisCache) IncrTop(ctx context.Context, metric string, postID uint64, inc int64) error {
	member := strconv.FormatUint(postID, 10)
	return c.rdb.ZIncrBy(ctx, topKey(metric), float64(inc), member).Err()
}

func (c *RedisCache) GetTop(ctx context.Context, metric string, limit uint32) ([]svc.TopItem, error) {
	if limit == 0 {
		return []svc.TopItem{}, nil
	}

	items, err := c.rdb.ZRevRangeWithScores(ctx, topKey(metric), 0, int64(limit)-1).Result()
	if err != nil {
		return nil, err
	}

	out := make([]svc.TopItem, 0, len(items))
	for _, z := range items {
		s, ok := z.Member.(string)
		if !ok {
			continue
		}
		postID, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			continue
		}
		out = append(out, svc.TopItem{
			PostID: postID,
			Value:  uint64(z.Score),
		})
	}
	return out, nil
}

func (c *RedisCache) MarkDirty(ctx context.Context, postID uint64) error {
	return c.rdb.SAdd(ctx, dirtyKey(), strconv.FormatUint(postID, 10)).Err()
}

func (c *RedisCache) GetDirtyBatch(ctx context.Context, limit int) ([]uint64, error) {
	if limit <= 0 {
		return []uint64{}, nil
	}

	members, err := c.rdb.SPopN(ctx, dirtyKey(), int64(limit)).Result()
	if err != nil {
		return nil, err
	}

	out := make([]uint64, 0, len(members))
	for _, s := range members {
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			continue
		}
		out = append(out, id)
	}
	return out, nil
}

func (c *RedisCache) SetLastSyncedAt(ctx context.Context, ts string) error {
	return c.rdb.Set(ctx, lastSyncedKey(), ts, 0).Err()
}

func (c *RedisCache) GetLastSyncedAt(ctx context.Context) (string, bool, error) {
	val, err := c.rdb.Get(ctx, lastSyncedKey()).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *RedisCache) ExpireUnique(ctx context.Context, postID uint64, ttl time.Duration) error {
	return c.rdb.Expire(ctx, uniqKey(postID), ttl).Err()
}
