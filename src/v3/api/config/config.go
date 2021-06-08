package config

// DB config
type DB struct {
	Dialect  string `envconfig:"DIALECT" default:"mysql"`
	Host     string `envconfig:"DB_HOST"`
	Username string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	Name     string `envconfig:"DB_DATABASE"`
	Charset  string `envconfig:"CHARSET" default:"utf8"`
	Timezone string `envconfig:"TIMEZONE"`
}
