package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB
	RabbitURL string `env:"RABBITMQ_URL" env-default:"amqp://guest:guest@localhost:5672/"`
}
type DB struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-default:"5488"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
	SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}
