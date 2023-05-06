package registry

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	e "github.com/stevensopilidis/dora/errors"
)

const (
	maxOperatingWaitTime = 5
)

type RedisRegistryClient struct {
	client *redis.Client
}

var registry *RedisRegistryClient = nil

func InitRegistry(opts *redis.Options) {
	registry = createRedisRegistryClient(opts)
}

func GetRedisRegistryClient() *RedisRegistryClient {
	if registry == nil {
		log.Panic("Registry was not initialized")
		return nil
	}
	return registry
}

func createRedisRegistryClient(opts *redis.Options) *RedisRegistryClient {
	return &RedisRegistryClient{
		client: redis.NewClient(opts),
	}
}

func (r *RedisRegistryClient) Append(ctx context.Context, name string, service Service) error {
	ctx, cancel := context.WithTimeout(ctx, maxOperatingWaitTime)
	defer cancel()
	json, err := json.Marshal(service)
	if err != nil {
		return err
	}
	result := r.client.WithContext(ctx).Set(name, json, 0)
	if err = result.Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRegistryClient) Remove(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, maxOperatingWaitTime)
	defer cancel()
	err := r.client.WithContext(ctx).Del(name).Err()
	if err == redis.Nil {
		return &e.ServiceNotFoundError{}
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRegistryClient) Get(ctx context.Context, name string) (error, Service) {
	ctx, cancel := context.WithTimeout(ctx, maxOperatingWaitTime)
	defer cancel()
	val, err := r.client.WithContext(ctx).Get(name).Result()
	if err == redis.Nil {
		return &e.ServiceNotFoundError{}, Service{}
	}
	if err != nil {
		return err, Service{}
	}
	var service Service
	err = json.Unmarshal([]byte(val), &service)
	if err != nil {
		return err, Service{}
	}
	return nil, service
}

func (r *RedisRegistryClient) CheckHealth(name string) error {
	err, service := r.Get(context.Background(), name)
	if err != nil {
		return &e.ServiceNotFoundError{}
	}
	_, err = http.Get(service.HealthCheckUrl)
	if err != nil {
		return &e.ServiceUnhealthyError{}
	}
	return err
}
