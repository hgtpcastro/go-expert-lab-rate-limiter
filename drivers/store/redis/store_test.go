package redis_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	limiter "github.com/hgtpcastro/go-expert-lab-rate-limiter"
	"github.com/hgtpcastro/go-expert-lab-rate-limiter/drivers/store/redis"
	"github.com/hgtpcastro/go-expert-lab-rate-limiter/drivers/store/tests"
	libredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	redisTestContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

var redisContainer *redisTestContainer.RedisContainer
var redisURL string

func TestRedisStoreSequentialAccess(t *testing.T) {
	is := require.New(t)
	ctx := context.Background()

	setup(ctx, t)
	defer func() {
		tearDown(t)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter:redis:sequential-test",
	})
	is.NoError(err)
	is.NotNil(store)

	tests.TestStoreSequentialAccess(t, store)
}

func TestRedisStoreConcurrentAccess(t *testing.T) {
	is := require.New(t)
	ctx := context.Background()

	setup(ctx, t)
	defer func() {
		tearDown(t)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter:redis:concurrent-test",
	})
	is.NoError(err)
	is.NotNil(store)

	tests.TestStoreConcurrentAccess(t, store)
}
func TestRedisClientExpiration(t *testing.T) {
	is := require.New(t)
	ctx := context.Background()

	setup(ctx, t)
	defer func() {
		tearDown(t)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	key := "my-key"
	value := 10
	keyNoExpiration := -1 * time.Nanosecond
	keyNotExist := -2 * time.Nanosecond

	delCmd := client.Del(ctx, key)
	_, err = delCmd.Result()
	is.NoError(err)

	expCmd := client.PTTL(ctx, key)
	ttl, err := expCmd.Result()
	is.NoError(err)
	is.Equal(keyNotExist, ttl)

	setCmd := client.Set(ctx, key, value, 0)
	_, err = setCmd.Result()
	is.NoError(err)

	expCmd = client.PTTL(ctx, key)
	ttl, err = expCmd.Result()
	is.NoError(err)
	is.Equal(keyNoExpiration, ttl)

	setCmd = client.Set(ctx, key, value, 1*time.Second)
	_, err = setCmd.Result()
	is.NoError(err)

	time.Sleep(100 * time.Millisecond)

	expCmd = client.PTTL(ctx, key)
	ttl, err = expCmd.Result()
	is.NoError(err)

	expected := int64(0)
	actual := int64(ttl)
	is.Greater(actual, expected)
}

func BenchmarkRedisStoreSequentialAccess(b *testing.B) {
	is := require.New(b)
	ctx := context.Background()

	setup(ctx, b)
	defer func() {
		tearDown(b)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter:redis:sequential-benchmark",
	})
	is.NoError(err)
	is.NotNil(store)

	tests.BenchmarkStoreSequentialAccess(b, store)
}

func BenchmarkRedisStoreConcurrentAccess(b *testing.B) {
	is := require.New(b)
	ctx := context.Background()

	setup(ctx, b)
	defer func() {
		tearDown(b)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter:redis:concurrent-benchmark",
	})
	is.NoError(err)
	is.NotNil(store)

	tests.BenchmarkStoreConcurrentAccess(b, store)
}

func newRedisClient(redisURL string) (*libredis.Client, error) {
	url := fmt.Sprintf("%s/0", redisURL)

	if os.Getenv("REDIS_URL") != "" {
		url = os.Getenv("REDIS_URL")
	}

	opt, err := libredis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := libredis.NewClient(opt)
	return client, nil
}

func runRedis(ctx context.Context, t testing.TB) *redisTestContainer.RedisContainer {
	redisContainer, err := redisTestContainer.Run(
		ctx,
		"redis:7.2.4-alpine3.19",
		redisTestContainer.WithSnapshotting(10, 1),
		redisTestContainer.WithLogLevel(redisTestContainer.LogLevelVerbose),
		testcontainers.WithHostPortAccess(6379),
		// testcontainers.WithWaitStrategy(
		// 	wait.ForLog("database system is ready to accept connections").
		// 		WithOccurrence(2).
		// 		WithStartupTimeout(10*time.Second)),
		//redisTestContainer.WithConfigFile(filepath.Join("test.data", "redis7.conf")),
	)

	// defer func() {
	// 	if err := testcontainers.TerminateContainer(redisContainer); err != nil {
	// 		t.Fatalf("failed to terminate container: %s", err)
	// 	}
	// }()

	if err != nil {
		t.Fatalf("failed to start container: %s", err)

	}

	return redisContainer
}

func setup(ctx context.Context, t testing.TB) {
	redisContainer = runRedis(ctx, t)

	connectionString, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	redisURL = connectionString
}

func tearDown(t testing.TB) {

	if err := testcontainers.TerminateContainer(redisContainer); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}

}
