package clients

type Client interface {
	Backup() error
	Update() (bool, string, error)
}
