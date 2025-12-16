package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	RabbitMQ RabbitMQConfig
	Zipkin   ZipkinConfig
}

type ServerConfig struct {
	Port        string
	Environment string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type ZipkinConfig struct {
	Endpoint string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8082"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
			User:     getEnv("RABBITMQ_USER", "guest"),
			Password: getEnv("RABBITMQ_PASS", "guest"),
		},
		Zipkin: ZipkinConfig{
			Endpoint: getEnv("ZIPKIN_ENDPOINT", "http://localhost:9411/api/v2/spans"),
		},
	}
}

// GetRabbitMQURL returns RabbitMQ connection URL
func (c *RabbitMQConfig) GetRabbitMQURL() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		c.User, c.Password, c.Host, c.Port,
	)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

