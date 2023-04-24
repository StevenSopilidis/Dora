package registry

import (
	"context"
	"testing"

	"github.com/go-redis/redis"
	e "github.com/stevensopilidis/dora/errors"
	"github.com/stretchr/testify/require"
)

func setUpTest(t *testing.T) (
	client *RedisRegistryClient,
	teardown func(),
) {
	t.Helper()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	client = &RedisRegistryClient{
		client: redisClient,
	}
	teardown = func() {
		redisClient.Close()
	}
	_, err := client.client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client, teardown
}

func TestRedisRegis(t *testing.T) {
	r, teardown := setUpTest(t)
	defer teardown()
	for scenario, fn := range map[string]func(
		t *testing.T,
		r *RedisRegistryClient,
	){
		"append entry to registry":   testAppend,
		"get entry from registry":    testGet,
		"remove entry from registry": testRemove,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, r)
		})
	}
}

func testAppend(t *testing.T, r *RedisRegistryClient) {
	service := Service{
		Addr: "127.0.0.1",
		Port: 80,
	}
	service_name := "test_service"
	err := r.Append(context.Background(), service_name, service)
	require.NoError(t, err)
}

func testGet(t *testing.T, r *RedisRegistryClient) {
	service := Service{
		Addr: "127.0.0.1",
		Port: 80,
	}
	service_name := "my_service"
	err := r.Append(context.Background(), service_name, service)
	require.NoError(t, err)
	err, data := r.Get(context.Background(), service_name)
	require.NoError(t, err)
	require.Equal(t, service.Addr, data.Addr)
	require.Equal(t, service.Port, data.Port)
}

func testRemove(t *testing.T, r *RedisRegistryClient) {
	service := Service{
		Addr: "127.0.0.1",
		Port: 80,
	}
	service_name := "test_service"
	err := r.Append(context.Background(), service_name, service)
	require.NoError(t, err)
	err = r.Remove(context.Background(), service_name)
	require.NoError(t, err)
	// redis.Nil error gets returned when key not valid
	err, _ = r.Get(context.Background(), service_name)
	require.Equal(t, err, &e.ServiceNotFoundError{})
}
