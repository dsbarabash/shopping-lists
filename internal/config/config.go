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
