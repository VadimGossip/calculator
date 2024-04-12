package domain

type NetServerConfig struct {
	Host string
	Port int
}

type DbCfg struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}
type Config struct {
	Db            DbCfg
	AppGrpcServer NetServerConfig
}
