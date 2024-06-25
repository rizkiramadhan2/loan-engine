package config

// Config struct
type Config struct {
	App       AppConfig       `yaml:"app"`
	Databases DatabasesConfig `yaml:"databases"`
}

// AppConfig struct
type AppConfig struct {
	Port string `yaml:"port"`
}

// DatabasesConfig struct
type DatabasesConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
	MySQL    MySQLConfig    `yaml:"mysql"`
	Redis    RedisConfig    `yaml:"redis"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Timeout  int    `yaml:"timeout"`
	MaxIdle  int    `yaml:"max_idle"`
	MaxAct   int    `yaml:"max_act"`
}

// PostgresConfig struct
type PostgresConfig struct {
	Master string `yaml:"master"`
	Slave  string `yaml:"slave"`
	MaxCon int    `yaml:"max_con"`
	Retry  int    `yaml:"retry"`
}

// MySQLConfig struct
type MySQLConfig struct {
	Master string `yaml:"master"`
	Slave  string `yaml:"slave"`
}

// ElasticConfig struct
type ElasticConfig struct {
	Host string `yaml:"host"`
}
