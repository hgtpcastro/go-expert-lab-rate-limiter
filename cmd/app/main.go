package main

import (
	"fmt"
	"log"
	"net/http"

	limiter "github.com/hgtpcastro/go-expert-lab-rate-limiter"
	"github.com/hgtpcastro/go-expert-lab-rate-limiter/config"
	mhttp "github.com/hgtpcastro/go-expert-lab-rate-limiter/drivers/middleware/stdlib"
	stdlib "github.com/hgtpcastro/go-expert-lab-rate-limiter/drivers/middleware/stdlib"
	sredis "github.com/hgtpcastro/go-expert-lab-rate-limiter/drivers/store/redis"
	libredis "github.com/redis/go-redis/v9"
)

func main() {
	//cfg, err := config.Load("./deployments/docker-compose") // <- Use em tempo de execução
	cfg, err := config.Load(".") // <- Use para debug | docker
	if err != nil {
		panic(err)
	}

	//redisUrl := "redis://rate-limiter-redis:6379/0"
	redisUrl := fmt.Sprintf("redis://%v:%v/%v", cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	option, err := libredis.ParseURL(redisUrl)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := libredis.NewClient(option)

	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "limiter_http_example",
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	rateByIP := limiter.NewRate(int64(cfg.RateMaxRequestsByIP), cfg.RatePeriodWindowSeconds)
	rateByToken := limiter.NewRate(int64(cfg.RateMaxRequestsByToken), cfg.RatePeriodWindowSeconds)

	limiterByIP := limiter.NewLimiter(store, rateByIP)
	limiterByToken := limiter.NewLimiter(store, rateByToken)

	middlewareByIP := mhttp.NewMiddleware(
		limiterByIP,
		stdlib.WithKeyGetter(stdlib.WithIPKeyGetter(limiterByIP)),
	)

	middlewareByToken := mhttp.NewMiddleware(
		limiterByToken,
		stdlib.WithKeyGetter(stdlib.WithTokenKeyGetter(limiterByToken)),
	)

	http.Handle("/", middlewareByToken.Handler(middlewareByIP.Handler(http.HandlerFunc(index))))
	//http.Handle("/", middlewareByIP.Handler(middlewareByToken.Handler(http.HandlerFunc(index))))
	fmt.Println(fmt.Sprintf("Server is running on port %d...", cfg.AppPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.AppPort), nil))

}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write([]byte(`{"message": "ok"}`))
	if err != nil {
		log.Fatal(err)
	}
}
