package client

type Client struct {
}

type Getter interface {
	Get() error
}
