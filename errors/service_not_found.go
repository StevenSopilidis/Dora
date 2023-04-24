package errors

type ServiceNotFoundError struct{}

func (s *ServiceNotFoundError) Error() string {
	return "service not found"
}
