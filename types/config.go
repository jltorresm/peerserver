package types

type Config struct {
	Server *server
}

type server struct {
	Url string
}
