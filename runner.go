package easymotion

type Runner interface {
	Run() error
	Close() error
}
