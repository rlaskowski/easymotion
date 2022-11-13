package easymotion

type Runner interface {
	RegisterServices() error
	Run() error
	Close() error
}
