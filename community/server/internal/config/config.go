package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Mysqldb     Mysqldb             `yaml:"mysql"`
	Redis       Redis               `yaml:"redis"`
	Logger      Logger              `yaml:"logger"`
	Server      Server              `yaml:"server"`
	XxlJob      XxlJobConfig        `yaml:"xxlJob"`
	File        File                `yaml:"file"`
	AI          AIConfig            `yaml:"ai"`
	SpotFilters map[string][]string `yaml:"spot_filters"`
}

type Server struct {
	Prot string `yaml:"port"`
}

type Mysqldb struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
	PoolSize int    `yaml:"pool_size"`
}

type Logger struct {
	FileName   string `yaml:"fileName"`
	Path       string `yaml:"path"`
	MaxAge     int    `yaml:"maxAge"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
}

type XxlJobConfig struct {
	IsEnable    bool          `yaml:"isEnable"`
	ServerAddrs []string      `yaml:"serverAddrs"`
	AccessToken string        `yaml:"accessToken"`
	AppName     string        `yaml:"appName"`
	ClientPort  int           `yaml:"clientPort"`
	Timeout     time.Duration `yaml:"timeout"`
	BeatTime    time.Duration `yaml:"beatTime"`
	LogLevel    int           `yaml:"logLevel"`
}

type File struct {
	Path         string `yaml:"path"`
	ExternalPath string `yaml:"externalPath"`
}

type AIConfig struct {
	Provider string `yaml:"provider"`
	ApiKey   string `yaml:"api_key"`
	Url      string `yaml:"url"`
	Model    string `yaml:"model"`
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
		zap.S().Errorf("Error reading config file: %v", zap.Error(err))
		return nil, err
	}

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

	spotFilters := make(map[string][]string)
	err = viper.UnmarshalKey("spot_filters", &spotFilters)
	if err != nil {
		zap.L().Error("unmarshal spot_filters error: ", zap.Error(err))
	}

	mc := &Config{}
	err = viper.Unmarshal(&mc)
	if err != nil {
		zap.L().Error("unmarshal error: ", zap.Error(err))
		return nil, err
	}

	mc.SpotFilters = spotFilters

	mc.AI = AIConfig{
		ApiKey: viper.GetString("AI_API_KEY"),
		Url:    viper.GetString("AI_BASE_URL"),
		Model:  viper.GetString("AI_MODEL"),
	}

	return mc, nil
}

func replaceEnvVariables(value string) string {
	return os.Expand(value, func(key string) string {
		return os.Getenv(key)
	})
}

func trimSpace(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, " ")
	return s
}
