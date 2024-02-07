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

type Config struct {
	AppHttpServer    NetServerConfig
	AMPQServerConfig AMPQServerConfig
	AMPQStructCfg    AMPQStructCfg
}
