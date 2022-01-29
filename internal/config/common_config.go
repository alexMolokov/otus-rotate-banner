package config

type LoggerConf struct {
	Level    string `config:"logger-level"`
	Encoding string `config:"logger-encoding"`
	Output   string `config:"logger-output"`
}

type DBConf struct {
	Driver            string `config:"db-driver"`
	Host              string `config:"db-host"`
	Port              int    `config:"db-port"`
	Name              string `config:"db-name"`
	User              string `config:"db-user"`
	Password          string `config:"db-password"`
	MaxConnectionPool int    `config:"db-maxConnectionPool"`
	SslMode           string `config:"db-sslMode"`
}

type GRPCConf struct {
	Host string `config:"grpc-host"`
	Port int    `config:"grpc-port"`
}

type QueueConf struct {
	URI          string `config:"queue-uri"`
	Name         string `config:"queue-name"`
	ExchangeName string `config:"queue-exchangeName"`
	ExchangeType string `config:"queue-exchangeType"`
	// When setting up the channel after a channel exception in seconds
	ReInitDelay int `config:"queue-reInitDelay"`
	// When reconnecting to the server after connection failure in seconds
	ReconnectDelay int `config:"queue-reconnectDelay"`
	// resending messages the server didn't confirm in seconds
	ResendDelay int `config:"queue-resendDelay"`
}
