package service

type Runner interface {
	Start() error
	Stop() error
}
