package registry

import (
	"context"
)

type Service struct {
	// Address of the service
	Addr string
	// Port where the service runs
	Port uint16
	// health check url where to check health of the service
	// It should be a get endpoint and the return status
	// code is 2xx its healthy
	HealthCheckUrl string
	// how often do we health check the service
	HealthCheckTimeInterval uint8
	// whether the service is alive
	IsAlive bool
}

type Registry interface {
	// appends service on a registry
	// @param name: name of service
	// @param service: Service details
	// returns bool indicating whether or not the append was successfull
	Append(ctx context.Context, name string, service Service) error
	// removes service from registry
	// @param name: name of the service we want to remove
	// returns bool indicating wether the removal was successful
	Remove(ctx context.Context, name string) error
	// retrives service details from registry
	// @param name: name of service to retreive
	// returns the Service or nil if it does not exist
	Get(ctx context.Context, name string) (error, Service)
	// function that will check the health of the service
	// by making a request to the service itself
	// if its not alive it will Set it as not alive
	// @param service: service to check
	// retutns ServiceNotFoundError if service is not found
	// of if service its not health it will return
	// ServiceNotHealthy soif it ServiceNotHealthy is returned
	// we can set the service in the db
	CheckHealth(service string) (error, bool)
}
