package config

import "os"

type Config struct {
    TraefikNetwork   string
    Domain           string
    LetsEncryptEmail string
    AuthCredentials  string // Basic auth in htpasswd format
    DefaultVNCConfig VNCConfig
    BehindProxy      bool
    // Container environment variables
    ContainerEnvVars EnvVars
}

type VNCConfig struct {
    Password   string `json:"password"`
    Resolution string `json:"resolution"`
    ColDepth   int    `json:"colDepth"`
    ViewOnly   bool   `json:"viewOnly"`
    Display    string `json:"display"`
}

// EnvVars holds environment variables for containers
type EnvVars struct {
    Width                        string
    Height                       string
    PlaywrightChromiumPath       string
    LogLevel                     string
    RabbitMQQueue                string
    OpenAIAPIKey                 string
    AnthropicAPIKey              string
    Port                         string
    RabbitMQUser                 string
    RabbitMQPassword             string
    RabbitMQHost                 string
    RabbitMQPort                 string
}

func Load() *Config {
    return &Config{
        TraefikNetwork:   "traefik_network",
        Domain:           "yourdomain.com",
        LetsEncryptEmail: "admin@yourdomain.com",
        AuthCredentials:  "user:$apr1$...", // Generated htpasswd
        BehindProxy:      os.Getenv("BEHIND_PROXY") == "true",
        DefaultVNCConfig: VNCConfig{
            Password:   "headless",
            Resolution: "1360x768",
            ColDepth:   24,
            ViewOnly:   false,
            Display:    ":1",
        },
        ContainerEnvVars: EnvVars{
            Width:                  getEnvWithDefault("WIDTH", "1024"),
            Height:                 getEnvWithDefault("HEIGHT", "768"),
            PlaywrightChromiumPath: getEnvWithDefault("PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH", "/usr/bin/chromium"),
            LogLevel:               getEnvWithDefault("LOG_LEVEL", "debug"),
            RabbitMQQueue:          getEnvWithDefault("RABBITMQ_QUEUE", "browser_tasks"),
            OpenAIAPIKey:           os.Getenv("OPENAI_API_KEY"),
            AnthropicAPIKey:        os.Getenv("ANTHROPIC_API_KEY"),
            Port:                   getEnvWithDefault("PORT", "3000"),
            RabbitMQUser:           getEnvWithDefault("RABBITMQ_USER", "admin"),
            RabbitMQPassword:       getEnvWithDefault("RABBITMQ_PASSWORD", "admin"),
            RabbitMQHost:           getEnvWithDefault("RABBITMQ_HOST", "rabbitmq"),
            RabbitMQPort:           getEnvWithDefault("RABBITMQ_PORT", "5672"),
        },
    }
}

// Helper function to get environment variable with default value
func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}