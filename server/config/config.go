package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config struct stores entire project configurations
type Config struct {
	ServerConfig        ServerConfig        `mapstructure:"server"`
	APIConfig           APIConfig           `mapstructure:"api"`
	APPConfig           APPConfig           `mapstructure:"app"`
	KafkaConfig         KafkaConfig         `mapstructure:"kafka"`
	LoggerConfig        LoggerConfig        `mapstructure:"logger"`
	DatabaseConfig      DatabaseConfig      `mapstructure:"database"`
	RedisConfig         RedisConfig         `mapstructure:"redis"`
	MiddlewareConfig    MiddlewareConfig    `mapstructure:"middleware"`
	TokenAuthConfig     TokenAuthConfig     `mapstructure:"token"`
	SNSConfig           SNSConfig           `mapstructure:"sns"`
	SESConfig           SESConfig           `mapstructure:"ses"`
	ElasticsearchConfig ElasticsearchConfig `mapstructure:"elasticsearch"`
	GoogleOAuth         GoogleOAuth         `mapstructure:"googleOAuth"`
	SessionConfig       SessionConfig       `mapstructure:"session"`
	SentryConfig        SentryConfig        `mapstructure:"sentry"`
	KaleyraConfig       KaleyraConfig       `mapstructure:"kaleyra"`
	MSGPlatformConfig   MSGPlatformConfig   `mapstructure:"message_platform"`
	GoKwikConfig        GoKwikConfig        `mapstructure:"goKwik"`
	// CommissionOrderListeneronfig ListenerConfig      `mapstructure:"commissionOrderListenerConsumer"`
}

// ServerConfig has only server specific configuration
type ServerConfig struct {
	ListenAddr     string        `mapstructure:"listenAddr"`
	Port           string        `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
	CloseTimeout   time.Duration `mapstructure:"closeTimeout"`
	Env            string        `mapstructure:"env"`
	UseMemoryStore bool          `mapstructure:"useMemoryStore"`
	CORSConfig     CORSConfig    `mapstructure:"cors"`
	HypdApiConfig  HypdApiConfig `mapstructure:"hypdApi"`
}

type GoogleOAuth struct {
	ClientID     string   `mapstructure:"clientID"`
	ClientSecret string   `mapstructure:"clientSecret"`
	RedirectURL  string   `mapstructure:"redirectURL"`
	Scopes       []string `mapstructure:"scopes"`
	State        string   `mapstructure:"state"`
}

// SentryConfig contains sentry related configuration
type SentryConfig struct {
	EnableSentry bool   `mapstructure:"enable"`
	DSN          string `mapstructure:"dsn"`
}

// APIConfig contains api package related configurations
type APIConfig struct {
	Mode                   string `mapstructure:"mode"`
	EnableTestRoute        bool   `mapstructure:"enableTestRoute"`
	EnableMediaRoute       bool   `mapstructure:"enableMediaRoute"`
	EnableStaticRoute      bool   `mapstructure:"enableStaticRoute"`
	MaxRequestDataSize     int    `mapstructure:"maxRequestDataSize"`
	KeeperLoginRedirectURL string `mapstructure:"keeperLoginRedirectURL"`
	HypdApiConfig          HypdApiConfig
	GoKwikConfig           GoKwikConfig
}

type SessionConfig struct {
	CookieConfig CookieConfig `mapstructure:"cookie"`
	RedisConfig  RedisConfig
}

type CookieConfig struct {
	Name     string `mapstructure:"name"`
	Path     string `mapstructure:"path"`
	HttpOnly bool   `mapstructure:"httpOnly"`
	Domain   string `mapstructure:"domain"`
	Secure   bool   `mapstructure:"secure"`
}

// APPConfig contains api package related configurations
type APPConfig struct {
	DatabaseConfig      DatabaseConfig
	TokenAuthConfig     TokenAuthConfig
	SNSConfig           SNSConfig
	SESConfig           SESConfig
	HypdApiConfig       HypdApiConfig
	GoKwikConfig        GoKwikConfig
	ElasticsearchConfig ElasticsearchConfig
	GoogleOAuth         GoogleOAuth
	Kaleyra             KaleyraConfig
	MSGPlatform         MSGPlatformConfig

	UserConfig              ServiceConfig `mapstructure:"user"`
	CustomerConfig          ServiceConfig `mapstructure:"customer"`
	BrandConfig             ServiceConfig `mapstructure:"brand"`
	InfluencerConfig        ServiceConfig `mapstructure:"influencer"`
	CartConfig              ServiceConfig `mapstructure:"cart"`
	ExpressCheckoutConfig   ServiceConfig `mapstructure:"expressCheckout"`
	WishlistConfig          ServiceConfig `mapstructure:"wishlist"`
	SizeProfileConfig       ServiceConfig `mapstructure:"sizeProfile"`
	CommissionInvoiceConfig ServiceConfig `mapstructure:"commissionInvoice"`

	CustomerChangeConfig          ListenerConfig `mapstructure:"customerChangeConsumer"`
	DiscountChangeConfig          ListenerConfig `mapstructure:"discountChangeConsumer"`
	CatalogChangeConfig           ListenerConfig `mapstructure:"catalogChangeConsumer"`
	InventoryChangeConfig         ListenerConfig `mapstructure:"inventoryChangeConsumer"`
	BrandChangeConfig             ListenerConfig `mapstructure:"brandChangeConsumer"`
	InfluencerChangeConfig        ListenerConfig `mapstructure:"influencerChangeConsumer"`
	CouponChangeConfig            ListenerConfig `mapstructure:"couponChangeConsumer"`
	CommissionOrderListenerConfig ListenerConfig `mapstructure:"commissionOrderListenerConsumer"`

	BrandFullProduceConfig       ProducerConfig `mapstructure:"brandFullProducer"`
	InfluencerFullProducerConfig ProducerConfig `mapstructure:"influencerFullProducer"`
}

// ElasticsearchConfig contains elasticsearch related configurations
type ElasticsearchConfig struct {
	Endpoint                  string `mapstructure:"endpoint"`
	Username                  string `mapstructure:"username"`
	Password                  string `mapstructure:"password"`
	BrandFullIndex            string `mapstructure:"brandFullIndex"`
	InfluencerFullIndex       string `mapstructure:"influencerFullIndex"`
	InfluencerCollectionIndex string `mapstructure:"influencerCollectionIndex"`
	InfluencerProductIndex    string `mapstructure:"influencerProductIndex"`
	ContentFullIndex          string `mapstructure:"contentFullIndex"`
}

//HypdApiConfig contains config related to other services
type HypdApiConfig struct {
	CatalogApi string `mapstructure:"catalogApi"`
	OrderApi   string `mapstructure:"orderApi"`
	CouponApi  string `mapstructure:"couponApi"`
	Token      string `mapstructure:"token"`
	CatalogURL string `mapstructure:"catalogURL"`
}

// GoKwikConfig has GoKwik configurations
type GoKwikConfig struct {
	RTOApi     string `mapstructure:"rtoApi"`
	AppID      string `mapstructure:"appID"`
	MerchantID string `mapstructure:"merchantID"`
	AppSecret  string `mapstructure:"appSecret"`
}

// ServiceConfig contains app service related config
type ServiceConfig struct {
	DBName string `mapstructure:"dbName"`
}

// ListenerConfig contains app kafka topic listener related config
type ListenerConfig struct {
	GroupID  string   `mapstructure:"groupId"`
	Brokers  []string `mapstructure:"brokers"`
	Topic    string   `mapstructure:"topic"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

// ProducerConfig contains app kafka topic producer related config
type ProducerConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	Topic    string   `mapstructure:"topic"`
	Async    bool     `mapstructure:"async"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

// TokenAuthConfig contains token authentication related configuration
type TokenAuthConfig struct {
	OTPLength        int    `mapstructure:"otpLength"`
	HashPasswordCost int    `mapstructure:"hashPasswordCost"`
	JWTSignKey       string `mapstructure:"jwtSignKey"`
	JWTExpiresAt     int64  `mapstructure:"expiresAt"`
	AuthServerUrl    string `mapstructure:"authServerUrl"`
}

// CORSConfig contains cors related config
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowedOrigins"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
}

// KafkaConfig has kafka cluster specific configuration
type KafkaConfig struct {
	EnableKafka bool     `mapstructure:"enableKafka"`
	BrokerDial  string   `mapstructure:"brokerDial"`
	BrokerURL   string   `mapstructure:"brokerUrl"`
	BrokerPort  string   `mapstructure:"brokerPort"`
	Brokers     []string `mapstructure:"brokers"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
}

// LoggerConfig contains different logger configurations
type LoggerConfig struct {
	KafkaLoggerConfig   `mapstructure:"kafkaLog"`
	FileLoggerConfig    `mapstructure:"fileLog"`
	ConsoleLoggerConfig `mapstructure:"consoleLog"`
}

// KafkaLoggerConfig contains kafka logger specific configuration
type KafkaLoggerConfig struct {
	EnableKafkaLogger bool   `mapstructure:"enableKafkaLog"`
	KafkaTopic        string `mapstructure:"kafkaTopic"`
	KafkaPartition    string `mapstructure:"kafkaPartition"`
}

// ConsoleLoggerConfig contains file console logging specific configuration
type ConsoleLoggerConfig struct {
	EnableConsoleLogger bool `mapstructure:"enableConsoleLog"`
}

// FileLoggerConfig contains file logging specific configuration
type FileLoggerConfig struct {
	FileName         string `mapstructure:"fileName"`
	Path             string `mapstructure:"path"`
	EnableFileLogger bool   `mapstructure:"enableFileLog"`
	MaxBackupsFile   int    `mapstructure:"maxBackupFile"`
	MaxSize          int    `mapstructure:"maxFileSize"`
	MaxAge           int    `mapstructure:"maxAge"`
	Compress         bool   `mapstructure:"compress"`
}

// DatabaseConfig contains mongodb related configuration
type DatabaseConfig struct {
	Scheme string `mapstructure:"scheme"`
	Host   string `mapstructure:"host"`
	// Name     string `mapstructure:"name"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	ReplicaSet string `mapstructure:"replicaSet"`
}

// SNSConfig contains aws sns service related configuration
type SNSConfig struct {
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"accessKeyId"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
}

// KaleyraConfig contains aws kaleyra service related configuration
type KaleyraConfig struct {
	Name       string `mapstructure:"name"`
	APIKey     string `mapstructure:"apiKey"`
	SID        string `mapstructure:"sid"`
	TemplateID string `mapstructure:"templateID"`
}

// MSGPlatformConfig contains aws Message service related configuration
type MSGPlatformConfig struct {
	Name string `mapstructure:"name"`
}

// SESConfig contains aws ses service related configuration
type SESConfig struct {
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"accessKeyId"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
}

// ConnectionURL returns connection string to of mongodb storage
func (d *DatabaseConfig) ConnectionURL() string {
	url := fmt.Sprintf("%s://", d.Scheme)
	if d.Username != "" && d.Password != "" {
		url += fmt.Sprintf("%s:%s@", d.Username, d.Password)
	}
	url += fmt.Sprintf("%s", d.Host)
	if d.ReplicaSet != "" {
		url += fmt.Sprintf("/?replicaSet=%s", d.ReplicaSet)
	}
	return url
}

// RedisConfig has cache related configuration.
type RedisConfig struct {
	Network  string `mapstructure:"network"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// ConnectionURL returns connection string to of mongodb storage
func (r *RedisConfig) ConnectionURL() string {
	var url string
	if r.Username != "" {
		url += fmt.Sprintf("%s", r.Username)
	}
	if r.Password != "" {
		url += fmt.Sprintf(":%s@", r.Password)
	}
	url += fmt.Sprintf("%s", r.Host)
	if r.Port != "" {
		url += fmt.Sprintf(":%s", r.Port)
	}
	return url
}

// MiddlewareConfig has middlewares related configuration
type MiddlewareConfig struct {
	EnableRequestLog bool `mapstructure:"enableRequestLog"`
}

// GetConfig returns entire project configuration
func GetConfig() *Config {
	return GetConfigFromFile("default")
}

// GetConfigFromFile returns configuration from specific file object
func GetConfigFromFile(fileName string) *Config {
	if fileName == "" {
		fileName = "default"
	}

	// looking for filename `default` inside `src/server` dir with `.toml` extension
	viper.SetConfigName(fileName)
	viper.AddConfigPath("../conf/")
	viper.AddConfigPath("../../conf/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf/")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("couldn't load config: %s", err)
		os.Exit(1)
	}
	config := &Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("couldn't read config: %s", err)
		os.Exit(1)
	}
	config.APPConfig.TokenAuthConfig = config.TokenAuthConfig
	config.APPConfig.SNSConfig = config.SNSConfig
	config.APPConfig.SESConfig = config.SESConfig
	config.APIConfig.HypdApiConfig = config.ServerConfig.HypdApiConfig
	config.APPConfig.HypdApiConfig = config.ServerConfig.HypdApiConfig
	config.APPConfig.ElasticsearchConfig = config.ElasticsearchConfig
	config.APPConfig.GoogleOAuth = config.GoogleOAuth
	config.SessionConfig.RedisConfig = config.RedisConfig
	config.APPConfig.Kaleyra = config.KaleyraConfig
	config.APPConfig.MSGPlatform = config.MSGPlatformConfig
	config.APIConfig.GoKwikConfig = config.GoKwikConfig
	config.APPConfig.GoKwikConfig = config.GoKwikConfig

	return config
}
