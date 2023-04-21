package registry

import "context"

type Service struct {
	Addr string
	Port uint16
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
}
