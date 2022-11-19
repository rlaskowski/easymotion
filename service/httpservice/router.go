package httpservice

type Router interface {
	StreamVideo() <-chan VideoResponse
}

type VideoResponse struct {
	Data []byte
	Err  error
}
