package api

type Api interface {
	Start(port string) error
	Stop()
}
