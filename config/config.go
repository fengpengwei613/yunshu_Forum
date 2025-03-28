package config

import (
	"fmt"
	"path"

	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 项目配置，监听地址和运行模式
type ServerConfig struct {
	Addr string `mapstructure:"addr"`
	Mode string `mapstructure:"mode"`
}

// 主从数据库配置
type DatabaseConfig struct {
	Master struct {
		Dsn string `mapstructure:"dsn"`
	}
	Slave struct {
		Dsn string `mapstructure:"dsn"`
	}
	Maxidleconn int `mapstructure:"max_idle_conns"`
	Maxopenconn int `mapstructure:"max_open_conns"`
	Maxidletime int `mapstructure:"max_idle_time"`
	Maxlifetime int `mapstructure:"max_life_time"`
}

// Redis配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Poolsize int    `mapstructure:"pool_size"`
}

// jwt配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

type OssConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access_key"`
	AccessKeySecret string `mapstructure:"access_secret"`
	Bucket          string `mapstructure:"bucket"`
	Url             string `mapstructure:"cdn_domain"`
}
type EmailConfig struct {
	Privider   string `mapstructure:"privider"`
	Host       string `mapstructure:"smtp_host"`
	Port       int    `mapstructure:"smtp_port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	From       string `mapstructure:"from_name"`
	Enable_ssl bool   `mapstructure:"enable_ssl"`
}

var (
	Server   = &ServerConfig{}
	Database = &DatabaseConfig{}
	Redis    = &RedisConfig{}
	JWT      = &JWTConfig{}
	Oss      = &OssConfig{}
	Email    = &EmailConfig{}
)

// 初始化配置
func InitConfig() {
	//获取项目目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("获取项目目录失败")
	}
	root := path.Dir(path.Dir(filename))
	//拼接配置文件路径
	configPath := path.Join(root, "/config")

	setDefaults()
	viper.SetConfigName(".config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("配置文件读取失败，原因：%w", err))
	}
	//检测配置文件更新
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		reloadConfig()
	})
	reloadConfig()
}

func setDefaults() {
	viper.SetDefault("server.addr", ":8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.database.master.maxidleconn", 10)
	viper.SetDefault("server.database.master.maxopenconn", 100)
	viper.SetDefault("server.database.slave.maxidleconn", 10)
	viper.SetDefault("server.database.slave.maxopenconn", 100)
}

func reloadConfig() {
	// 解析全部配置项到结构体
	if err := viper.UnmarshalKey("server", &Server); err != nil {
		panic(fmt.Errorf("服务配置错误: %w", err))
	}
	if err := viper.UnmarshalKey("database", &Database); err != nil {
		panic(fmt.Errorf("数据库配置错误: %w", err))
	}
	if err := viper.UnmarshalKey("redis", &Redis); err != nil {
		panic(fmt.Errorf("Redis配置错误: %w", err))
	}
	if err := viper.UnmarshalKey("jwt", &JWT); err != nil {
		panic(fmt.Errorf("Jwt配置错误: %w", err))
	}
	if err := viper.UnmarshalKey("oss", &Oss); err != nil {
		panic(fmt.Errorf("Oss配置错误: %w", err))
	}
	if err := viper.UnmarshalKey("email", &Email); err != nil {
		panic(fmt.Errorf("Email配置错误: %w", err))
	}

	// 强制校验必要配置项
	if JWT.Secret == "" {
		panic("jwt密钥不能为空")
	}
	if Database.Master.Dsn == "" {
		panic("数据库主库dsn不能为空")
	}
	if Database.Slave.Dsn == "" {
		panic("数据库从库dsn不能为空")
	}
	if Redis.Addr == "" {
		panic("Redis地址不能为空")
	}
	if Oss.Endpoint == "" {
		panic("Oss地址不能为空")
	}
	if Email.Host == "" {
		panic("Email地址不能为空")
	}
}
