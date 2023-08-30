package model

type Config struct {
	AppPort string `mapstructure:"APP_PORT"`
	AppHost string `mapstructure:"APP_HOST"`
	AppPrefix string `mapstructure:"APP_PREFIX"`
	AppSecret string `mapstructure:"APP_SECRET"`
	AppEnv string `mapstructure:"APP_ENV"`
	MongoDBURI string `mapstructure:"MONGODB_URI"`
	MongoDBRetry int `mapstructure:"MONGODB_RETRY"`
	MongoDBRetryInterval int `mapstructure:"MONGODB_RETRY_INTERVAL"`
}

type Response struct {
	Message string `json:"message"`
	Status int `json:"status"`
	Data interface{} `json:"data"`
}