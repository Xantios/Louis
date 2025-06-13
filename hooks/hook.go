package hooks

type Hook interface {
	Send(string) error
}
