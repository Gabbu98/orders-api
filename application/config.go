package application

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddress	string
	MongoAddress	string
	ServerPort 		uint16
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddress: "localhost:6379",
		MongoAddress: "mongodb://localhost:27017",
		ServerPort: 3000,
	}

	if redisAddr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddress = redisAddr
	}

	if mongoAddr, exists := os.LookupEnv("MONGO_ADDR"); exists {
		cfg.MongoAddress = mongoAddr
	}

	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	return cfg
}