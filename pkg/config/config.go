package config

// Config for the server
type Config struct {
	GRPCPort             string `mapstructure:"grpc_port"`
	HTTPPort             string `mapstructure:"http_port"`
	DatabaseHost         string `mapstructure:"db_host"`
	DatabasePort         string `mapstructure:"db_port"`
	DatabaseUser         string `mapstructure:"db_user"`
	DatabasePassword     string `mapstructure:"db_password"`
	DatabaseName         string `mapstructure:"db_name"`
	DatabaseSsl          string `mapstructure:"db_ssl"`
	LogLevel             int    `mapstructure:"log_level"`
	LogTimeFormat        string `mapstructure:"log_time_format"`
	JWTSecret            string `mapstructure:"jwt_secret"`
	JWTExpiration        int    `mapstructure:"jwt_expiration"`
	JWTRefreshExpiration int    `mapstructure:"jwt_refresh_expiration"`
}
