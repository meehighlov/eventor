package config

import (
	"log"
	"os"

	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	PROD  = "prod"
	LOCAL = "local"
)

type Config struct {
	ENV                        string `env:"ENV" env-default:"local"`
	Users                      string `env:"USERS" env-required:"true"`
	ReportChatId               string `env:"REPORT_CHAT_ID" env-required:"true"`
	Timezone                   string `env:"TIMEZONE" env-default:"Europe/Moscow"`
	PostgresDSN                string `env:"POSTGRES_DSN" env-required:"true"`
	MigrationsDir              string `env:"MIGRATIONS_DIR" env-default:"eventor-migrations"`
	RunMigrations              bool   `env:"RUN_MIGRATIONS" env-default:"true"`
	WatcherCheckIntervalSec    int    `env:"WATCHER_CHECK_INTERVAL_SEC" env-default:"10"`
	ChatCacheExpirationMinutes int    `env:"CHAT_CACHE_EXPIRATION_MINUTES" env-default:"10"`
	LoggingFileName            string `env:"LOGGING_FILE_NAME" env-default:"eventor.log"`

	TelegramToken              string `env:"TELEGRAM_TOKEN" env-required:"true"`
	TelegramUseWebook          bool   `env:"TELEGRAM_USE_WEBHOOK" env-default:"false"`
	TelegramWebhookToken       string `env:"TELEGRAM_WEBHOOK_TOKEN" env-default:""`
	TelegramWebhookAddress     string `env:"TELEGRAM_WEBHOOK_ADDRESS" env-default:":8080"`
	TelegramWebhookTLSAddress  string `env:"TELEGRAM_WEBHOOK_TLS_ADDRESS" env-default:":443"`
	TelegramWebhookTLSCertFile string `env:"TELEGRAM_WEBHOOK_TLS_CERT_FILE" env-default:""`
	TelegramWebhookTLSKeyFile  string `env:"TELEGRAM_WEBHOOK_TLS_KEY_FILE" env-default:""`
	TelegramUseTLS             bool   `env:"TELEGRAM_USE_TLS" env-default:"false"`
	TelegramHandlerTimeoutSec  int    `env:"TELEGRAM_HANDLER_TIMEOUT_SEC" env-default:"2"`

	RedisAddr     string `env:"REDIS_ADDR" env-default:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" env-default:""`
	RedisDB       int    `env:"REDIS_DB" env-default:"0"`
}

func (cfg *Config) AllowedUsers() []string {
	return strings.Split(cfg.Users, ",")
}

// loads config from .env
// panics on any read error
// also sets TZ env variable from according .env value
func MustLoad() *Config {
	if _, err := os.Stat("env/eventor/.env"); os.IsNotExist(err) {
		log.Fatal("Not found .env file")
	}

	cfg := Config{}
	err := cleanenv.ReadConfig("env/eventor/.env", &cfg)
	if err != nil {
		log.Fatal("Failed to read envs:", err.Error())
	}

	os.Setenv("TZ", cfg.Timezone)

	return &cfg
}
