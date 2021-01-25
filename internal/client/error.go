package client

type ClientError struct {
	Code int
	msg  string
}

type ServerError struct {
	Code int
	msg  string
}

func (e ClientError) Error() string {
	return e.msg
}

func (e ServerError) Error() string {
	return e.msg
}
