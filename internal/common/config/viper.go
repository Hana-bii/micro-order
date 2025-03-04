package config

import (
	"strings"

	"github.com/spf13/viper"
)

func NewViperConfig() error {
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../common/config")
	// 将global.yaml配置里的key中的'-'替换为'_'
	viper.EnvKeyReplacer(strings.NewReplacer("_", "-"))
	// _ = viper.BindEnv("stripe-key", "STRIPE_KEY")
	viper.AutomaticEnv()
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET")
	return viper.ReadInConfig()
}
