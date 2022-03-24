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
	ElasticsearchConfig ElasticsearchConfig `mapstructure:"elasticsearch"`
	SessionConfig       SessionConfig       `mapstructure:"session"`
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
	HypdApiConfig  HypdApiConfig `mapstructure:"hypdApi"`
	CORSConfig     CORSConfig    `mapstructure:"cors"`
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

// CORSConfig contains cors related config
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowedOrigins"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
}

// APIConfig contains api package related configurations
type APIConfig struct {
	Mode               string `mapstructure:"mode"`
	EnableTestRoute    bool   `mapstructure:"enableTestRoute"`
	EnableMediaRoute   bool   `mapstructure:"enableMediaRoute"`
	EnableStaticRoute  bool   `mapstructure:"enableStaticRoute"`
	MaxRequestDataSize int    `mapstructure:"maxRequestDataSize"`
	HypdApiConfig      HypdApiConfig
}

// APPConfig contains api package related configurations
type APPConfig struct {
	DatabaseConfig             DatabaseConfig
	HypdApiConfig              HypdApiConfig
	ElasticsearchConfig        ElasticsearchConfig
	KeeperCatalogConfig        ServiceConfig `mapstructure:"keeperCatalog"`
	CategoryConfig             ServiceConfig `mapstructure:"category"`
	BrandConfig                ServiceConfig `mapstructure:"brand"`
	DiscountConfig             ServiceConfig `mapstructure:"discount"`
	GroupConfig                ServiceConfig `mapstructure:"group"`
	CollectionConfig           ServiceConfig `mapstructure:"collection"`
	InfluencerCollectionConfig ServiceConfig `mapstructure:"influencer_collection"`
	InfluencerProductsConfig   ServiceConfig `mapstructure:"influencer_products"`
	ReviewConfig               ServiceConfig `mapstructure:"review"`
	InventoryConfig            ServiceConfig `mapstructure:"inventory"`
	PageSize                   int           `mapstructure:"pageSize"`

	CatalogChangeConfig              ListenerConfig `mapstructure:"catalogChangeConsumer"`
	CollectionCatalogChangeConfig    ListenerConfig `mapstructure:"collectionCatalogChangeConsumer"`
	CollectionChangeConfig           ListenerConfig `mapstructure:"collectionChangeConsumer"`
	DiscountChangeConfig             ListenerConfig `mapstructure:"discountChangeConsumer"`
	InventoryChangeConfig            ListenerConfig `mapstructure:"inventoryChangeConsumer"`
	BrandChangeConfig                ListenerConfig `mapstructure:"brandChangeConsumer"`
	ContentChangeConfig              ListenerConfig `mapstructure:"contentChangeConsumer"`
	GroupChangeConfig                ListenerConfig `mapstructure:"groupChangeConsumer"`
	ReviewChangeConfig               ListenerConfig `mapstructure:"reviewChangeConsumer"`
	InfluencerCollectionChangeConfig ListenerConfig `mapstructure:"influencerCollectionChangeConsumer"`
	InfluencerProductChangeConfig    ListenerConfig `mapstructure:"influencerProductChangeConfig"`

	CatalogFullProducerConfig          ProducerConfig `mapstructure:"catalogFullProducer"`
	CollectionFullProducerConfig       ProducerConfig `mapstructure:"collectionFullProducer"`
	InfluencerCollectionProducerConfig ProducerConfig `mapstructure:"influencerCollectionProducer"`
	InfluencerProductProducerConfig    ProducerConfig `mapstructure:"influencerProductProducer"`
	ReviewFullProducerConfig           ProducerConfig `mapstructure:"reviewFullProducer"`
}

// ElasticsearchConfig contains elasticsearch related configurations
type ElasticsearchConfig struct {
	Endpoint                  string `mapstructure:"endpoint"`
	Username                  string `mapstructure:"username"`
	Password                  string `mapstructure:"password"`
	CollectionFullIndex       string `mapstructure:"collectionFullIndex"`
	ReviewFullIndex           string `mapstructure:"reviewFullIndex"`
	CatalogFullIndex          string `mapstructure:"catalogFullIndex"`
	BrandFullIndex            string `mapstructure:"brandFullIndex"`
	InfluencerFullIndex       string `mapstructure:"influencerFullIndex"`
	ContentFullIndex          string `mapstructure:"contentFullIndex"`
	SeriesFullIndex           string `mapstructure:"seriesFullIndex"`
	InfluencerCollectionIndex string `mapstructure:"influencerCollectionIndex"`
	InfluencerProductIndex    string `mapstructure:"influencerProductIndex"`
}
type HypdApiConfig struct {
	CmsApi    string `mapstructure:"cmsApi"`
	EntityApi string `mapstructure:"entityApi"`
	Token     string `mapstructure:"token"`
}

// ServiceConfig contains app service related config
type ServiceConfig struct {
	DBName            string `mapstructure:"dbName"`
	CatalogContentURL string `mapstructure:"catalogContentUrl"`
}

// TokenAuthConfig contains token authentication related configuration
type TokenAuthConfig struct {
	JWTSignKey    string `mapstructure:"jwtSignKey"`
	JWTExpiresAt  int64  `mapstructure:"expiresAt"`
	AuthServerUrl string `mapstructure:"authServerUrl"`
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

// ListenerConfig contains app kafka topic listener related config
type ListenerConfig struct {
	GroupID  string   `mapstructure:"groupId"`
	Brokers  []string `mapstructure:"brokers"`
	Topic    string   `mapstructure:"topic"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

// ListenerConfig contains app kafka topic producer related config
type ProducerConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	Topic    string   `mapstructure:"topic"`
	Async    bool     `mapstructure:"async"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
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
	config.APIConfig.HypdApiConfig = config.ServerConfig.HypdApiConfig
	config.APPConfig.HypdApiConfig = config.ServerConfig.HypdApiConfig
	config.APPConfig.ElasticsearchConfig = config.ElasticsearchConfig
	config.SessionConfig.RedisConfig = config.RedisConfig
	return config
}
