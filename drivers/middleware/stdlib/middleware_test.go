package stdlib_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	limiter "github.com/hgtpcastro/go-expert-lab-rate-limiter"
	stdlib "github.com/hgtpcastro/go-expert-lab-rate-limiter/drivers/middleware/stdlib"
	"github.com/hgtpcastro/go-expert-lab-rate-limiter/drivers/store/redis"
	libredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	redisTestContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

var redisContainer *redisTestContainer.RedisContainer
var redisURL string

func TestRateLimiterByIPWithSequentialAccess(t *testing.T) {
	is := require.New(t)
	ctx := context.Background()

	setup(ctx, t)
	defer func() {
		tearDown(t)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	request, err := http.NewRequest("GET", "/", nil)

	is.NoError(err)
	is.NotNil(request)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
	})

	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter:redis:sequential-test",
	})
	is.NoError(err)
	is.NotNil(store)

	limiter := limiter.NewLimiter(store, limiter.Rate{
		Limit:  10,
		Period: 1 * time.Minute,
	})

	middleware := stdlib.NewMiddleware(limiter).Handler(handler)
	is.NotZero(middleware)

	success := int64(10)
	clients := int64(100)

	for i := int64(1); i <= clients; i++ {

		resp := httptest.NewRecorder()
		middleware.ServeHTTP(resp, request)

		if i <= success {
			is.Equal(resp.Code, http.StatusOK)
		} else {
			is.Equal(resp.Code, http.StatusTooManyRequests)
		}
	}
}

func TestRateLimiterByIPWithConcurrentAccess(t *testing.T) {
	is := require.New(t)
	ctx := context.Background()

	setup(ctx, t)
	defer func() {
		tearDown(t)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	request, err := http.NewRequest("GET", "/", nil)
	is.NoError(err)
	is.NotNil(request)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
	})

	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter:redis:sequential-test",
	})
	is.NoError(err)
	is.NotNil(store)

	limiter := limiter.NewLimiter(store, limiter.Rate{
		Limit:  10,
		Period: 1 * time.Minute,
	})

	middleware := stdlib.NewMiddleware(limiter).Handler(handler)
	is.NotZero(middleware)

	success := int64(10)
	clients := int64(100)
	counter := int64(0)

	wg := &sync.WaitGroup{}

	for i := int64(1); i <= clients; i++ {
		wg.Add(1)
		go func() {

			resp := httptest.NewRecorder()
			middleware.ServeHTTP(resp, request)

			if resp.Code == http.StatusOK {
				atomic.AddInt64(&counter, 1)
			}

			wg.Done()
		}()
	}

	wg.Wait()
	is.Equal(success, atomic.LoadInt64(&counter))
}

func TestRateLimiterByTokenWithSequentialAccess(t *testing.T) {
	apiKey := "any-api-key"
	is := require.New(t)
	ctx := context.Background()

	setup(ctx, t)
	defer func() {
		tearDown(t)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	request, err := http.NewRequest("GET", "/", nil)
	request.Header.Set("API_KEY", apiKey)
	is.NoError(err)
	is.NotNil(request)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
	})

	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter:redis:sequential-test",
	})
	is.NoError(err)
	is.NotNil(store)

	limiter := limiter.NewLimiter(store, limiter.Rate{
		Limit:  10,
		Period: 1 * time.Minute,
	})

	middleware := stdlib.NewMiddleware(
		limiter,
		stdlib.WithKeyGetter(stdlib.WithTokenKeyGetter(limiter)),
	).Handler(handler)
	is.NotZero(middleware)

	success := int64(10)
	clients := int64(100)

	for i := int64(1); i <= clients; i++ {

		resp := httptest.NewRecorder()
		middleware.ServeHTTP(resp, request)

		if i <= success {
			is.Equal(resp.Code, http.StatusOK)
		} else {
			is.Equal(resp.Code, http.StatusTooManyRequests)
		}
	}
}

func TestRateLimiterByTokenWithConcurrentAccess(t *testing.T) {
	apiKey := "any-api-key"
	is := require.New(t)
	ctx := context.Background()

	setup(ctx, t)
	defer func() {
		tearDown(t)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	request, err := http.NewRequest("GET", "/", nil)
	request.Header.Set("API_KEY", apiKey)
	is.NoError(err)
	is.NotNil(request)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
	})

	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter:redis:sequential-test",
	})
	is.NoError(err)
	is.NotNil(store)

	limiter := limiter.NewLimiter(store, limiter.Rate{
		Limit:  10,
		Period: 1 * time.Minute,
	})

	middleware := stdlib.NewMiddleware(
		limiter,
		stdlib.WithKeyGetter(stdlib.WithTokenKeyGetter(limiter)),
	).Handler(handler)
	is.NotZero(middleware)

	success := int64(10)
	clients := int64(100)
	counter := int64(0)

	wg := &sync.WaitGroup{}

	for i := int64(1); i <= clients; i++ {
		wg.Add(1)
		go func() {

			resp := httptest.NewRecorder()
			middleware.ServeHTTP(resp, request)

			if resp.Code == http.StatusOK {
				atomic.AddInt64(&counter, 1)
			}

			wg.Done()
		}()
	}

	wg.Wait()
	is.Equal(success, atomic.LoadInt64(&counter))
}

func TestRateLimiterByTokenAndIPWithSequentialAccess(t *testing.T) {
	apiKey := "any-api-key"
	is := require.New(t)
	ctx := context.Background()

	setup(ctx, t)
	defer func() {
		tearDown(t)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
	})

	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter:redis:sequential-test",
	})
	is.NoError(err)
	is.NotNil(store)

	limiterByIP := limiter.NewLimiter(store, limiter.Rate{
		Limit:  10,
		Period: 1 * time.Minute,
	})

	middlewareByIP := stdlib.NewMiddleware(
		limiterByIP,
		stdlib.WithKeyGetter(stdlib.WithTokenAndIPKeyGetter(limiterByIP)),
	).Handler(handler)
	is.NotZero(middlewareByIP)

	limiterByToken := limiter.NewLimiter(store, limiter.Rate{
		Limit:  15,
		Period: 1 * time.Minute,
	})

	middlewareByToken := stdlib.NewMiddleware(
		limiterByToken,
		stdlib.WithKeyGetter(stdlib.WithTokenAndIPKeyGetter(limiterByToken)),
	).Handler(handler)
	is.NotZero(middlewareByToken)

	successByIP := int64(10)
	successByToken := int64(15)
	clients := int64(100)
	j := int64(0)
	k := int64(0)

	requestByIP, err := http.NewRequest("GET", "/", nil)
	is.NoError(err)
	is.NotNil(requestByIP)

	requestByToken, err := http.NewRequest("GET", "/", nil)
	requestByToken.Header.Set("API_KEY", apiKey)
	is.NoError(err)
	is.NotNil(requestByToken)

	for i := int64(1); i <= clients; i++ {

		resp := httptest.NewRecorder()
		if i%2 == 0 {
			middlewareByIP.ServeHTTP(resp, requestByIP)
			j++
			if j <= successByIP {
				is.Equal(resp.Code, http.StatusOK)
			} else {
				is.Equal(resp.Code, http.StatusTooManyRequests)
			}
		} else {
			middlewareByToken.ServeHTTP(resp, requestByToken)
			k++
			if k <= successByToken {
				is.Equal(resp.Code, http.StatusOK)
			} else {
				is.Equal(resp.Code, http.StatusTooManyRequests)
			}
		}
	}
}

func TestRateLimiterByTokenAndIPWithConcurrentAccess(t *testing.T) {
	apiKey := "any-api-key"
	is := require.New(t)
	ctx := context.Background()

	setup(ctx, t)
	defer func() {
		tearDown(t)
	}()

	client, err := newRedisClient(redisURL)
	is.NoError(err)
	is.NotNil(client)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
	})

	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter:redis:sequential-test",
	})
	is.NoError(err)
	is.NotNil(store)

	limiterByIP := limiter.NewLimiter(store, limiter.Rate{
		Limit:  10,
		Period: 1 * time.Minute,
	})

	middlewareByIP := stdlib.NewMiddleware(
		limiterByIP,
		stdlib.WithKeyGetter(stdlib.WithTokenAndIPKeyGetter(limiterByIP)),
	).Handler(handler)
	is.NotZero(middlewareByIP)

	limiterByToken := limiter.NewLimiter(store, limiter.Rate{
		Limit:  15,
		Period: 1 * time.Minute,
	})

	middlewareByToken := stdlib.NewMiddleware(
		limiterByToken,
		stdlib.WithKeyGetter(stdlib.WithTokenAndIPKeyGetter(limiterByToken)),
	).Handler(handler)
	is.NotZero(middlewareByToken)

	clients := int64(100)
	successByIP := int64(10)
	successByToken := int64(15)
	counterByIP := int64(0)
	counterByToken := int64(0)

	requestByIP, err := http.NewRequest("GET", "/", nil)
	is.NoError(err)
	is.NotNil(requestByIP)

	requestByToken, err := http.NewRequest("GET", "/", nil)
	requestByToken.Header.Set("API_KEY", apiKey)
	is.NoError(err)
	is.NotNil(requestByToken)

	wg := &sync.WaitGroup{}

	for i := int64(1); i <= clients; i++ {
		wg.Add(1)
		go func(i int64) {

			resp := httptest.NewRecorder()
			if i%2 == 0 {
				middlewareByIP.ServeHTTP(resp, requestByIP)

				if resp.Code == http.StatusOK {
					atomic.AddInt64(&counterByIP, 1)
				}
			} else {
				middlewareByToken.ServeHTTP(resp, requestByToken)

				if resp.Code == http.StatusOK {
					atomic.AddInt64(&counterByToken, 1)
				}
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
	is.Equal(successByIP, atomic.LoadInt64(&counterByIP))
	is.Equal(successByToken, atomic.LoadInt64(&counterByToken))
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

func runRedisWithTestContainer(ctx context.Context, t testing.TB) *redisTestContainer.RedisContainer {
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
	redisContainer = runRedisWithTestContainer(ctx, t)

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
