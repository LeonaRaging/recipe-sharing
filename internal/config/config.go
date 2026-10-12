package config

type Config struct {
	HTTPAddress string
	DatabaseURL string
}

func Load() Config {
	return Config{
		HTTPAddress: "localhost:8080",
		DatabaseURL: "host=localhost user=postgres password=postgres dbname=recipe_sharing port=5432 sslmode=disable",
	}
}
