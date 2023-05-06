package main

import (
	"github.com/go-redis/redis"
	r "github.com/stevensopilidis/dora/registry"
	s "github.com/stevensopilidis/dora/server"
)

func main() {
	r.CreateRedisRegistryClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 200,
	})
	s.InitServer()
}
