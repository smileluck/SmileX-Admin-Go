// Package conf 配置加载（viper），对应 Kratos 的 internal/conf。
package conf

import (
	"github.com/spf13/viper"
)

// Bootstrap 应用根配置
type Bootstrap struct {
	Server Server `mapstructure:"server"`
	DB     DB     `mapstructure:"db"`
	JWT    JWT    `mapstructure:"jwt"`
}

type Server struct {
	Port      int    `mapstructure:"port"`
	Mode      string `mapstructure:"mode"`
	StaticDir string `mapstructure:"staticDir"` // 前端产物目录（如 web/dist），存在则由后端托管 SPA
}

type DB struct {
	Driver       string   `mapstructure:"driver"`
	AutoMigrate  bool     `mapstructure:"autoMigrate"`
	MySQL        MySQL    `mapstructure:"mysql"`
	Postgres     Postgres `mapstructure:"postgres"`
	SQLite       SQLite   `mapstructure:"sqlite"`
}

type MySQL struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DBName      string `mapstructure:"dbname"`
	Charset     string `mapstructure:"charset"`
	MaxOpenConns int  `mapstructure:"maxOpenConns"`
	MaxIdleConns int  `mapstructure:"maxIdleConns"`
}

type Postgres struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DBName      string `mapstructure:"dbname"`
	SSLMode     string `mapstructure:"sslmode"`
	MaxOpenConns int  `mapstructure:"maxOpenConns"`
	MaxIdleConns int  `mapstructure:"maxIdleConns"`
}

type SQLite struct {
	Path string `mapstructure:"path"`
}

type JWT struct {
	Secret        string `mapstructure:"secret"`
	Issuer        string `mapstructure:"issuer"`
	ExpireHours   int    `mapstructure:"expireHours"`
	RefreshHours  int    `mapstructure:"refreshHours"`
}

// Load 从指定路径加载配置
func Load(path string) (*Bootstrap, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var c Bootstrap
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

// LoadDefault 从默认路径列表加载（项目根 configs/config.yaml 或 -conf 指定路径）
func LoadDefault() (*Bootstrap, error) {
	paths := []string{"configs/config.yaml", "conf/config.yaml", "config.yaml"}
	for _, p := range paths {
		if c, err := Load(p); err == nil {
			return c, nil
		}
	}
	return nil, viper.ConfigFileNotFoundError{}
}
