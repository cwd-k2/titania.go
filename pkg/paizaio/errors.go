package paizaio

type ClientError struct {
	Code int
	msg  string
}

type ServerError struct {
	Code int
	msg  string
}

type RunnerError struct {
	msg string
}

func (e ClientError) Error() string {
	return e.msg
}

func (e ServerError) Error() string {
	return e.msg
}

func (e RunnerError) Error() string {
	return e.msg
}
