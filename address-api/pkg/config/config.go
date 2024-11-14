package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	defaultConfigPath     = "pkg/config/"
	tagName               = "mapstructure"
	configFileType        = "yaml"
	defaultConfigFileName = "config-dev"
	prodConfigFileName    = "config-prod"
	environmentKey        = "environment"
)

type Config struct {
	Server   ServerConfig         `mapstructure:"server"`
	Postgres PostgresConfig       `mapstructure:"postgres"`
	Logger   LoggerConfig         `mapstructure:"logger"`
	Auth     AuthenticationConfig `mapstructure:"auth"`
	Metric   MetricConfig         `mapstructure:"metric"`
}

type ServerConfig struct {
	AppVersion     string        `mapstructure:"appVersion"`
	Host           string        `mapstructure:"host"`
	Port           string        `mapstructure:"port"`
	Mode           string        `mapstructure:"mode"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
	SSL            bool          `mapstructure:"ssl"`
	MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
	CtxTimeout     time.Duration `mapstructure:"ctxTimeout"`
}

type LoggerConfig struct {
	Development bool   `mapstructure:"development"`
	Encoding    string `mapstructure:"encoding"`
	Level       string `mapstructure:"level"`
}

type PostgresConfig struct {
	Host               string `mapstructure:"host"`
	Port               string `mapstructure:"port"`
	UserName           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	DbName             string `mapstructure:"dbname"`
	SSLMode            bool   `mapstructure:"ssl-mode"`
	Driver             string `mapstructure:"driver"`
	MaxOpenConnections int    `mapstructure:"maxOpenConnections"`
	ConnMaxLifeTime    int    `mapstructure:"connMaxLifetime"`
	MaxIdleConnections int    `mapstructure:"maxIdleConnections"`
	ConnMaxIdleTime    int    `mapstructure:"connMaxIdleTime"`
}

type AuthenticationConfig struct {
	JwtSecret string `mapstructure:"jwtSecret"`
}

type MetricConfig struct {
	Url         string `mapstructure:"url"`
	ServiceName string `mapstructure:"serviceName"`
}

func NewConfig() *Config {
	env, _ := os.LookupEnv(environmentKey)
	fmt.Println("Environment: [" + env + "] was successfully read from runtime arguments [" + environmentKey + "].")

	return ReadConfig(&Config{}, strings.ToUpper(env))
}

func addKeysToViper(v *viper.Viper) {
	var reply interface{} = Config{}
	t := reflect.TypeOf(reply)
	keys := getAllKeys(t)
	for _, key := range keys {
		v.SetDefault(key, "")
	}
}

func getAllKeys(t reflect.Type) []string {
	var result []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		n := strings.ToUpper(f.Tag.Get(tagName))
		if reflect.Struct == f.Type.Kind() {
			subKeys := getAllKeys(f.Type)
			for _, k := range subKeys {
				result = append(result, n+"."+k)
			}
		} else {
			result = append(result, n)
		}
	}
	return result
}

var readFromEnv = func(v *viper.Viper) *viper.Viper {
	fmt.Println("Reading environment configuration")
	addKeysToViper(v)
	v.AutomaticEnv()
	return v
}

var readFromAppYml = func(v *viper.Viper) *viper.Viper {
	fmt.Println("Reading application yml configuration")
	v.SetConfigName(defaultConfigFileName)
	v.SetTypeByDefaultValue(true)
	v.SetConfigType(configFileType)
	v.AddConfigPath("./" + defaultConfigPath)
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Viper read config has an error; %e\n", err)
	}

	return v
}

var readFromAppProdYml = func(v *viper.Viper) *viper.Viper {
	fmt.Println("Reading application yml configuration")
	v.SetConfigName(prodConfigFileName)
	v.SetTypeByDefaultValue(true)
	v.SetConfigType(configFileType)
	v.AddConfigPath("./" + defaultConfigPath)
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Viper read config has an error; %e\n", err)
	}

	return v
}

var ReadConfig = func(c *Config, env string) *Config {
	fmt.Println("Configuration read initiated...")
	v := viper.New()
	switch {
	case env == "DEV":
		v = readFromAppYml(v)
	case env == "PROD":
		v = readFromAppProdYml(v)
	default:
		v = readFromEnv(v)
	}
	if err := v.Unmarshal(&c); err != nil {
		panic(any("Configuration unmarshalling occurred an error, terminating: " + err.Error()))
	}

	return c
}
