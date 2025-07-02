package config

import (
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Redis     RedisConfig     `yaml:"redis"`
	JWT       JWTConfig       `yaml:"jwt"`
	WebSocket WebSocketConfig `yaml:"websocket"`
	Log       LogConfig       `yaml:"log"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	TimeZone string `yaml:"timezone"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	Secret      string        `yaml:"secret"`
	ExpireHours time.Duration `yaml:"expire_hours"`
}

type WebSocketConfig struct {
	ReadBufferSize  int  `yaml:"read_buffer_size"`
	WriteBufferSize int  `yaml:"write_buffer_size"`
	CheckOrigin     bool `yaml:"check_origin"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func Load() (*Config, error) {
	// 加载 .env 文件（如果存在）
	_ = godotenv.Load()

	// 先从 YAML 文件加载配置
	cfg := &Config{}
	if data, err := ioutil.ReadFile("configs/config.yaml"); err == nil {
		_ = yaml.Unmarshal(data, cfg)
	}

	// 环境变量覆盖配置文件值
	if port := getEnv("SERVER_PORT", ""); port != "" {
		cfg.Server.Port = port
	} else if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	if mode := getEnv("SERVER_MODE", ""); mode != "" {
		cfg.Server.Mode = mode
	} else if cfg.Server.Mode == "" {
		cfg.Server.Mode = "debug"
	}

	// 数据库配置
	if host := getEnv("DB_HOST", ""); host != "" {
		cfg.Database.Host = host
	} else if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}

	if port := getEnvAsInt("DB_PORT", 0); port != 0 {
		cfg.Database.Port = port
	} else if cfg.Database.Port == 0 {
		cfg.Database.Port = 3306
	}

	if user := getEnv("DB_USER", ""); user != "" {
		cfg.Database.User = user
	} else if cfg.Database.User == "" {
		cfg.Database.User = "root"
	}

	if password := getEnv("DB_PASSWORD", ""); password != "" {
		cfg.Database.Password = password
	}
	// 注意：这里不设置默认密码，保持YAML文件中的值

	if dbname := getEnv("DB_NAME", ""); dbname != "" {
		cfg.Database.DBName = dbname
	} else if cfg.Database.DBName == "" {
		cfg.Database.DBName = "gochat"
	}

	if timezone := getEnv("DB_TIMEZONE", ""); timezone != "" {
		cfg.Database.TimeZone = timezone
	} else if cfg.Database.TimeZone == "" {
		cfg.Database.TimeZone = "Asia/Shanghai"
	}

	// JWT 配置
	if secret := getEnv("JWT_SECRET", ""); secret != "" {
		cfg.JWT.Secret = secret
	} else if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "your-secret-key-change-in-production"
	}

	if hours := getEnvAsInt("JWT_EXPIRE_HOURS", 1); hours != 0 {
		cfg.JWT.ExpireHours = time.Duration(hours) * time.Hour
	} else if cfg.JWT.ExpireHours == 0 {
		cfg.JWT.ExpireHours = time.Duration(cfg.JWT.ExpireHours) * time.Hour
		if cfg.JWT.ExpireHours == 0 {
			cfg.JWT.ExpireHours = 24 * time.Hour
		}
	}

	// 日志配置
	if level := getEnv("LOG_LEVEL", ""); level != "" {
		cfg.Log.Level = level
	} else if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}

	if format := getEnv("LOG_FORMAT", ""); format != "" {
		cfg.Log.Format = format
	} else if cfg.Log.Format == "" {
		cfg.Log.Format = "json"
	}

	// WebSocket 配置设置默认值
	if cfg.WebSocket.ReadBufferSize == 0 {
		cfg.WebSocket.ReadBufferSize = 1024
	}
	if cfg.WebSocket.WriteBufferSize == 0 {
		cfg.WebSocket.WriteBufferSize = 1024
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
