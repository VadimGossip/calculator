package calculatorapi

type ClientService interface {
	SendStartEvalRequest(req *StartSubExpressionEvalRequest) (*StartSubExpressionEvalResponse, error)
	SendStopEvalRequest(req *StopSubExpressionEvalRequest) (*CommonResponse, error)
}

type HttpClient interface {
	SendPostRequest(endpoint string, request, response interface{}) error
}

type client struct {
	httpClient HttpClient
}

func NewClient(httpClient HttpClient) *client {
	return &client{httpClient: httpClient}
}

func (c *client) SendStartEvalRequest(req *StartSubExpressionEvalRequest) (*StartSubExpressionEvalResponse, error) {
	res := &StartSubExpressionEvalResponse{}
	err := c.httpClient.SendPostRequest("/sub_expression/start", req, res)
	return res, err
}

func (c *client) SendStopEvalRequest(req *StopSubExpressionEvalRequest) (*CommonResponse, error) {
	res := &CommonResponse{}
	err := c.httpClient.SendPostRequest("/sub_expression/stop", req, res)
	return res, err
}
