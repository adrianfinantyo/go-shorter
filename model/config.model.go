package model

type Config struct {
	AppPort string `mapstructure:"APP_PORT"`
	AppHost string `mapstructure:"APP_HOST"`
	AppPrefix string `mapstructure:"APP_PREFIX"`
	AppSecret string `mapstructure:"APP_SECRET"`
	AppEnv string `mapstructure:"APP_ENV"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName string `mapstructure:"DB_NAME"`
	MongoDBURI string `mapstructure:"MONGODB_URI"`
}

type Response struct {
	Message string `json:"message"`
	Status int `json:"status"`
	Data interface{} `json:"data"`
}