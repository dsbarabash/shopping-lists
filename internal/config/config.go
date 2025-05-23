package config

type Config struct {
	Host   string
	Port   int
	Secret string
}

var Cfg Config

func NewConfig() *Config {
	Cfg := Config{
		Host:   "0.0.0.0",
		Port:   8989,
		Secret: "sadffhu8210jt678s7ghb8dnde4r4eergcnu8ndennot8u",
	}
	return &Cfg
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
	}
}

type MongoConfig struct {
	Host string
	Port string
}

func NewMongoConfig() *MongoConfig {
	return &MongoConfig{
		Host: "localhost",
		Port: "27017",
	}
}

type RedisConfig struct {
	Host string
	Port string
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host: "localhost",
		Port: "6379",
	}
}
