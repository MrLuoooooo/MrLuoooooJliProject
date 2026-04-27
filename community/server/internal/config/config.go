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
	Mysqldb Mysqldb      `yaml:""mysql`
	Redis   Redis        `yaml:""redis`
	Logger  Logger       `yaml:""logger`
	Server  Server       `yaml:""server`
	XxlJob  XxlJobConfig `yaml:""xxlJob`
	File    File         `yaml:""file`
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

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}
func replaceEnvVariables(value string) string {
	return os.Expand(value, func(key string) string {
		return os.Getenv(key)
	})
}
func NewConfigFromPath(path ...string) (*Config, error) {
	pathStr := ""
	if isEmptyStringArray(path) {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		pathStr = dir
		viper.SetConfigName("config")
		viper.AddConfigPath(pathStr)
		zap.S.Infof("路径将会从 %s:%s 加载", pathStr, "config.yaml")
	} else {
		pathStr = path[0]
		viper.SetConfigFile(pathStr)
		zap.S.Infof("路径将会从 %s 加载", pathStr)
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		zap.S().Errorf("Error reading config file: %v", zap.Error(err))
		return nil, err
	}
	viper.AllKeys()
	for _, v := range viper.AllKeys() {
		value = viper.GetString(v)
		if value == "" {
			continue
		}
		replacedValue := replaceEnvVariables(value)
		replacedValue = trimSpace(replacedValue)
		replacedValue = strings.ReplaceAll(replacedValue, "'", "")
		replacedValue = strings.ReplaceAll(replacedValue, "\"", "")
		viper.Set(key, replacedValue)
}
potFilters := make(map[string][]string, 0)
	err = viper.UnmarshalKey("spot_filters", &spotFilters)
	if err != nil {
		zap.L().Error("unmarshal contract-role error: ", zap.Error(err))
	}

	mc := &Config{}
	err = viper.Unmarshal(&mc)
	if err != nil {
		zap.L().Error("unmarshal error: ", zap.Error(err))
		return nil, err
	}

	mc.SpotFilters = spotFilters
	return mc, nil
		
	}

}
func isEmptyStringArray(arr []string) bool {
	for len(arr) == 0 {
		return true
	}
	for _, v := range arr {
		if strings.TrimSpace(v) != "" {
			return false
		}
	}
	return true
}

func trimSpace(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, " ")
	return s
}
