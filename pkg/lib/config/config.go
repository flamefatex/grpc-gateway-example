package config

import (
	"strings"
	"time"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
	"github.com/spf13/viper"
)

var globalConfig *viper.Viper

func init() {
	globalConfig = viper.New()
}

// Provider defines a set of read-only methods for accessing the application
type Provider interface {
	ConfigFileUsed() string

	Get(key string) interface{}
	GetBool(key string) bool
	GetInt(key string) int
	GetIntSlice(key string) []int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetUint(key string) uint
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetString(key string) string
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetSizeInBytes(key string) uint
	UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error
	Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error
	UnmarshalExact(rawVal interface{}, opts ...viper.DecoderConfigOption) error
	IsSet(key string) bool
	InConfig(key string) bool
	AllKeys() []string
	AllSettings() map[string]interface{}
}

func Config() Provider {
	return globalConfig
}

// Init 初始化配置
func Init(serviceName string) error {
	v, err := readViperConfig(serviceName)
	if err != nil {
		return err
	}
	log.Infof("load config from file: %s", v.ConfigFileUsed())

	// 输出配置内容
	if v.GetBool("config.enableLog") {
		log.Infof("config all settings: %v", v.AllSettings())
	}

	globalConfig = v
	return nil
}

// readViperConfig 从约定的路径读取配置
func readViperConfig(serviceName string) (*viper.Viper, error) {
	v := viper.New()

	// env
	envServiceName := strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_"))
	v.SetEnvPrefix(envServiceName)
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	// config file
	lowerServiceName := strings.ToLower(serviceName)
	v.SetConfigName("config")                         // name of config file (without extension)
	v.AddConfigPath("/etc/" + lowerServiceName + "/") // path to look for the config file in
	v.AddConfigPath("$HOME/." + lowerServiceName)     // call multiple times to add many search paths
	v.AddConfigPath(".")                              // optionally look for config in the working directory
	err := v.ReadInConfig()                           // Find and read the config file
	if err != nil {                                   // Handle errors reading the config file
		return nil, err
	}
	return v, nil
}
