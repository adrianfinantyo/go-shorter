package model

type Config struct {
	AppPort string `mapstructure:"APP_PORT"`
	AppHost string `mapstructure:"APP_HOST"`
	AppPrefix string `mapstructure:"APP_PREFIX"`
	AppSecret string `mapstructure:"APP_SECRET"`
	AppEnv string `mapstructure:"APP_ENV"`
	MongoDBURI string
	MongoDBHost string `mapstructure:"MONGODB_HOST"`
	MongoDBPort string `mapstructure:"MONGODB_PORT"`
	MongoDBName string `mapstructure:"MONGODB_NAME"`
	MongoDBUser string `mapstructure:"MONGODB_USER"`
	MongoDBPassword string `mapstructure:"MONGODB_PASS"`
	MongoDBRetry int `mapstructure:"MONGODB_RETRY"`
	MongoDBRetryInterval int `mapstructure:"MONGODB_RETRY_INTERVAL"`
}

type Response struct {
	Message string `json:"message"`
	Status int `json:"status"`
	Data interface{} `json:"data"`
}