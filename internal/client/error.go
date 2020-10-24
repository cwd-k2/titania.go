package client

type TitaniaClientError struct {
	Code int
	Err  error
}

func (e TitaniaClientError) Error() string {
	return e.Err.Error()
}
