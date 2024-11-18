package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
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
	Logger LoggerConfig `mapstructure:"logger"`
	Kafka  KafkaConfig  `mapstructure:"kafka"`
}

type LoggerConfig struct {
	Development      bool   `mapstructure:"development"`
	Encoding         string `mapstructure:"encoding"`
	Level            string `mapstructure:"level"`
	IndexName        string `mapstructure:"indexName"`
	ElasticsearchUrl string `mapstructure:"elasticSearchUrl"`
}

type KafkaConfig struct {
	Brokers        []string `mapstructure:"brokers"`
	ConsumerGroup  string   `mapstructure:"consumerGroup"`
	MaxPollRecords int      `mapstructure:"maxPollRecords"`
	GroupID        string   `mapstructure:"groupId"`
	AutoCommit     bool     `mapstructure:"autoCommit"`
	FetchMaxWaitMs int      `mapstructure:"fetchMaxWaitMs"`
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
