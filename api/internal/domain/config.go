package domain

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

type Config struct {
	Expression       ExpressionCfg
	AppHttpServer    NetServerConfig
	DbAgentGrpc      DbAgentCfg
	AMPQServerConfig AMPQServerConfig
	AMPQStructCfg    AMPQStructCfg
}
