package writer

const (
	NewState = iota + 1
	ErrorState
	InProgressState
	OkState
)

type Expression struct {
	Id             int64  `json:"expression_id"`
	Value          string `json:"expression_value"`
	Result         string `json:"result"`
	State          string `json:"status"`
	CreatedAt      int64  `json:"created_at"`
	EvalStartedAt  int64  `json:"eval_started_at"`
	EvalFinishedAt int64  `json:"eval_finished_at"`
}
