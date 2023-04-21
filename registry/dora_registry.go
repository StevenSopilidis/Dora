package registry

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis"
)

const (
	maxOperatingWaitTime = 5
)

type RedisRegistryClient struct {
	client *redis.Client
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
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRegistryClient) Get(ctx context.Context, name string) (error, Service) {
	ctx, cancel := context.WithTimeout(ctx, maxOperatingWaitTime)
	defer cancel()
	val, err := r.client.WithContext(ctx).Get(name).Result()
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
