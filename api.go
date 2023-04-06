package srouter

type Api interface {
	Start(port string) error
	Stop()
}
