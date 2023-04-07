package srouter

type Api interface {
	Start(addr string) error
	Stop()
}
