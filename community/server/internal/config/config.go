package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Mysqldb Mysqldb  `yaml:"mysql" mapstructure:"mysql"`
	Redis   Redis    `yaml:"redis" mapstructure:"redis"`
	Logger  Logger   `yaml:"logger" mapstructure:"logger"`
	Server  Server   `yaml:"server" mapstructure:"server"`
	File    File     `yaml:"file" mapstructure:"file"`
	AI      AIConfig `yaml:"ai" mapstructure:"ai"`
	IM      IMConfig `yaml:"im" mapstructure:"im"`
}

type Server struct {
	Port string `yaml:"port" mapstructure:"port"`
}

type Mysqldb struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Dbname   string `yaml:"dbname" mapstructure:"dbname"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
}

type Redis struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	DBname   string `yaml:"dbname" mapstructure:"dbname"`
	PoolSize int    `yaml:"pool_size" mapstructure:"pool_size"`
}

type Logger struct {
	FileName   string `yaml:"fileName" mapstructure:"fileName"`
	Path       string `yaml:"path" mapstructure:"path"`
	MaxAge     int    `yaml:"maxAge" mapstructure:"maxAge"`
	MaxSize    int    `yaml:"maxSize" mapstructure:"maxSize"`
	MaxBackups int    `yaml:"maxBackups" mapstructure:"maxBackups"`
}

type File struct {
	Path         string `yaml:"path"`
	ExternalPath string `yaml:"externalPath"`
}

type IMConfig struct {
	BaseURL   string `yaml:"base_url"`
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
}

type AIConfig struct {
	ApiKey string `yaml:"api_key"`
	Url    string `yaml:"url"`
	Model  string `yaml:"model"`
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(dir)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	log.Printf("Config file loaded from: %s", viper.ConfigFileUsed())

	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		if value == "" {
			continue
		}
		replacedValue := replaceEnvVariables(value)
		replacedValue = trimSpace(replacedValue)
		replacedValue = strings.ReplaceAll(replacedValue, "'", "")
		replacedValue = strings.ReplaceAll(replacedValue, "\"", "")
		viper.Set(key, replacedValue)
	}

	mc := &Config{}
	err = viper.Unmarshal(&mc)
	if err != nil {
		log.Printf("unmarshal error: %v", err)
		return nil, err
	}

	log.Printf("MySQL Config: host=%s, port=%s, dbname=%s", mc.Mysqldb.Host, mc.Mysqldb.Port, mc.Mysqldb.Dbname)

	mc.AI = AIConfig{
		ApiKey: viper.GetString("AI_API_KEY"),
		Url:    viper.GetString("AI_BASE_URL"),
		Model:  viper.GetString("AI_MODEL"),
	}

	return mc, nil
}

func replaceEnvVariables(value string) string {
	return os.Expand(value, func(key string) string {
		if idx := strings.Index(key, ":-"); idx != -1 {
			env := os.Getenv(key[:idx])
			if env != "" {
				return env
			}
			return key[idx+2:]
		}
		return os.Getenv(key)
	})
}

func trimSpace(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, " ")
	return s
}
