package errors

type ServiceUnhealthyError struct{}

func (s *ServiceUnhealthyError) Error() string {
	return "Service not healthy"
}
