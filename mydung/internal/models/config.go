package models

// Config pure system ki settings ko wrap karta hai
type Config struct {
	App    AppConfig    `mapstructure:"app" yaml:"app"`
	Worker WorkerConfig `mapstructure:"worker" yaml:"worker"`
	Redis  RedisConfig  `mapstructure:"redis" yaml:"redis"`
}

type AppConfig struct {
	Port int `mapstructure:"port" yaml:"port"`
}

type WorkerConfig struct {
	Count     int `mapstructure:"count" yaml:"count"`
	QueueSize int `mapstructure:"queue_size" yaml:"queue_size"`
}

type RedisConfig struct {
    Host        string `mapstructure:"host" yaml:"host"`
    Port        int    `mapstructure:"port" yaml:"port"`
    Password    string `mapstructure:"password" yaml:"password"`
    DB          int    `mapstructure:"db" yaml:"db"`
    MaxRetries  int    `mapstructure:"max_retries" yaml:"max_retries"`
    RetryDelay  int    `mapstructure:"retry_delay" yaml:"retry_delay"` // Seconds mein
}