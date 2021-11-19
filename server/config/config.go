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
	S3Config            S3Config            `mapstructure:"s3"`
	IVSConfig           IVSConfig           `mapstructure:"ivs"`
	ElasticsearchConfig ElasticsearchConfig `mapstructure:"elasticsearch"`
	HypdAPIConfig       HypdAPIConfig       `mapstructure:"hypdAPI"`
	SessionConfig       SessionConfig       `mapstructure:"session"`
	SentryConfig        SentryConfig        `mapstructure:"sentry"`
}

// HypdAPIConfig contains other HYPD service apis
type HypdAPIConfig struct {
	EntityAPI  string `mapstructure:"entityAPI"`
	CatalogAPI string `mapstructure:"catalogAPI"`
	Token      string `mapstructure:"token"`
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
}

// SentryConfig contains sentry related configuration
type SentryConfig struct {
	EnableSentry bool   `mapstructure:"enable"`
	DSN          string `mapstructure:"dsn"`
}

// APIConfig contains api package related configurations
type APIConfig struct {
	Mode               string `mapstructure:"mode"`
	EnableTestRoute    bool   `mapstructure:"enableTestRoute"`
	EnableMediaRoute   bool   `mapstructure:"enableMediaRoute"`
	EnableStaticRoute  bool   `mapstructure:"enableStaticRoute"`
	MaxRequestDataSize int    `mapstructure:"maxRequestDataSize"`
}

// APPConfig contains api package related configurations
type APPConfig struct {
	DatabaseConfig      DatabaseConfig
	S3Config            S3Config
	ElasticsearchConfig ElasticsearchConfig
	IVSConfig           IVSConfig
	HypdAPIConfig       HypdAPIConfig
	MediaConfig         ServiceConfig `mapstructure:"media"`
	ContentConfig       ServiceConfig `mapstructure:"content"`
	LiveConfig          ServiceConfig `mapstructure:"live"`
	SeriesConfig        ServiceConfig `mapstructure:"series"`
	CollectionConfig    ServiceConfig `mapstructure:"collection"`
	CategoryConfig      ServiceConfig `mapstructure:"category"`

	LiveCommentProducerConfig    ProducerConfig `mapstructure:"liveCommentProducer"`
	ContentFullProducerConfig    ProducerConfig `mapstructure:"contentFullProducer"`
	SeriesFullProducerConfig     ProducerConfig `mapstructure:"seriesFullProducer"`
	CollectionFullProducerConfig ProducerConfig `mapstructure:"collectionFullProducer"`

	LikeChangeConfig                  ListenerConfig `mapstructure:"likeChangesConsumer"`
	LikeChangeForSeriesConfig         ListenerConfig `mapstructure:"likeChangeForSeriesConsumer"`
	CommentChangeConfig               ListenerConfig `mapstructure:"commentChangesConsumer"`
	ViewChangeConfig                  ListenerConfig `mapstructure:"viewChangesConsumer"`
	BrandChangesConfig                ListenerConfig `mapstructure:"brandChangesConsumer"`
	InfluencerChangesConfig           ListenerConfig `mapstructure:"influencerChangesConsumer"`
	CatalogChangesConfig              ListenerConfig `mapstructure:"catalogChangesConsumer"`
	ContentChangesConfig              ListenerConfig `mapstructure:"contentChangesConsumer"`
	PebbleStatusChangeForSeriesConfig ListenerConfig `mapstructure:"pebbleStatusChangeForSeriesConsumer"`
	LiveCommentChangesConfig          ListenerConfig `mapstructure:"liveCommentChangesConsumer"`
	SeriesConsumerConfig              ListenerConfig `mapstructure:"seriesConsumer"`
	CollectionConsumerConfig          ListenerConfig `mapstructure:"collectionConsumer"`
}

type ContentStreamProcessorConfig struct {
	TopicProcessorName    string   `mapstructure:"topicProcessorName"`
	InputTopics           []string `mapstructure:"inputTopics"`
	InputPartitions       []int    `mapstructure:"inputPartitions"`
	CatalogTopicStream    string   `mapstructure:"catalogTopicStream"`
	BrandTopicStream      string   `mapstructure:"brandTopicStream"`
	InfluencerTopicStream string   `mapstructure:"influencerTopicStream"`
}

// ElasticsearchConfig contains elasticsearch related configurations
type ElasticsearchConfig struct {
	Endpoint                  string `mapstructure:"endpoint"`
	Username                  string `mapstructure:"username"`
	Password                  string `mapstructure:"password"`
	ContentFullIndex          string `mapstructure:"contentFullIndex"`
	PebbleSeriesFullIndex     string `mapstructure:"pebbleFullIndex"`
	PebbleCollectionFullIndex string `mapstructure:"pebbleCollectionFullIndex"`
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

// ListenerConfig contains app kafka topic producer related config
type ProducerConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	Topic    string   `mapstructure:"topic"`
	Async    bool     `mapstructure:"async"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

// TokenAuthConfig contains token authentication related configuration
type TokenAuthConfig struct {
	JWTSignKey   string `mapstructure:"jwtSignKey"`
	JWTExpiresAt int64  `mapstructure:"expiresAt"`
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

// CORSConfig contains cors related config
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowedOrigins"`
	AllowedMethods   []string `mapstructure:"allowedMethods"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	AllowedHeaders   []string `mapstructure:"allowedHeaders"`
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

// S3Config stores s3 configurations
type S3Config struct {
	Region               string        `mapstructure:"region"`
	AccessKeyID          string        `mapstructure:"accessKeyID"`
	SecretAccessKey      string        `mapstructure:"secretAccessKey"`
	ImageCloudfrontURL   string        `mapstructure:"imageCloudfrontUrl"`
	ImageUploadBucket    string        `mapstructure:"imageUploadBucket"`
	VideoUploadPath      string        `mapstructure:"videoUploadPath"`
	VideoUploadBucket    string        `mapstructure:"videoUploadBucket"`
	PresignedURLValidity time.Duration `mapstructure:"presignedURLValidity"`
}

// IVSConfig contains aws ivs related configuration
type IVSConfig struct {
	Region           string `mapstructure:"region"`
	ARN              string `mapstructure:"arn"`
	AccessKeyID      string `mapstructure:"accessKeyId"`
	SecretAccessKey  string `mapstructure:"secretAccessKey"`
	AuthorizeChannel bool   `mapstructure:"authorizeChannel"`
	LatencyMode      string `mapstructure:"latencyMode"`
	ChannelType      string `mapstructure:"channelType"`
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
	config.APPConfig.S3Config = config.S3Config
	config.APPConfig.IVSConfig = config.IVSConfig
	config.APPConfig.ElasticsearchConfig = config.ElasticsearchConfig
	config.APPConfig.HypdAPIConfig = config.HypdAPIConfig
	config.SessionConfig.RedisConfig = config.RedisConfig
	return config
}
