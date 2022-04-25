package configprovider

type Config struct {
	PSQLConfig
	HTTPKubeConfig
	AUTHKubeConfig
	ExternalServicesConfig
}

type PSQLConfig struct {
	DBString           string `envconfig:"PSQL_DB_URL"`
	MaxConnections     int    `envconfig:"PSQL_DB_MAX_CONNECTIONS"`
	MaxIdleConnections int    `envconfig:"PSQL_DB_MAX_IDLE_CONNECTIONS"`
}

type HTTPKubeConfig struct {
	Port string
}

type AUTHKubeConfig struct {
	FireKey string `envconfig:"FIREBASE_KEY"`
}

type ExternalServicesConfig struct {
	SendGridAPIKey string `envconfig:"SENDGRID_API_KEY"`
}
