package easymotion

type Runner interface {
	RegisterServices()
	Run() error
	Close() error
}
