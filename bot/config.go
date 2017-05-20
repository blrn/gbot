package bot

import "fmt"

const (
	DEFAULT_SERVER_PORT  uint = 8000
	DEFAULT_SERVER_IP         = "127.0.0.1"
	DEFAULT_CHANNEL_NAME      = "town-square"
)

type Config struct {
	HttpEnabled        bool
	HttpServerPort     uint
	HttpServerIp       string
	DefaultChannelName string
}

func defaultConfig() Config {
	return Config{
		HttpEnabled:        true,
		HttpServerPort:     DEFAULT_SERVER_PORT,
		HttpServerIp:       DEFAULT_SERVER_IP,
		DefaultChannelName: DEFAULT_CHANNEL_NAME,
	}
}

func (c *Config) getHttpAddress() string {
	port := c.HttpServerPort
	host := c.HttpServerIp
	if port == 0 {
		port = DEFAULT_SERVER_PORT
	}
	if len(host) == 0 {
		host = DEFAULT_SERVER_IP
	}
	return fmt.Sprintf("%s:%d", host, port)
}
