package errors

type InvalidArgument struct {
	Message string
}

func (e *InvalidArgument) Error() string {
	return "Invalid Argument Passed on " + e.Message
}
