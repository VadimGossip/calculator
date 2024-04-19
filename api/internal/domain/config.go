package domain

import "time"

type NetServerConfig struct {
	Host string
	Port int
}

type AMPQServerConfig struct {
	Url string
}

type Exchange struct {
	Name string
}

type Query struct {
	Name string
	DLX  string
	TTL  int
}

type QueryBind struct {
	ExchangeName string `mapstructure:"exchange_name"`
	QueryName    string `mapstructure:"query_name"`
	Key          string
}

type ConsumerCfg struct {
	ExchangeName string `mapstructure:"exchange_name"`
	RoutingKey   string
	QueryName    string `mapstructure:"query_name"`
}

type AMPQStructCfg struct {
	WorkExchange  Exchange `mapstructure:"work_exchange"`
	RetryExchange Exchange `mapstructure:"retry_exchange"`
	Queries       []Query
	QueryBinds    []QueryBind `mapstructure:"binds"`
	ConsumerCfg   ConsumerCfg `mapstructure:"consumer"`
}

type DbAgentCfg struct {
	Host string
	Port int
}

type ExpressionCfg struct {
	MaxLength            int
	HungTimeout          int
	AgentDownTimeout     int
	HungCheckPeriod      int
	AgentDownCheckPeriod int
}

type AuthConfig struct {
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
	Salt            string
	Secret          string
}

type Config struct {
	Expression       ExpressionCfg
	AppHttpServer    NetServerConfig
	DbAgentGrpc      DbAgentCfg
	Auth             AuthConfig
	AMPQServerConfig AMPQServerConfig
	AMPQStructCfg    AMPQStructCfg
}
