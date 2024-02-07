package writer

const (
	NewState        = "New"
	ErrorState      = "Error"
	InProgressState = "In Progress"
	OkState         = "Ok"
)

type Expression struct {
	Id             int64  `json:"expression_id"`
	Value          string `json:"expression_value"`
	Result         *int   `json:"result"`
	State          string `json:"status"`
	CreatedAt      int64  `json:"created_at"`
	EvalStartedAt  int64  `json:"eval_started_at"`
	EvalFinishedAt int64  `json:"eval_finished_at"`
}
